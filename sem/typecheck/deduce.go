package typecheck

import (
	"fmt"

	"github.com/mewkiz/pkg/errutil"
	"github.com/mewmew/uc/ast"
	"github.com/mewmew/uc/ast/astutil"
	"github.com/mewmew/uc/sem/errors"
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
	case *ast.BinaryExpr:
		x, err := typeOf(n.X)
		if err != nil {
			return nil, errutil.Err(err)
		}
		y, err := typeOf(n.Y)
		if err != nil {
			return nil, errutil.Err(err)
		}
		if n.Op == token.Assign {
			if !isAssignable(n.X) {
				return nil, errors.Newf(n.OpPos, "cannot assign to %q of type %q", n.X, x)
			}
			if !isCompatible(x, y) {
				return nil, errors.Newf(n.OpPos, "cannot assign to %q (type mismatch between %q and %q)", n.X, x, y)
			}
			// TODO: higherPrecision(x,y) != x could be used for loss of
			// percision warning.
			return x, nil
		} else if isVoid(x) || isVoid(y) {
			return nil, errors.Newf(n.OpPos, "invalid operands to binary expression: %v (%q and %q)", n, x, y)
		} else if !isCompatible(x, y) {
			return nil, errors.Newf(n.OpPos, "invalid operation: %v (type mismatch between %q and %q)", n, x, y)
		}
		// TODO: Implement better implicit conversion.
		return higherPrecision(x, y), nil
	case *ast.CallExpr:
		typ := n.Name.Decl.Type()
		if typ, ok := typ.(*types.Func); ok {
			return typ.Result, nil
		}
		return nil, errors.Newf(n.Lparen, "cannot call non-function %q of type %q", n.Name, typ)
	case *ast.Ident:
		// TODO: Make sure that type declarations are handled correctly for
		// keyword types such as "int".
		return n.Decl.Type(), nil
	case *ast.IndexExpr:
		typ := n.Name.Decl.Type()
		if typ, ok := typ.(*types.Array); ok {
			return typ.Elem, nil
		}
		return nil, errors.Newf(n.Lbracket, "invalid operation: %v (type %q does not support indexing)", n, typ)
	case *ast.ParenExpr:
		return typeOf(n.X)
	case *ast.UnaryExpr:
		// TODO: Future. Add support for pointer
		return typeOf(n.X)
	default:
		panic(fmt.Sprintf("support for type %T not yet implemented.", n))
	}
}

// TODO: Verify isAssignable against the definition of lvale in the C spec (I
// tried and failed).

// isAssignable reports whether the given expression is assignable (i.e. a valid
// lvalue).
func isAssignable(x ast.Expr) bool {
	switch x := x.(type) {
	case *ast.BasicLit:
		return false
	case *ast.BinaryExpr:
		return false
	case *ast.CallExpr:
		return false
	case *ast.Ident:
		switch typ := x.Decl.Type().(type) {
		case *types.Basic:
			return true
		case *types.Array:
			return false
		case *types.Func:
			return false
		default:
			panic(fmt.Sprintf("support for declaration type %T not yet implemented", typ))
		}
	case *ast.IndexExpr:
		return true
	case *ast.ParenExpr:
		return isAssignable(x.X)
	case *ast.UnaryExpr:
		// TODO: Add support for pointer dereferences.
		return false
	default:
		panic(fmt.Sprintf("support for expression type %T not yet implemented", x))
	}
}

// isVoid reports whether the given type is a void type.
func isVoid(typ types.Type) bool {
	if typ, ok := typ.(*types.Basic); ok {
		return typ.Kind == types.Void
	}
	return false
}

// higherPrecision returns the type of higher precision.
func higherPrecision(x, y types.Type) types.Type {
	// TODO: Implement with a list of types sorted by precision when more
	// types are added.
	if x, ok := x.(*types.Basic); ok {
		if y, ok := y.(*types.Basic); ok {
			if x.Kind == types.Void || y.Kind == types.Void {
				panic(fmt.Sprint(`incorrect use of higherPrecision; "void" does not have precision.`))
			}
			// Check for types in order of highest precision
			if x.Kind == types.Int || y.Kind == types.Int {
				return &types.Basic{Kind: types.Int}
			}
			if x.Kind == types.Char || y.Kind == types.Char {
				return &types.Basic{Kind: types.Char}
			}
		}
	}
	panic(fmt.Sprintf(`support for type "%v" or "%v" not yet implemented.`, x, y))
}
