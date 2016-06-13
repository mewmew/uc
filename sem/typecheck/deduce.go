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

// deduce performs type deduction of expressions, and store the result in
// exprTypes.
func deduce(file *ast.File, exprTypes map[ast.Expr]types.Type) error {
	// deduce performs type deduction of the given expression.
	deduce := func(n ast.Node) error {
		if expr, ok := n.(ast.Expr); ok {
			typ, err := typeOf(expr)
			if err != nil {
				return errutil.Err(err)
			}
			exprTypes[expr] = typ
		}
		return nil
	}

	// Walk the AST of the given file to deduce the types of expression nodes.
	if err := astutil.Walk(file, deduce); err != nil {
		return errutil.Err(err)
	}

	return nil
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
			panic(fmt.Sprintf("support for basic literal type %v not yet implemented", n.Kind))
		}
	case *ast.BinaryExpr:
		xType, err := typeOf(n.X)
		if err != nil {
			return nil, errutil.Err(err)
		}
		yType, err := typeOf(n.Y)
		if err != nil {
			return nil, errutil.Err(err)
		}
		if n.Op == token.Assign {
			if !isAssignable(n.X) {
				return nil, errors.Newf(n.OpPos, "cannot assign to %q of type %q", n.X, xType)
			}
			if !isCompatible(xType, yType) {
				return nil, errors.Newf(n.OpPos, "cannot assign to %q (type mismatch between %q and %q)", n.X, xType, yType)
			}
			// TODO: !types.Equal(higherPrecision(xType, yType), xType) could be
			// used for loss of percision warning.
			return xType, nil
		}
		if types.IsVoid(xType) || types.IsVoid(yType) {
			return nil, errors.Newf(n.OpPos, "invalid operands to binary expression: %v (%q and %q)", n, xType, yType)
		}
		if !isCompatible(xType, yType) {
			return nil, errors.Newf(n.OpPos, "invalid operation: %v (type mismatch between %q and %q)", n, xType, yType)
		}
		// TODO: Implement better implicit conversion. Future: Make sure to
		// promote types early when implementing signed/unsigned types and
		// types need to be promoted anyway later. Be careful of bug:
		// https://youtu.be/Ux0YnVEaI6A?t=279
		return higherPrecision(xType, yType), nil
	case *ast.CallExpr:
		typ := n.Name.Decl.Type()
		if typ, ok := typ.(*types.Func); ok {
			return typ.Result, nil
		}
		return nil, errors.Newf(n.Lparen, "cannot call non-function %q of type %q", n.Name, typ)
	case *ast.Ident:
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
		// TODO: Add support for pointers.
		return typeOf(n.X)
	default:
		panic(fmt.Sprintf("support for type %T not yet implemented.", n))
	}
}

// TODO: Verify isAssignable against the definition of lvalue in the C spec (I
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

// higherPrecision returns the type of higher precision.
func higherPrecision(t, u types.Type) types.Type {
	// TODO: Implement with a list of types sorted by precision when support for
	// more types are added.
	if t, ok := t.(*types.Basic); ok {
		if u, ok := u.(*types.Basic); ok {
			if t.Kind == types.Void || u.Kind == types.Void {
				panic(fmt.Sprint(`incorrect use of higherPrecision; "void" does not have precision.`))
			}
			// Check for types in order of highest precision.
			if t.Kind == types.Int || u.Kind == types.Int {
				return &types.Basic{Kind: types.Int}
			}
			if t.Kind == types.Char || u.Kind == types.Char {
				return &types.Basic{Kind: types.Char}
			}
		}
	}
	panic(fmt.Sprintf(`support for type "%v" or "%v" not yet implemented.`, t, u))
}
