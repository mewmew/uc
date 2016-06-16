package irgen

// TODO: Add convenience functions for creating instruction in emit.go, to
// remove if err != nil { panic("foo") } from the irgen code.

import (
	"fmt"
	"strconv"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/instruction"
	irtypes "github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
	"github.com/mewmew/uc/ast"
	"github.com/mewmew/uc/ast/astutil"
	"github.com/mewmew/uc/sem"
	"github.com/mewmew/uc/token"
)

// Gen generates LLVM IR based on the syntax tree of the given file.
func Gen(file *ast.File, info *sem.Info) *ir.Module {
	return gen(file, info)
}

// === [ File scope ] ==========================================================

// gen generates LLVM IR based on the syntax tree of the given file.
func gen(file *ast.File, info *sem.Info) *ir.Module {
	m := NewModule(info)
	for _, decl := range file.Decls {
		// Ignore tentative definitions.
		if isTentativeDef(decl) {
			dbg.Printf("ignoring tentative definition of %q", decl.Name())
			continue
		}
		switch decl := decl.(type) {
		case *ast.FuncDecl:
			m.funcDecl(decl)
		case *ast.VarDecl:
			m.globalVarDecl(decl)
		case *ast.TypeDef:
			m.typeDef(decl)
		default:
			panic(fmt.Sprintf("support for %T not yet implemented", decl))
		}
	}
	return m.Module
}

// --- [ Function declaration ] ------------------------------------------------

// funcDecl lowers the given function declaration to LLVM IR, emitting code to
// m.
func (m *Module) funcDecl(n *ast.FuncDecl) {
	// Generate function signature.
	name := n.Name().String()
	typ := toIrType(n.Type())
	sig, ok := typ.(*irtypes.Func)
	if !ok {
		panic(fmt.Sprintf("invalid function type; expected *types.Func, got %T", typ))
	}
	f := NewFunction(name, sig)
	if !astutil.IsDef(n) {
		dbg.Printf("create function declaration: %v", n)
		// Emit function declaration.
		m.emitFunc(f)
		return
	}

	// Generate function body.
	dbg.Printf("create function definition: %v", n)
	m.funcBody(f, n.FuncType.Params, n.Body)
}

// funcBody lowers the given function declaration to LLVM IR, emitting code to
// m.
func (m *Module) funcBody(f *Function, params []*ast.VarDecl, body *ast.BlockStmt) {
	// Initialize function body.
	f.startBody()

	// TODO: Figure out a cleaner way to handle function parameters. The current
	// implementation makes use of and cross-references information from both
	// ast.FuncType.Params (for ast.Poc to be used in locals map) and
	// irtypes.Func.Params (value.Value). It should be possible to find a better
	// approach which only needs one of these two.

	// Emit local variable declarations for function parameters.
	for i, param := range f.Sig().Params() {
		p := m.funcParam(f, param)
		// Add mapping from parameter name to the corresponding allocated local
		// variable; i.e.
		//
		// For the following C code
		//
		//    void f(int a) {}
		//
		// with the following LLVM IR code
		//
		//    define void @f(i32 %a) {
		//    0:
		//       %1 = alloca i32
		//       store i32 %a, i32* %1
		//    }
		//
		// map from the parameter "a" to the allocated local variable "%1".
		dbg.Printf("create function parameter: %v", params[i])
		ident := params[i].Name()
		got := f.genUnique(ident)
		if ident.Name != got {
			panic(fmt.Sprintf("unable to generate identical function parameter name; expected %q, got %q", ident, got))
		}
		f.setIdentValue(ident, p)
	}

	// Generate function body.
	m.stmt(f, body)

	// Finalize function body.
	if err := f.endBody(); err != nil {
		panic(fmt.Sprintf("unable to finalize function body; %v", err))
	}

	// Emit function definition.
	m.emitFunc(f)
}

