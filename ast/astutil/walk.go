// Package astutil implements utility functions for walking parse trees.
package astutil

import (
	"fmt"

	"github.com/mewkiz/pkg/errutil"
	"github.com/mewmew/uc/ast"
)

// Walk traverses the given parse tree, calling f(n) for each node n in the
// tree, in a bottom-up traversal.
func Walk(node ast.Node, f func(ast.Node) error) error {
	nop := func(n ast.Node) error { return nil }
	return WalkBeforeAfter(node, nop, f)
}

// WalkBeforeAfter traverses the given parse tree, calling before(n) before
// traversing the node's children, and after(n) afterwards, in a bottom-up
// traversal.
func WalkBeforeAfter(node ast.Node, before, after func(ast.Node) error) error {
	switch n := node.(type) {

	// Source file.
	case *ast.File:
		if n != nil {
			return walkFile(n, before, after)
		}

	// Declarations.
	case *ast.FuncDecl:
		if n != nil {
			return walkFuncDecl(n, before, after)
		}
	case *ast.VarDecl:
		if n != nil {
			return walkVarDecl(n, before, after)
		}

	// Statements.
	case *ast.BlockStmt:
		if n != nil {
			return walkBlockStmt(n, before, after)
		}
	case *ast.EmptyStmt:
		if n != nil {
			return walkEmptyStmt(n, before, after)
		}
	case *ast.ExprStmt:
		if n != nil {
			return walkExprStmt(n, before, after)
		}
	case *ast.IfStmt:
		if n != nil {
			return walkIfStmt(n, before, after)
		}
	case *ast.ReturnStmt:
		if n != nil {
			return walkReturnStmt(n, before, after)
		}
	case *ast.WhileStmt:
		if n != nil {
			return walkWhileStmt(n, before, after)
		}

	// Expressions.
	case *ast.BasicLit:
		if n != nil {
			return walkBasicLit(n, before, after)
		}
	case *ast.BinaryExpr:
		if n != nil {
			return walkBinaryExpr(n, before, after)
		}
	case *ast.CallExpr:
		if n != nil {
			return walkCallExpr(n, before, after)
		}
	case *ast.Ident:
		if n != nil {
			return walkIdent(n, before, after)
		}
	case *ast.IndexExpr:
		if n != nil {
			return walkIndexExpr(n, before, after)
		}
	case *ast.ParenExpr:
		if n != nil {
			return walkParenExpr(n, before, after)
		}
	case *ast.UnaryExpr:
		if n != nil {
			return walkUnaryExpr(n, before, after)
		}

	// Types.
	case *ast.ArrayType:
		if n != nil {
			return walkArrayType(n, before, after)
		}
	case *ast.FuncType:
		if n != nil {
			return walkFuncType(n, before, after)
		}
	case *ast.Field:
		if n != nil {
			return walkTypeField(n, before, after)
		}

	case nil:
		// Nothing to do.
		return nil
	default:
		panic(fmt.Sprintf("support for walking node of type %T not yet implemented", node))
	}

	return nil
}

// === [ Source file ] ===

// walkFile walks the parse tree of the given source file in depth first order.
func walkFile(file *ast.File, before, after func(ast.Node) error) error {
	if err := before(file); err != nil {
		return errutil.Err(err)
	}
	for _, decl := range file.Decls {
		if err := WalkBeforeAfter(decl, before, after); err != nil {
			return errutil.Err(err)
		}
	}
	if err := after(file); err != nil {
		return errutil.Err(err)
	}
	return nil
}

// === [ Top-level declarations ] ===

// walkFuncDecl walks the parse tree of the given function declaration in depth
// first order.
func walkFuncDecl(decl *ast.FuncDecl, before, after func(ast.Node) error) error {
	if err := before(decl); err != nil {
		return errutil.Err(err)
	}
	if err := WalkBeforeAfter(decl.FuncName, before, after); err != nil {
		return errutil.Err(err)
	}
	if err := WalkBeforeAfter(decl.FuncType, before, after); err != nil {
		return errutil.Err(err)
	}
	if err := WalkBeforeAfter(decl.Body, before, after); err != nil {
		return errutil.Err(err)
	}
	if err := after(decl); err != nil {
		return errutil.Err(err)
	}
	return nil
}

