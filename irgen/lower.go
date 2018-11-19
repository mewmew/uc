package irgen

// TODO: Add convenience functions for creating instruction in emit.go, to
// remove if err != nil { panic("foo") } from the irgen code.

import (
	"fmt"
	"strconv"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/enum"
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
	ident := n.Name()
	name := ident.String()
	typ := toIrType(n.Type())
	sig, ok := typ.(*irtypes.FuncType)
	if !ok {
		panic(fmt.Sprintf("invalid function type; expected *types.FuncType, got %T", typ))
	}
	var params []*ir.Param
	for i, p := range n.FuncType.Params {
		paramType := sig.Params[i]
		param := ir.NewParam(p.Name().String(), paramType)
		params = append(params, param)
	}
	f := NewFunction(name, sig.RetType, params...)
	if !astutil.IsDef(n) {
		dbg.Printf("create function declaration: %v", n)
		// Emit function declaration.
		m.emitFunc(f)
		return
	}
	m.setIdentValue(ident, f.Function)

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
	for i, param := range f.Params {
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
func (m *Module) funcParam(f *Function, param *ir.Param) value.Value {
	// Input:
	//    void f(int a) {
	//    }
	// Output:
	//    %1 = alloca i32
	//    store i32 %a, i32* %1
	addr := f.curBlock.NewAlloca(param.Type())
	f.curBlock.NewStore(param, addr)
	return addr
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
	var val constant.Constant
	switch {
	case n.Val != nil:
		panic("support for global variable initializer not yet implemented")
	default:
		if intType, ok := typ.(*irtypes.IntType); ok {
			val = constant.NewInt(intType, 0)
		} else {
			val = constant.NewZeroInitializer(typ)
		}
	}
	global := ir.NewGlobalDef(ident.Name, val)
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
	allocaInst := f.curBlock.NewAlloca(typ)
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
	trueBranch := f.NewBlock("")
	end := f.NewBlock("")
	falseBranch := end
	if stmt.Else != nil {
		falseBranch = f.NewBlock("")
	}
	termCondBr := ir.NewCondBr(cond, trueBranch.BasicBlock, falseBranch.BasicBlock)
	f.curBlock.SetTerm(termCondBr)
	f.curBlock = trueBranch
	m.stmt(f, stmt.Body)
	// Emit jump if body doesn't end with return statement (i.e. the current
	// basic block is none nil).
	if f.curBlock != nil {
		termBr := ir.NewBr(end.BasicBlock)
		f.curBlock.SetTerm(termBr)
	}
	if stmt.Else != nil {
		f.curBlock = falseBranch
		m.stmt(f, stmt.Else)
		// Emit jump if body doesn't end with return statement (i.e. the current
		// basic block is none nil).
		if f.curBlock != nil {
			termBr := ir.NewBr(end.BasicBlock)
			f.curBlock.SetTerm(termBr)
		}
	}
	f.curBlock = end
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
		termRet := ir.NewRet(nil)
		f.curBlock.SetTerm(termRet)
		f.curBlock = nil
		return
	}
	result := m.expr(f, stmt.Result)
	// Implicit conversion.
	resultType := f.Sig.RetType
	result = m.convert(f, result, resultType)
	termRet := ir.NewRet(result)
	f.curBlock.SetTerm(termRet)
	f.curBlock = nil
}

// whileStmt lowers the given while statement to LLVM IR, emitting code to f.
func (m *Module) whileStmt(f *Function, stmt *ast.WhileStmt) {
	condBranch := f.NewBlock("")
	termBr := ir.NewBr(condBranch.BasicBlock)
	f.curBlock.SetTerm(termBr)
	f.curBlock = condBranch
	cond := m.cond(f, stmt.Cond)
	bodyBranch := f.NewBlock("")
	endBranch := f.NewBlock("")
	termCondBr := ir.NewCondBr(cond, bodyBranch.BasicBlock, endBranch.BasicBlock)
	f.curBlock.SetTerm(termCondBr)
	f.curBlock = bodyBranch
	m.stmt(f, stmt.Body)
	// Emit jump if body doesn't end with return statement (i.e. the current
	// basic block is none nil).
	if f.curBlock != nil {
		termBr := ir.NewBr(condBranch.BasicBlock)
		f.curBlock.SetTerm(termBr)
	}
	f.curBlock = endBranch
}

// --- [ Expressions ] ----------------------------------------------------------

// cond lowers the given condition expression to LLVM IR, emitting code to f.
func (m *Module) cond(f *Function, expr ast.Expr) value.Value {
	cond := m.expr(f, expr)
	if cond.Type().Equal(irtypes.I1) {
		return cond
	}
	// Create boolean expression if cond is not already of boolean type.
	//
	//    cond != 0
	// zero is the integer constant 0.
	zero := constZero(cond.Type())
	return f.curBlock.NewICmp(enum.IPredNE, cond, zero)
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
		intType, ok := typ.(*irtypes.IntType)
		if !ok {
			panic(fmt.Errorf("invalid character literal type; expected *types.IntType, got %T", typ))
		}
		return constant.NewInt(intType, int64(s[0]))
	case token.IntLit:
		intType, ok := typ.(*irtypes.IntType)
		if !ok {
			panic(fmt.Errorf("invalid integer literal type; expected *types.IntType, got %T", typ))
		}
		c, err := constant.NewIntFromString(intType, n.Val)
		if err != nil {
			panic(fmt.Errorf("unable to parse integer literal %q; %v", n.Val, err))
		}
		return c
	default:
		panic(fmt.Sprintf("support for basic literal kind %v not yet implemented", n.Kind))
	}
}