// funcParam lowers the given function parameter to LLVM IR, emitting code to f.
func (m *Module) funcParam(f *Function, param *irtypes.Param) value.Value {
	// Input:
	//    void f(int a) {
	//    }
	// Output:
	//    %1 = alloca i32
	//    store i32 %a, i32* %1
	allocaInst, err := instruction.NewAlloca(param.Type(), 1)
	if err != nil {
		panic(fmt.Sprintf("unable to create alloca instruction; %v", err))
	}
	// Emit local variable definition for the given function parameter.
	p := f.emitInst(allocaInst)
	storeInst, err := instruction.NewStore(param, p)
	if err != nil {
		panic(fmt.Sprintf("unable to create store instruction; %v", err))
	}
	f.curBlock.AppendInst(storeInst)
	return p
}

// --- [ Global variable declaration ] -----------------------------------------

// globalVarDecl lowers the given global variable declaration to LLVM IR,
// emitting code to m.
func (m *Module) globalVarDecl(n *ast.VarDecl) {
	// Input:
	//    int x;
	// Output:
	//    @x = global i32 0
	ident := n.Name()
	dbg.Printf("create global variable: %v", n)
	typ := toIrType(n.Type())
	var val value.Value
	switch {
	case n.Val != nil:
		panic("support for global variable initializer not yet implemented")
	case irtypes.IsInt(typ):
		var err error
		val, err = constant.NewInt(typ, "0")
		if err != nil {
			panic(fmt.Sprintf("unable to create integer constant; %v", err))
		}
	default:
		val = constant.NewZeroInitializer(typ)
	}
	global, err := ir.NewGlobalDef(ident.Name, val, false)
	if err != nil {
		panic(fmt.Sprintf("unable to create global variable definition %q", ident))
	}
	m.setIdentValue(ident, global)
	// Emit global variable definition.
	m.emitGlobal(global)
}

// --- [ Type definition ] -----------------------------------------------------

// typeDef lowers the given type definition to LLVM IR, emitting code to m.
func (m *Module) typeDef(def *ast.TypeDef) {
	panic("not yet implemented")
}

// === [ Function scope ] ======================================================

// --- [ Local variable definition ] -------------------------------------------

// localVarDef lowers the given local variable definition to LLVM IR, emitting
// code to f.
func (m *Module) localVarDef(f *Function, n *ast.VarDecl) {
	// Input:
	//    void f() {
	//       int a;           // <-- relevant line
	//    }
	// Output:
	//    %a = alloca i32
	ident := n.Name()
	dbg.Printf("create local variable: %v", n)
	typ := toIrType(n.Type())
	allocaInst, err := instruction.NewAlloca(typ, 1)
	if err != nil {
		panic(fmt.Sprintf("unable to create alloca instruction; %v", err))
	}
	// Emit local variable definition.
	f.emitLocal(ident, allocaInst)
	if n.Val != nil {
		panic("support for local variable definition initializer not yet implemented")
	}
}

// --- [ Statements ] ----------------------------------------------------------

// stmt lowers the given statement to LLVM IR, emitting code to f.
func (m *Module) stmt(f *Function, stmt ast.Stmt) {
	switch stmt := stmt.(type) {
	case *ast.BlockStmt:
		m.blockStmt(f, stmt)
		return
	case *ast.EmptyStmt:
		// nothing to do.
		return
	case *ast.ExprStmt:
		m.exprStmt(f, stmt)
		return
	case *ast.IfStmt:
		m.ifStmt(f, stmt)
		return
	case *ast.ReturnStmt:
		m.returnStmt(f, stmt)
		return
	case *ast.WhileStmt:
		m.whileStmt(f, stmt)
		return
	default:
		panic(fmt.Sprintf("support for %T not yet implemented", stmt))
	}
	panic("unreachable")
}

