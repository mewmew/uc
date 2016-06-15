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
	m.funcBody(f, n.Body)
}

// funcBody lowers the given function declaration to LLVM IR, emitting code to
// m.
func (m *Module) funcBody(f *Function, body *ast.BlockStmt) {
	// Initialize function body.
	f.startBody()

	// Generate function body.
	m.stmt(f, body)

	// Finalize function body.
	if err := f.endBody(); err != nil {
		panic(fmt.Sprintf("unable to finalize function body; %v", err))
	}

	// Emit function definition.
	m.emitFunc(f)
}

// --- [ Global variable declaration ] -----------------------------------------

// globalVarDecl lowers the given global variable declaration to LLVM IR,
// emitting code to m.
func (m *Module) globalVarDecl(n *ast.VarDecl) {
	// Input:
	//    int x;
	// Output:
	//    @x = global i32 0
	name := n.Name().Name
	dbg.Printf("create global variable: %v", n)
	typ := toIrType(n.Type())
	var val value.Value
	switch {
	case n.Val != nil:
		val = m.constExpr(n.Val)
	case irtypes.IsInt(typ):
		var err error
		val, err = constant.NewInt(typ, "0")
		if err != nil {
			panic(fmt.Sprintf("unable to create integer constant; %v", err))
		}
	default:
		val = constant.NewZeroInitializer(typ)
	}
	global := ir.NewGlobalDef(name, val, false)
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
	name := n.Name().Name
	dbg.Printf("create local variable: %v", n)
	typ := toIrType(n.Type())
	inst, err := instruction.NewAlloca(typ, 1)
	if err != nil {
		panic(fmt.Sprintf("unable to create alloca instruction; %v", err))
	}
	// Emit local variable definition.
	f.emitLocal(name, inst)
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
	panic("not yet implemented")
}

// ifStmt lowers the given if statement to LLVM IR, emitting code to f.
func (m *Module) ifStmt(f *Function, stmt *ast.IfStmt) {
	cond := m.expr(f, stmt.Cond)
	// Create boolean expression if cond is not already of boolean type.
	//
	//    cond != 0
	if !irtypes.IsBool(cond.Type()) {
		// zero is the integer constant 0.
		zero := constZero(cond.Type())
		inst, err := instruction.NewICmp(instruction.ICondNE, cond, zero)
		if err != nil {
			panic(fmt.Sprintf("unable to create icmp instruction; %v", err))
		}
		cond = f.emitInst(inst)
	}
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
	panic("not yet implemented")
}

// --- [ Expressions ] ----------------------------------------------------------

// TODO: Consider merging expr and constExpr, and using type assertion on
// constant.Constant to verify that the expression is constant where needed
// (e.g. initializer of global variable definition).

// expr lowers the given expression to LLVM IR, emitting code to f.
func (m *Module) expr(f *Function, expr ast.Expr) value.Value {
	typ := m.typeOf(expr)
	switch expr := expr.(type) {
	case *ast.BasicLit:
		return m.basicLit(f, expr, typ)
	case *ast.BinaryExpr:
		panic(fmt.Sprintf("support for type %T not yet implemented", expr))
	case *ast.CallExpr:
		panic(fmt.Sprintf("support for type %T not yet implemented", expr))
	case *ast.Ident:
		return m.ident(f, expr, typ)
	case *ast.IndexExpr:
		panic(fmt.Sprintf("support for type %T not yet implemented", expr))
	case *ast.ParenExpr:
		panic(fmt.Sprintf("support for type %T not yet implemented", expr))
	case *ast.UnaryExpr:
		panic(fmt.Sprintf("support for type %T not yet implemented", expr))
	default:
		panic(fmt.Sprintf("support for type %T not yet implemented", expr))
	}
	panic("unreachable")
}

// basicLit lowers the given basic literal to LLVM IR, emitting code to f.
func (m *Module) basicLit(f *Function, n *ast.BasicLit, typ irtypes.Type) value.Value {
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
}

// ident lowers the given identifier to LLVM IR, emitting code to f.
func (m *Module) ident(f *Function, ident *ast.Ident, typ irtypes.Type) value.Value {
	// Input:
	//    void f() {
	//       int x;
	//       x;               // <-- relevant line
	//    }
	// Output:
	//    %1 = load i32, i32* %x
	addr := f.local(ident.String())
	inst, err := instruction.NewLoad(typ, addr)
	if err != nil {
		panic(fmt.Sprintf("unable to create local instruction; %v", err))
	}

	// Emit load instruction.
	return f.emitLocal("", inst)
}

// constExpr converts the given expression to an LLVM IR constant expression.
func (m *Module) constExpr(expr ast.Expr) constant.Constant {
	typ := m.typeOf(expr)
	switch expr := expr.(type) {
	case *ast.BasicLit:
		switch expr.Kind {
		case token.CharLit:
			s, err := strconv.Unquote(expr.Val)
			if err != nil {
				panic(fmt.Sprintf("unable to unquote character literal; %v", err))
			}
			val, err := constant.NewInt(typ, strconv.Itoa(int(s[0])))
			if err != nil {
				panic(fmt.Sprintf("unable to create integer constant; %v", err))
			}
			return val
		case token.IntLit:
			val, err := constant.NewInt(typ, expr.Val)
			if err != nil {
				panic(fmt.Sprintf("unable to create integer constant; %v", err))
			}
			return val
		default:
			panic(fmt.Sprintf("support for basic literal kind %v not yet implemented", expr.Kind))
		}
	//case *ast.BinaryExpr:
	//case *ast.CallExpr:
	//case *ast.Ident:
	//case *ast.IndexExpr:
	//case *ast.ParenExpr:
	//case *ast.UnaryExpr:
	default:
		panic(fmt.Sprintf("support for type %T not yet implemented", expr))
	}
	panic("unreachable")
}