// binaryExpr lowers the given binary expression to LLVM IR, emitting code to f.
func (m *Module) binaryExpr(f *Function, n *ast.BinaryExpr) value.Value {
	switch n.Op {
	// +
	case token.Add:
		x, y := m.expr(f, n.X), m.expr(f, n.Y)
		x, y = m.implicitConversion(f, x, y)
		return f.curBlock.NewAdd(x, y)

	// -
	case token.Sub:
		x, y := m.expr(f, n.X), m.expr(f, n.Y)
		x, y = m.implicitConversion(f, x, y)
		return f.curBlock.NewSub(x, y)

	// *
	case token.Mul:
		x, y := m.expr(f, n.X), m.expr(f, n.Y)
		x, y = m.implicitConversion(f, x, y)
		return f.curBlock.NewMul(x, y)

	// /
	case token.Div:
		x, y := m.expr(f, n.X), m.expr(f, n.Y)
		x, y = m.implicitConversion(f, x, y)
		// TODO: Add support for unsigned division.
		return f.curBlock.NewSDiv(x, y)

	// <
	case token.Lt:
		x, y := m.expr(f, n.X), m.expr(f, n.Y)
		x, y = m.implicitConversion(f, x, y)
		return f.curBlock.NewICmp(enum.IPredSLT, x, y)

	// >
	case token.Gt:
		x, y := m.expr(f, n.X), m.expr(f, n.Y)
		x, y = m.implicitConversion(f, x, y)
		return f.curBlock.NewICmp(enum.IPredSGT, x, y)

	// <=
	case token.Le:
		x, y := m.expr(f, n.X), m.expr(f, n.Y)
		x, y = m.implicitConversion(f, x, y)
		return f.curBlock.NewICmp(enum.IPredSLE, x, y)

	// >=
	case token.Ge:
		x, y := m.expr(f, n.X), m.expr(f, n.Y)
		x, y = m.implicitConversion(f, x, y)
		return f.curBlock.NewICmp(enum.IPredSGE, x, y)

	// !=
	case token.Ne:
		x, y := m.expr(f, n.X), m.expr(f, n.Y)
		x, y = m.implicitConversion(f, x, y)
		return f.curBlock.NewICmp(enum.IPredNE, x, y)

	// ==
	case token.Eq:
		x, y := m.expr(f, n.X), m.expr(f, n.Y)
		x, y = m.implicitConversion(f, x, y)
		return f.curBlock.NewICmp(enum.IPredEQ, x, y)

	// &&
	case token.Land:
		x := m.cond(f, n.X)

		start := f.curBlock
		trueBranch := f.NewBlock("")
		end := f.NewBlock("")
		termCondBr := ir.NewCondBr(x, trueBranch.BasicBlock, end.BasicBlock)
		f.curBlock.SetTerm(termCondBr)
		f.curBlock = trueBranch

		y := m.cond(f, n.Y)
		termBr := ir.NewBr(end.BasicBlock)
		trueBranch.SetTerm(termBr)
		f.curBlock = end

		var incs []*ir.Incoming
		zero := constZero(irtypes.I1)
		inc := ir.NewIncoming(zero, start.BasicBlock)
		incs = append(incs, inc)
		inc = ir.NewIncoming(y, trueBranch.BasicBlock)
		incs = append(incs, inc)
		return f.curBlock.NewPhi(incs...)

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
}

// callExpr lowers the given identifier to LLVM IR, emitting code to f.
func (m *Module) callExpr(f *Function, callExpr *ast.CallExpr) value.Value {
	typ := toIrType(callExpr.Name.Decl.Type())
	sig, ok := typ.(*irtypes.FuncType)
	if !ok {
		panic(fmt.Sprintf("invalid function type; expected *types.FuncType, got %T", typ))
	}
	params := sig.Params
	result := sig.RetType
	// TODO: Validate result against function return type.
	_ = result
	var args []value.Value
	// TODO: Add support for variadic arguments.
	for i, arg := range callExpr.Args {
		expr := m.expr(f, arg)
		expr = m.convert(f, expr, params[i])
		args = append(args, expr)
	}
	v := m.valueFromIdent(f, callExpr.Name)
	callee, ok := v.(*ir.Function)
	if !ok {
		panic(fmt.Sprintf("invalid callee type; expected *ir.Function, got %T", v))
	}
	return f.curBlock.NewCall(callee, args...)
}

// ident lowers the given identifier to LLVM IR, emitting code to f.
func (m *Module) ident(f *Function, ident *ast.Ident) value.Value {
	switch typ := m.typeOf(ident).(type) {
	case *irtypes.ArrayType:
		array := m.valueFromIdent(f, ident)
		zero := constZero(irtypes.I64)
		indices := []value.Value{zero, zero}

		// Emit getelementptr instruction.
		if m.isGlobal(ident) {
			var is []*constant.Index
			for _, index := range indices {
				i, ok := index.(constant.Constant)
				if !ok {
					break
				}
				idx := constant.NewIndex(i)
				is = append(is, idx)
			}
			if len(is) == len(indices) {
				// In accordance with Clang, emit getelementptr constant expressions
				// for global variables.
				// TODO: Validate typ against array.
				_ = typ
				if array, ok := array.(constant.Constant); ok {
					return constant.NewGetElementPtr(array, is...)
				}
				panic(fmt.Sprintf("invalid constant array type; expected constant.Constant, got %T", array))
			}
		}
		return f.curBlock.NewGetElementPtr(array, indices...)
	case *irtypes.PointerType:
		// Emit load instruction.
		// TODO: Validate typ against srcAddr.Elem().
		return f.curBlock.NewLoad(m.valueFromIdent(f, ident))
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
	// TODO: Validate typ against v.Elem()
	return f.curBlock.NewLoad(v)
}

// identDef lowers the given identifier definition to LLVM IR, emitting code to
// f.
func (m *Module) identDef(f *Function, ident *ast.Ident, v value.Value) {
	addr := m.ident(f, ident)
	addrType, ok := addr.Type().(*irtypes.PointerType)
	if !ok {
		panic(fmt.Sprintf("invalid pointer type; expected *types.PointerType, got %T", addr.Type()))
	}
	v = m.convert(f, v, addrType.ElemType)
	f.curBlock.NewStore(v, addr)
}

// indexExpr lowers the given index expression to LLVM IR, emitting code to f.
func (m *Module) indexExpr(f *Function, n *ast.IndexExpr) value.Value {
	index := m.expr(f, n.Index)
	// Extend the index to a 64-bit integer.
	if !irtypes.Equal(index.Type(), irtypes.I64) {
		index = m.convert(f, index, irtypes.I64)
	}
	typ := m.typeOf(n.Name)
	array := m.valueFromIdent(f, n.Name)

	// Dereference pointer pointer.
	elem := typ
	addr := array
	zero := constZero(irtypes.I64)
	indices := []value.Value{zero, index}
	if typ, ok := typ.(*irtypes.PointerType); ok {
		elem = typ.ElemType

		// Emit load instruction.
		// TODO: Validate typ against array.Elem().
		addr = f.curBlock.NewLoad(array)
		indices = []value.Value{index}
	}

	// Emit getelementptr instruction.
	if m.isGlobal(n.Name) {
		var is []*constant.Index
		for _, index := range indices {
			i, ok := index.(constant.Constant)
			if !ok {
				break
			}
			idx := constant.NewIndex(i)
			is = append(is, idx)
		}
		if len(is) == len(indices) {
			// In accordance with Clang, emit getelementptr constant expressions
			// for global variables.
			// TODO: Validate elem against addr.
			_ = elem
			if addr, ok := addr.(constant.Constant); ok {
				return constant.NewGetElementPtr(addr, is...)
			}
			panic(fmt.Sprintf("invalid constant address type; expected constant.Constant, got %T", addr))
		}
	}
	// TODO: Validate elem against array.Elem().
	return f.curBlock.NewGetElementPtr(addr, indices...)
}

// indexExprUse lowers the given index expression usage to LLVM IR, emitting
// code to f.
func (m *Module) indexExprUse(f *Function, n *ast.IndexExpr) value.Value {
	v := m.indexExpr(f, n)
	typ := m.typeOf(n)
	if isRef(typ) {
		return v
	}
	// TODO: Validate typ against v.Elem().
	return f.curBlock.NewLoad(v)
}

// indexExprDef lowers the given identifier expression definition to LLVM IR,
// emitting code to f.
func (m *Module) indexExprDef(f *Function, n *ast.IndexExpr, v value.Value) {
	addr := m.indexExpr(f, n)
	addrType, ok := addr.Type().(*irtypes.PointerType)
	if !ok {
		panic(fmt.Sprintf("invalid pointer type; expected *types.PointerType, got %T", addr.Type()))
	}
	v = m.convert(f, v, addrType.ElemType)
	f.curBlock.NewStore(v, addr)
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
		return f.curBlock.NewSub(zero, expr)
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
		notCond := f.curBlock.NewXor(cond, one)
		return f.curBlock.NewZExt(notCond, m.typeOf(n.X))
	default:
		panic(fmt.Sprintf("support for unary operator %v not yet implemented", n.Op))
	}
}