// blockStmt lowers the given block statement to LLVM IR, emitting code to f.
func (m *Module) blockStmt(f *Function, stmt *ast.BlockStmt) {
	for _, item := range stmt.Items {
		switch item := item.(type) {
		case ast.Decl:
			switch decl := item.(type) {
			case *ast.FuncDecl:
				panic(fmt.Sprintf("support for nested function declarations not yet implemented: %v", decl))
			case *ast.VarDecl:
				m.localVarDef(f, decl)
			case *ast.TypeDef:
				panic(fmt.Sprintf("support for scoped type definitions not yet implemented: %v", decl))
			}
		case ast.Stmt:
			m.stmt(f, item)
		}
	}
}

// exprStmt lowers the given expression statement to LLVM IR, emitting code to
// f.
func (m *Module) exprStmt(f *Function, stmt *ast.ExprStmt) {
	m.expr(f, stmt.X)
}

// ifStmt lowers the given if statement to LLVM IR, emitting code to f.
func (m *Module) ifStmt(f *Function, stmt *ast.IfStmt) {
	cond := m.cond(f, stmt.Cond)
	trueBranch := f.NewBasicBlock("")
	falseBranch := f.NewBasicBlock("")
	term, err := instruction.NewBr(cond, trueBranch, falseBranch)
	if err != nil {
		panic(fmt.Sprintf("unable to create br terminator; %v", err))
	}
	f.curBlock.SetTerm(term)
	f.curBlock = trueBranch

	m.stmt(f, stmt.Body)
	end := falseBranch
	f.curBlock = falseBranch

	if stmt.Else != nil {
		m.stmt(f, stmt.Else)
		end = f.NewBasicBlock("")
		falseBranch.emitJmp(end)
		f.curBlock = end
	}

	trueBranch.emitJmp(end)
}

// returnStmt lowers the given return statement to LLVM IR, emitting code to f.
func (m *Module) returnStmt(f *Function, stmt *ast.ReturnStmt) {
	// Input:
	//    int f() {
	//       return 42;       // <-- relevant line
	//    }
	// Output:
	//    ret i32 42
	if stmt.Result == nil {
		term, err := instruction.NewRet(irtypes.NewVoid(), nil)
		if err != nil {
			panic(fmt.Sprintf("unable to create ret terminator; %v", err))
		}
		f.curBlock.SetTerm(term)
		f.curBlock = nil
		return
	}
	result := m.expr(f, stmt.Result)
	term, err := instruction.NewRet(result.Type(), result)
	if err != nil {
		panic(fmt.Sprintf("unable to create ret terminator; %v", err))
	}
	f.curBlock.SetTerm(term)
	f.curBlock = nil
}

// whileStmt lowers the given while statement to LLVM IR, emitting code to f.
func (m *Module) whileStmt(f *Function, stmt *ast.WhileStmt) {
	condBranch := f.NewBasicBlock("")
	condJmp, err := instruction.NewJmp(condBranch)
	if err != nil {
		panic(fmt.Sprintf("unable to create jmp terminator; %v", err))
	}
	f.curBlock.SetTerm(condJmp)
	f.curBlock = condBranch
	cond := m.cond(f, stmt.Cond)
	bodyBranch := f.NewBasicBlock("")
	endBranch := f.NewBasicBlock("")
	term, err := instruction.NewBr(cond, bodyBranch, endBranch)
	if err != nil {
		panic(fmt.Sprintf("unable to create br terminator; %v", err))
	}
	f.curBlock.SetTerm(term)
	f.curBlock = bodyBranch

	m.stmt(f, stmt.Body)

	f.curBlock.SetTerm(condJmp)
	f.curBlock = endBranch
}

// --- [ Expressions ] ----------------------------------------------------------

// cond lowers the given condition expression to LLVM IR, emitting code to f.
func (m *Module) cond(f *Function, expr ast.Expr) value.Value {
	cond := m.expr(f, expr)
	if irtypes.IsBool(cond.Type()) {
		return cond
	}
	// Create boolean expression if cond is not already of boolean type.
	//
	//    cond != 0
	// zero is the integer constant 0.
	zero := constZero(cond.Type())
	icmpInst, err := instruction.NewICmp(instruction.ICondNE, cond, zero)
	if err != nil {
		panic(fmt.Sprintf("unable to create icmp instruction; %v", err))
	}
	return f.emitInst(icmpInst)
}

