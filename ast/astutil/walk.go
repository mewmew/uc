// Package astutil implements utility functions for walking parse trees.
package astutil

import (
	"fmt"

	"github.com/mewkiz/pkg/errutil"
	"github.com/mewmew/uc/ast"
)

// Walk walks the given parse tree in depth first order.
func Walk(node ast.Node, f func(ast.Node)) error {
	switch n := node.(type) {

	// Source file.
	case *ast.File:
		if n != nil {
			return walkFile(n, f)
		}

	// Declarations.
	case *ast.FuncDecl:
		if n != nil {
			return walkFuncDecl(n, f)
		}
	case *ast.VarDecl:
		if n != nil {
			return walkVarDecl(n, f)
		}
	// Statements.
	case *ast.BlockStmt:
		if n != nil {
			return walkBlockStmt(n, f)
		}
	case *ast.EmptyStmt:
		if n != nil {
			return walkEmptyStmt(n, f)
		}
	case *ast.ExprStmt:
		if n != nil {
			return walkExprStmt(n, f)
		}
	case *ast.IfStmt:
		if n != nil {
			return walkIfStmt(n, f)
		}
	case *ast.ReturnStmt:
		if n != nil {
			return walkReturnStmt(n, f)
		}
	case *ast.WhileStmt:
		if n != nil {
			return walkWhileStmt(n, f)
		}

	// Expressions.
	case *ast.BasicLit:
		if n != nil {
			return walkBasicLit(n, f)
		}
	case *ast.BinaryExpr:
		if n != nil {
			return walkBinaryExpr(n, f)
		}
	case *ast.CallExpr:
		if n != nil {
			return walkCallExpr(n, f)
		}
	case *ast.Ident:
		if n != nil {
			return walkIdent(n, f)
		}
	case *ast.IndexExpr:
		if n != nil {
			return walkIndexExpr(n, f)
		}
	case *ast.ParenExpr:
		if n != nil {
			return walkParenExpr(n, f)
		}
	case *ast.UnaryExpr:
		if n != nil {
			return walkUnaryExpr(n, f)
		}

	// Types.
	case *ast.ArrayType:
		if n != nil {
			return walkArrayType(n, f)
		}
	case *ast.FuncType:
		if n != nil {
			return walkFuncType(n, f)
		}
	case *ast.Field:
		if n != nil {
			return walkTypeField(n, f)
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
func walkFile(file *ast.File, f func(ast.Node)) error {
	f(file)
	for _, decl := range file.Decls {
		if err := Walk(decl, f); err != nil {
			return errutil.Err(err)
		}
	}
	return nil
}

// === [ Top-level declarations ] ===

// walkFuncDecl walks the parse tree of the given function declaration in depth
// first order.
func walkFuncDecl(decl *ast.FuncDecl, f func(ast.Node)) error {
	f(decl)
	if err := Walk(decl.FuncName, f); err != nil {
		return errutil.Err(err)
	}
	if err := Walk(decl.FuncType, f); err != nil {
		return errutil.Err(err)
	}
	if err := Walk(decl.Body, f); err != nil {
		return errutil.Err(err)
	}
	return nil
}

// walkVarDecl walks the parse tree of the given variable declaration in depth
// first order.
func walkVarDecl(decl *ast.VarDecl, f func(ast.Node)) error {
	f(decl)
	if err := Walk(decl.VarType, f); err != nil {
		return errutil.Err(err)
	}
	if err := Walk(decl.VarName, f); err != nil {
		return errutil.Err(err)
	}
	if err := Walk(decl.Val, f); err != nil {
		return errutil.Err(err)
	}
	return nil
}

// === [ Statements ] ===

// walkBlockStmt walks the parse tree of the given block statement in depth
// first order.
func walkBlockStmt(block *ast.BlockStmt, f func(ast.Node)) error {
	f(block)
	for _, item := range block.Items {
		if err := Walk(item, f); err != nil {
			return errutil.Err(err)
		}
	}
	return nil
}

// walkEmptyStmt walks the parse tree of the given empty statement in depth
// first order.
func walkEmptyStmt(stmt *ast.EmptyStmt, f func(ast.Node)) error {
	f(stmt)
	return nil
}

// walkExprStmt walks the parse tree of the given expression statement in depth
// first order.
func walkExprStmt(stmt *ast.ExprStmt, f func(ast.Node)) error {
	f(stmt)
	if err := Walk(stmt.X, f); err != nil {
		return errutil.Err(err)
	}
	return nil
}

// walkIfStmt walks the parse tree of the given if statement in depth first
// order.
func walkIfStmt(stmt *ast.IfStmt, f func(ast.Node)) error {
	f(stmt)
	if err := Walk(stmt.Cond, f); err != nil {
		return errutil.Err(err)
	}
	if err := Walk(stmt.Body, f); err != nil {
		return errutil.Err(err)
	}
	if err := Walk(stmt.Else, f); err != nil {
		return errutil.Err(err)
	}
	return nil
}

// walkReturnStmt walks the parse tree of the given return statement in depth
// first order.
func walkReturnStmt(stmt *ast.ReturnStmt, f func(ast.Node)) error {
	f(stmt)
	if err := Walk(stmt.Result, f); err != nil {
		return errutil.Err(err)
	}
	return nil
}

// walkWhileStmt walks the parse tree of the given while statement in depth
// first order.
func walkWhileStmt(stmt *ast.WhileStmt, f func(ast.Node)) error {
	f(stmt)
	if err := Walk(stmt.Cond, f); err != nil {
		return errutil.Err(err)
	}
	if err := Walk(stmt.Body, f); err != nil {
		return errutil.Err(err)
	}
	return nil
}

// === [ Expressions ] ===

// walkBasicLit walks the parse tree of the given basic literal expression in
// depth first order.
func walkBasicLit(lit *ast.BasicLit, f func(ast.Node)) error {
	f(lit)
	return nil
}

// walkBinaryExpr walks the parse tree of the given binary expression in depth
// first order.
func walkBinaryExpr(expr *ast.BinaryExpr, f func(ast.Node)) error {
	f(expr)
	if err := Walk(expr.X, f); err != nil {
		return errutil.Err(err)
	}
	if err := Walk(expr.Y, f); err != nil {
		return errutil.Err(err)
	}
	return nil
}

// walkCallExpr walks the parse tree of the given call expression in depth first
// order.
func walkCallExpr(call *ast.CallExpr, f func(ast.Node)) error {
	f(call)
	if err := Walk(call.Name, f); err != nil {
		return errutil.Err(err)
	}
	for _, arg := range call.Args {
		if err := Walk(arg, f); err != nil {
			return errutil.Err(err)
		}
	}
	return nil
}

// walkIdent walks the parse tree of the given identifier expression in depth
// first order.
func walkIdent(ident *ast.Ident, f func(ast.Node)) error {
	f(ident)
	return nil
}

// walkIndexExpr walks the parse tree of the given index expression in depth
// first order.
func walkIndexExpr(expr *ast.IndexExpr, f func(ast.Node)) error {
	f(expr)
	if err := Walk(expr.Name, f); err != nil {
		return errutil.Err(err)
	}
	if err := Walk(expr.Index, f); err != nil {
		return errutil.Err(err)
	}
	return nil
}

// walkParenExpr walks the parse tree of the given parenthesized expression in
// depth first order.
func walkParenExpr(expr *ast.ParenExpr, f func(ast.Node)) error {
	f(expr)
	if err := Walk(expr.X, f); err != nil {
		return errutil.Err(err)
	}
	return nil
}

// walkUnaryExpr walks the parse tree of the given unary expression in depth
// first order.
func walkUnaryExpr(expr *ast.UnaryExpr, f func(ast.Node)) error {
	f(expr)
	if err := Walk(expr.X, f); err != nil {
		return errutil.Err(err)
	}
	return nil
}

// === [ Types ] ===

// walkArrayType walks the parse tree of the given array type in depth first
// order.
func walkArrayType(arr *ast.ArrayType, f func(ast.Node)) error {
	f(arr)
	if err := Walk(arr.Elem, f); err != nil {
		return errutil.Err(err)
	}
	return nil
}

// walkFuncType walks the parse tree of the given function signature in depth
// first order.
func walkFuncType(fn *ast.FuncType, f func(ast.Node)) error {
	f(fn)
	if err := Walk(fn.Result, f); err != nil {
		return errutil.Err(err)
	}
	for _, param := range fn.Params {
		if err := Walk(param, f); err != nil {
			return errutil.Err(err)
		}
	}
	return nil
}

// walkTypeField walks the parse tree of the given type field in depth first
// order.
func walkTypeField(field *ast.Field, f func(ast.Node)) error {
	f(field)
	if err := Walk(field.Type, f); err != nil {
		return errutil.Err(err)
	}
	if err := Walk(field.Name, f); err != nil {
		return errutil.Err(err)
	}
	return nil
}