// walkVarDecl walks the parse tree of the given variable declaration in depth
// first order.
func walkVarDecl(decl *ast.VarDecl, before, after func(ast.Node) error) error {
	if err := before(decl); err != nil {
		return errutil.Err(err)
	}
	if err := WalkBeforeAfter(decl.VarType, before, after); err != nil {
		return errutil.Err(err)
	}
	if err := WalkBeforeAfter(decl.VarName, before, after); err != nil {
		return errutil.Err(err)
	}
	if err := WalkBeforeAfter(decl.Val, before, after); err != nil {
		return errutil.Err(err)
	}
	if err := after(decl); err != nil {
		return errutil.Err(err)
	}
	return nil
}

// === [ Statements ] ===

// walkBlockStmt walks the parse tree of the given block statement in depth
// first order.
func walkBlockStmt(block *ast.BlockStmt, before, after func(ast.Node) error) error {
	if err := before(block); err != nil {
		return errutil.Err(err)
	}
	for _, item := range block.Items {
		if err := WalkBeforeAfter(item, before, after); err != nil {
			return errutil.Err(err)
		}
	}
	if err := after(block); err != nil {
		return errutil.Err(err)
	}
	return nil
}

// walkEmptyStmt walks the parse tree of the given empty statement in depth
// first order.
func walkEmptyStmt(stmt *ast.EmptyStmt, before, after func(ast.Node) error) error {
	if err := before(stmt); err != nil {
		return errutil.Err(err)
	}
	if err := after(stmt); err != nil {
		return errutil.Err(err)
	}
	return nil
}

// walkExprStmt walks the parse tree of the given expression statement in depth
// first order.
func walkExprStmt(stmt *ast.ExprStmt, before, after func(ast.Node) error) error {
	if err := before(stmt); err != nil {
		return errutil.Err(err)
	}
	if err := WalkBeforeAfter(stmt.X, before, after); err != nil {
		return errutil.Err(err)
	}
	if err := after(stmt); err != nil {
		return errutil.Err(err)
	}
	return nil
}

// walkIfStmt walks the parse tree of the given if statement in depth first
// order.
func walkIfStmt(stmt *ast.IfStmt, before, after func(ast.Node) error) error {
	if err := before(stmt); err != nil {
		return errutil.Err(err)
	}
	if err := WalkBeforeAfter(stmt.Cond, before, after); err != nil {
		return errutil.Err(err)
	}
	if err := WalkBeforeAfter(stmt.Body, before, after); err != nil {
		return errutil.Err(err)
	}
	if err := WalkBeforeAfter(stmt.Else, before, after); err != nil {
		return errutil.Err(err)
	}
	if err := after(stmt); err != nil {
		return errutil.Err(err)
	}
	return nil
}

// walkReturnStmt walks the parse tree of the given return statement in depth
// first order.
func walkReturnStmt(stmt *ast.ReturnStmt, before, after func(ast.Node) error) error {
	if err := before(stmt); err != nil {
		return errutil.Err(err)
	}
	if err := WalkBeforeAfter(stmt.Result, before, after); err != nil {
		return errutil.Err(err)
	}
	if err := after(stmt); err != nil {
		return errutil.Err(err)
	}
	return nil
}

// walkWhileStmt walks the parse tree of the given while statement in depth
// first order.
func walkWhileStmt(stmt *ast.WhileStmt, before, after func(ast.Node) error) error {
	if err := before(stmt); err != nil {
		return errutil.Err(err)
	}
	if err := WalkBeforeAfter(stmt.Cond, before, after); err != nil {
		return errutil.Err(err)
	}
	if err := WalkBeforeAfter(stmt.Body, before, after); err != nil {
		return errutil.Err(err)
	}
	if err := after(stmt); err != nil {
		return errutil.Err(err)
	}
	return nil
}

// === [ Expressions ] ===

// walkBasicLit walks the parse tree of the given basic literal expression in
// depth first order.
func walkBasicLit(lit *ast.BasicLit, before, after func(ast.Node) error) error {
	if err := before(lit); err != nil {
		return errutil.Err(err)
	}
	if err := after(lit); err != nil {
		return errutil.Err(err)
	}
	return nil
}