// expr lowers the given expression to LLVM IR, emitting code to f.
func (m *Module) expr(f *Function, expr ast.Expr) value.Value {
	switch expr := expr.(type) {
	case *ast.BasicLit:
		return m.basicLit(f, expr)
	case *ast.BinaryExpr:
		return m.binaryExpr(f, expr)
	case *ast.CallExpr:
		return m.callExpr(f, expr)
	case *ast.Ident:
		return m.identUse(f, expr)
	case *ast.IndexExpr:
		return m.indexExprUse(f, expr)
	case *ast.ParenExpr:
		return m.expr(f, expr.X)
	case *ast.UnaryExpr:
		return m.unaryExpr(f, expr)
	default:
		panic(fmt.Sprintf("support for type %T not yet implemented", expr))
	}
	panic("unreachable")
}

// basicLit lowers the given basic literal to LLVM IR, emitting code to f.
func (m *Module) basicLit(f *Function, n *ast.BasicLit) value.Value {
	typ := m.typeOf(n)
	switch n.Kind {
	case token.CharLit:
		s, err := strconv.Unquote(n.Val)
		if err != nil {
			panic(fmt.Sprintf("unable to unquote character literal; %v", err))
		}
		val, err := constant.NewInt(typ, strconv.Itoa(int(s[0])))
		if err != nil {
			panic(fmt.Sprintf("unable to create integer constant; %v", err))
		}
		return val
	case token.IntLit:
		val, err := constant.NewInt(typ, n.Val)
		if err != nil {
			panic(fmt.Sprintf("unable to create integer constant; %v", err))
		}
		return val
	default:
		panic(fmt.Sprintf("support for basic literal kind %v not yet implemented", n.Kind))
	}
	panic("unreachable")
}

