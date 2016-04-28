package typecheck

import (
	"fmt"

	"github.com/mewkiz/pkg/errutil"
	"github.com/mewmew/uc/ast"
	"github.com/mewmew/uc/ast/astutil"
	"github.com/mewmew/uc/token"
	"github.com/mewmew/uc/types"
)

// deduce performs type deduction of expressions to annotate the AST.
func deduce(file *ast.File) (exprType map[ast.Expr]types.Type, err error) {
	// Map expression nodes to types.
	exprType = make(map[ast.Expr]types.Type)
	deduce := func(n ast.Node) error {
		if expr, ok := n.(ast.Expr); ok {
			typ, err := typeOf(expr)
			if err != nil {
				return errutil.Err(err)
			}
			exprType[expr] = typ
		}
		return nil
	}
	if err := astutil.Walk(file, deduce); err != nil {
		return nil, errutil.Err(err)
	}
	return exprType, nil
}

// typeOf returns the type of the given expression.
func typeOf(n ast.Expr) (types.Type, error) {
	switch n := n.(type) {
	case *ast.BasicLit:
		switch n.Kind {
		case token.CharLit:
			return &types.Basic{Kind: types.Char}, nil
		case token.IntLit:
			return &types.Basic{Kind: types.Int}, nil
		default:
			panic(fmt.Sprintf("support for basic type kind %v not yet implemented", n.Kind))
		}
		panic(fmt.Sprintf("support for type %T not yet implemented.", n))
	case *ast.BinaryExpr:
		x, err := typeOf(n.X)
		if err != nil {
			return nil, errutil.Err(err)
		}
		y, err := typeOf(n.Y)
		if err != nil {
			return nil, errutil.Err(err)
		}
		if !isCompatible(x, y) {
			return nil, errutil.Newf("invalid operation: %v (type mismatch between %q and %q)", n, x, y)
		}
		// TODO: Implement implicit conversion.
		return x, nil
		panic(fmt.Sprintf("support for type %T not yet implemented.", n))
	case *ast.CallExpr:
		typ := n.Name.Decl.Type()
		if typ, ok := typ.(*types.Func); ok {
			return typ.Result, nil
		}
		return nil, errutil.Newf("%d: cannot call non-function %q of type %q", n.Start(), n.Name, typ)
	case *ast.Ident:
		// TODO: Make sure that type declarations are handled correctly for
		// keyword types such as "int".
		return n.Decl.Type(), nil
	case *ast.IndexExpr:
		typ := n.Name.Decl.Type()
		if typ, ok := typ.(*types.Array); ok {
			return typ.Elem, nil
		}
		return nil, errutil.Newf("%d: invalid operation: %v (type %q does not support indexing)", n.Lbracket, n, typ)
	case *ast.ParenExpr:
		return typeOf(n.X)
	case *ast.UnaryExpr:
		panic(fmt.Sprintf("support for type %T not yet implemented.", n))
	default:
		panic(fmt.Sprintf("support for type %T not yet implemented.", n))
	}
}