// walkBinaryExpr walks the parse tree of the given binary expression in depth
// first order.
func walkBinaryExpr(expr *ast.BinaryExpr, before, after func(ast.Node) error) error {
	if err := before(expr); err != nil {
		return errutil.Err(err)
	}
	if err := WalkBeforeAfter(expr.X, before, after); err != nil {
		return errutil.Err(err)
	}
	if err := WalkBeforeAfter(expr.Y, before, after); err != nil {
		return errutil.Err(err)
	}
	if err := after(expr); err != nil {
		return errutil.Err(err)
	}
	return nil
}

// walkCallExpr walks the parse tree of the given call expression in depth first
// order.
func walkCallExpr(call *ast.CallExpr, before, after func(ast.Node) error) error {
	if err := before(call); err != nil {
		return errutil.Err(err)
	}
	if err := WalkBeforeAfter(call.Name, before, after); err != nil {
		return errutil.Err(err)
	}
	for _, arg := range call.Args {
		if err := WalkBeforeAfter(arg, before, after); err != nil {
			return errutil.Err(err)
		}
	}
	if err := after(call); err != nil {
		return errutil.Err(err)
	}
	return nil
}

// walkIdent walks the parse tree of the given identifier expression in depth
// first order.
func walkIdent(ident *ast.Ident, before, after func(ast.Node) error) error {
	if err := before(ident); err != nil {
		return errutil.Err(err)
	}
	if err := after(ident); err != nil {
		return errutil.Err(err)
	}
	return nil
}

// walkIndexExpr walks the parse tree of the given index expression in depth
// first order.
func walkIndexExpr(expr *ast.IndexExpr, before, after func(ast.Node) error) error {
	if err := before(expr); err != nil {
		return errutil.Err(err)
	}
	if err := WalkBeforeAfter(expr.Name, before, after); err != nil {
		return errutil.Err(err)
	}
	if err := WalkBeforeAfter(expr.Index, before, after); err != nil {
		return errutil.Err(err)
	}
	if err := after(expr); err != nil {
		return errutil.Err(err)
	}
	return nil
}

// walkParenExpr walks the parse tree of the given parenthesized expression in
// depth first order.
func walkParenExpr(expr *ast.ParenExpr, before, after func(ast.Node) error) error {
	if err := before(expr); err != nil {
		return errutil.Err(err)
	}
	if err := WalkBeforeAfter(expr.X, before, after); err != nil {
		return errutil.Err(err)
	}
	if err := after(expr); err != nil {
		return errutil.Err(err)
	}
	return nil
}

// walkUnaryExpr walks the parse tree of the given unary expression in depth
// first order.
func walkUnaryExpr(expr *ast.UnaryExpr, before, after func(ast.Node) error) error {
	if err := before(expr); err != nil {
		return errutil.Err(err)
	}
	if err := WalkBeforeAfter(expr.X, before, after); err != nil {
		return errutil.Err(err)
	}
	if err := after(expr); err != nil {
		return errutil.Err(err)
	}
	return nil
}

// === [ Types ] ===

// walkArrayType walks the parse tree of the given array type in depth first
// order.
func walkArrayType(arr *ast.ArrayType, before, after func(ast.Node) error) error {
	if err := before(arr); err != nil {
		return errutil.Err(err)
	}
	if err := WalkBeforeAfter(arr.Elem, before, after); err != nil {
		return errutil.Err(err)
	}
	if err := after(arr); err != nil {
		return errutil.Err(err)
	}
	return nil
}

// walkFuncType walks the parse tree of the given function signature in depth
// first order.
func walkFuncType(fn *ast.FuncType, before, after func(ast.Node) error) error {
	if err := before(fn); err != nil {
		return errutil.Err(err)
	}
	if err := WalkBeforeAfter(fn.Result, before, after); err != nil {
		return errutil.Err(err)
	}
	for _, param := range fn.Params {
		if err := WalkBeforeAfter(param, before, after); err != nil {
			return errutil.Err(err)
		}
	}
	if err := after(fn); err != nil {
		return errutil.Err(err)
	}
	return nil
}

// walkTypeField walks the parse tree of the given type field in depth first
// order.
func walkTypeField(field *ast.Field, before, after func(ast.Node) error) error {
	if err := before(field); err != nil {
		return errutil.Err(err)
	}
	if err := WalkBeforeAfter(field.Type, before, after); err != nil {
		return errutil.Err(err)
	}
	if err := WalkBeforeAfter(field.Name, before, after); err != nil {
		return errutil.Err(err)
	}
	if err := after(field); err != nil {
		return errutil.Err(err)
	}
	return nil
}