// binaryExpr lowers the given binary expression to LLVM IR, emitting code to f.
func (m *Module) binaryExpr(f *Function, n *ast.BinaryExpr) value.Value {
	switch n.Op {
	// +
	case token.Add:
		x, y := m.expr(f, n.X), m.expr(f, n.Y)
		addInst, err := instruction.NewAdd(x, y)
		if err != nil {
			panic(fmt.Sprintf("unable to create add instruction; %v", err))
		}
		// Emit add instruction.
		return f.emitInst(addInst)

	// -
	case token.Sub:
		x, y := m.expr(f, n.X), m.expr(f, n.Y)
		subInst, err := instruction.NewSub(x, y)
		if err != nil {
			panic(fmt.Sprintf("unable to create sub instruction; %v", err))
		}
		// Emit sub instruction.
		return f.emitInst(subInst)

	// *
	case token.Mul:
		x, y := m.expr(f, n.X), m.expr(f, n.Y)
		mulInst, err := instruction.NewMul(x, y)
		if err != nil {
			panic(fmt.Sprintf("unable to create mul instruction; %v", err))
		}
		// Emit mul instruction.
		return f.emitInst(mulInst)

	// /
	case token.Div:
		x, y := m.expr(f, n.X), m.expr(f, n.Y)
		// TODO: Add support for unsigned division.
		sdivInst, err := instruction.NewSDiv(x, y)
		if err != nil {
			panic(fmt.Sprintf("unable to create sdiv instruction; %v", err))
		}
		// Emit sdiv instruction.
		return f.emitInst(sdivInst)

	// <
	case token.Lt:
		x, y := m.expr(f, n.X), m.expr(f, n.Y)
		icmpInst, err := instruction.NewICmp(instruction.ICondSLT, x, y)
		if err != nil {
			panic(fmt.Sprintf("unable to create icmp instruction; %v", err))
		}
		// Emit icmp instruction.
		cond := f.emitInst(icmpInst)
		zextInst, err := instruction.NewZExt(cond, x.Type())
		if err != nil {
			panic(fmt.Sprintf("unable to create zext instruction; %v", err))
		}
		// Emit zext instruction.
		return f.emitInst(zextInst)

	// >
	case token.Gt:
		x, y := m.expr(f, n.X), m.expr(f, n.Y)
		icmpInst, err := instruction.NewICmp(instruction.ICondSGT, x, y)
		if err != nil {
			panic(fmt.Sprintf("unable to create icmp instruction; %v", err))
		}
		// Emit icmp instruction.
		cond := f.emitInst(icmpInst)
		zextInst, err := instruction.NewZExt(cond, x.Type())
		if err != nil {
			panic(fmt.Sprintf("unable to create zext instruction; %v", err))
		}
		// Emit zext instruction.
		return f.emitInst(zextInst)

	// <=
	case token.Le:
		x, y := m.expr(f, n.X), m.expr(f, n.Y)
		icmpInst, err := instruction.NewICmp(instruction.ICondSLE, x, y)
		if err != nil {
			panic(fmt.Sprintf("unable to create icmp instruction; %v", err))
		}
		// Emit icmp instruction.
		cond := f.emitInst(icmpInst)
		zextInst, err := instruction.NewZExt(cond, x.Type())
		if err != nil {
			panic(fmt.Sprintf("unable to create zext instruction; %v", err))
		}
		// Emit zext instruction.
		return f.emitInst(zextInst)

	// >=
	case token.Ge:
		x, y := m.expr(f, n.X), m.expr(f, n.Y)
		icmpInst, err := instruction.NewICmp(instruction.ICondSGE, x, y)
		if err != nil {
			panic(fmt.Sprintf("unable to create icmp instruction; %v", err))
		}
		// Emit icmp instruction.
		cond := f.emitInst(icmpInst)
		zextInst, err := instruction.NewZExt(cond, x.Type())
		if err != nil {
			panic(fmt.Sprintf("unable to create zext instruction; %v", err))
		}
		// Emit zext instruction.
		return f.emitInst(zextInst)

	// !=
	case token.Ne:
		x, y := m.expr(f, n.X), m.expr(f, n.Y)
		icmpInst, err := instruction.NewICmp(instruction.ICondNE, x, y)
		if err != nil {
			panic(fmt.Sprintf("unable to create icmp instruction; %v", err))
		}
		// Emit icmp instruction.
		cond := f.emitInst(icmpInst)
		zextInst, err := instruction.NewZExt(cond, x.Type())
		if err != nil {
			panic(fmt.Sprintf("unable to create zext instruction; %v", err))
		}
		// Emit zext instruction.
		return f.emitInst(zextInst)

	// ==
	case token.Eq:
		x, y := m.expr(f, n.X), m.expr(f, n.Y)
		icmpInst, err := instruction.NewICmp(instruction.ICondEq, x, y)
		if err != nil {
			panic(fmt.Sprintf("unable to create icmp instruction; %v", err))
		}
		// Emit icmp instruction.
		cond := f.emitInst(icmpInst)
		zextInst, err := instruction.NewZExt(cond, x.Type())
		if err != nil {
			panic(fmt.Sprintf("unable to create zext instruction; %v", err))
		}
		// Emit zext instruction.
		return f.emitInst(zextInst)

	// &&
	case token.Land:
		x := m.cond(f, n.X)

		start := f.curBlock
		trueBranch := f.NewBasicBlock("")
		end := f.NewBasicBlock("")
		term, err := instruction.NewBr(x, trueBranch, end)
		if err != nil {
			panic(fmt.Sprintf("unable to create br terminator; %v", err))
		}
		f.curBlock.SetTerm(term)
		f.curBlock = trueBranch

		y := m.cond(f, n.Y)
		trueBranch.emitJmp(end)
		f.curBlock = end

		var incs []*instruction.Incoming
		zero := constZero(irtypes.I1)
		inc, err := instruction.NewIncoming(zero, start)
		if err != nil {
			panic(fmt.Sprintf("unable to create incoming value; %v", err))
		}
		incs = append(incs, inc)
		inc, err = instruction.NewIncoming(y, trueBranch)
		if err != nil {
			panic(fmt.Sprintf("unable to create incoming value; %v", err))
		}
		incs = append(incs, inc)
		phiInst, err := instruction.NewPHI(irtypes.I1, incs)
		if err != nil {
			panic(fmt.Sprintf("unable to create br terminator; %v", err))
		}
		// Emit phi instruction.
		result := f.emitInst(phiInst)
		zextInst, err := instruction.NewZExt(result, m.typeOf(n))
		if err != nil {
			panic(fmt.Sprintf("unable to create zext instruction; %v", err))
		}
		// Emit zext instruction.
		return f.emitInst(zextInst)

	// =
	case token.Assign:
		y := m.expr(f, n.Y)
		switch expr := n.X.(type) {
		case *ast.Ident:
			m.identDef(f, expr, y)
		case *ast.IndexExpr:
			m.indexExprDef(f, expr, y)
		default:
			panic(fmt.Sprintf("support for assignment to type %T not yet implemented", expr))
		}
		return y

	default:
		panic(fmt.Sprintf("support for binary operator %v not yet implemented", n.Op))
	}
	panic("unreachable")
}

// callExpr lowers the given identifier to LLVM IR, emitting code to f.
func (m *Module) callExpr(f *Function, callExpr *ast.CallExpr) value.Value {
	// val := m.valueFromIdent(f, callExpr.Name.Decl.Name())
	// typ := val.Type().(*irtypes.Func)
	typ := toIrType(callExpr.Name.Decl.Type()).(*irtypes.Func)
	var args []value.Value
	for _, arg := range callExpr.Args {
		fmt.Printf("arg: %T, %#v\n", arg, arg)
		expr := m.expr(f, arg)
		args = append(args, expr)
		// TODO: Add cast
	}
	inst, err := instruction.NewCall(typ.Result(), callExpr.Name.String(), args)
	if err != nil {
		panic(fmt.Sprintf("unable to create call instruction; %v", err))
	}
	return f.emitInst(inst)
}

// ident lowers the given identifier to LLVM IR, emitting code to f.
func (m *Module) ident(f *Function, ident *ast.Ident) value.Value {
	switch typ := m.typeOf(ident).(type) {
	case *irtypes.Array:
		array := m.valueFromIdent(f, ident)
		zero := constZero(irtypes.I64)
		indices := []value.Value{zero, zero}
		gepInst, err := instruction.NewGetElementPtr(typ, array, indices)
		if err != nil {
			panic(fmt.Sprintf("unable to create getelementptr instruction; %v", err))
		}
		return f.emitInst(gepInst)
	default:
		return m.valueFromIdent(f, ident)
	}
}

// identUse lowers the given identifier usage to LLVM IR, emitting code to f.
func (m *Module) identUse(f *Function, ident *ast.Ident) value.Value {
	v := m.ident(f, ident)
	typ := m.typeOf(ident)
	if isRef(typ) {
		return v
	}
	loadInst, err := instruction.NewLoad(typ, v)
	if err != nil {
		panic(fmt.Sprintf("unable to create load instruction; %v", err))
	}
	// Emit load instruction.
	return f.emitInst(loadInst)
}

// isRef reports whether the given type is a reference type; e.g. pointer or
// array.
func isRef(typ irtypes.Type) bool {
	switch typ.(type) {
	case *irtypes.Array:
		return true
	case *irtypes.Pointer:
		return true
	default:
		return false
	}
}

// identDef lowers the given identifier definition to LLVM IR, emitting code to
// f.
func (m *Module) identDef(f *Function, ident *ast.Ident, v value.Value) {
	addr := m.ident(f, ident)
	storeInst, err := instruction.NewStore(v, addr)
	if err != nil {
		panic(fmt.Sprintf("unable to create store instruction; %v", err))
	}
	f.curBlock.AppendInst(storeInst)
}

// indexExpr lowers the given index expression to LLVM IR, emitting code to f.
func (m *Module) indexExpr(f *Function, n *ast.IndexExpr) value.Value {
	index := m.expr(f, n.Index)
	// Extend the index to a 64-bit integer.
	if c, ok := index.(constant.Constant); ok {
		// The index is guaranteed to be of integer type.
		var err error
		index, err = constant.NewInt(irtypes.I64, c.ValueString())
		if err != nil {
			panic(fmt.Sprintf("unable to create integer constant; %v", err))
		}
	} else {
		// TODO: Use zext for unsigned values and sext for signed values.
		sextInst, err := instruction.NewSExt(index, irtypes.I64)
		if err != nil {
			panic(fmt.Sprintf("unable to create sext instruction; %v", err))
		}
		index = f.emitInst(sextInst)
	}
	typ := m.typeOf(n.Name)
	array := m.valueFromIdent(f, n.Name)
	zero := constZero(irtypes.I64)
	indices := []value.Value{zero, index}
	gepInst, err := instruction.NewGetElementPtr(typ, array, indices)
	if err != nil {
		panic(fmt.Sprintf("unable to create getelementptr instruction; %v", err))
	}
	return f.emitInst(gepInst)
}

// indexExprUse lowers the given index expression usage to LLVM IR, emitting
// code to f.
func (m *Module) indexExprUse(f *Function, n *ast.IndexExpr) value.Value {
	v := m.indexExpr(f, n)
	typ := m.typeOf(n)
	if isRef(typ) {
		return v
	}
	loadInst, err := instruction.NewLoad(typ, v)
	if err != nil {
		panic(fmt.Sprintf("unable to create load instruction; %v", err))
	}
	// Emit load instruction.
	return f.emitInst(loadInst)
}

// indexExprDef lowers the given identifier expression definition to LLVM IR,
// emitting code to f.
func (m *Module) indexExprDef(f *Function, n *ast.IndexExpr, v value.Value) {
	panic("indexExprDef: not yet implemented")
}

// unaryExpr lowers the given unary expression to LLVM IR, emitting code to f.
func (m *Module) unaryExpr(f *Function, n *ast.UnaryExpr) value.Value {
	switch n.Op {
	// -expr
	case token.Sub:
		// Input:
		//    void f() {
		//       int x;
		//       -x;              // <-- relevant line
		//    }
		// Output:
		//    %2 = sub i32 0, %1
		expr := m.expr(f, n.X)
		zero := constZero(expr.Type())
		subInst, err := instruction.NewSub(zero, expr)
		if err != nil {
			panic(fmt.Sprintf("unable to create sub instruction; %v", err))
		}
		// Emit sub instruction.
		return f.emitInst(subInst)
	// !expr
	case token.Not:
		// TODO: Replace `(x != 0) ^ 1` with `x == 0`. Using the former for now to
		// simplify test cases, as they are generated by Clang.

		// Input:
		//    int g() {
		//       int y;
		//       !y;              // <-- relevant line
		//    }
		// Output:
		//    %2 = icmp ne i32 %1, 0
		//    %3 = xor i1 %2, true
		cond := m.cond(f, n.X)
		one := constOne(cond.Type())
		xorInst, err := instruction.NewXor(cond, one)
		if err != nil {
			panic(fmt.Sprintf("unable to create xor instruction; %v", err))
		}
		// Emit xor instruction.
		notCond := f.emitInst(xorInst)
		zextInst, err := instruction.NewZExt(notCond, m.typeOf(n.X))
		if err != nil {
			panic(fmt.Sprintf("unable to create zext instruction; %v", err))
		}
		// Emit zext instruction.
		return f.emitInst(zextInst)
	default:
		panic(fmt.Sprintf("support for unary operator %v not yet implemented", n.Op))
	}
	panic("unreachable")
}
