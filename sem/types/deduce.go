package types

import (
	"fmt"

	"github.com/mewkiz/pkg/errutil"
	"github.com/mewmew/uc/ast"
	"github.com/mewmew/uc/ast/astutil"
	"github.com/mewmew/uc/token"
	"github.com/mewmew/uc/types"
)

// deduce performs type deduction of expressions to annotate the AST.
func deduce(file *ast.File) error {

	// funcs is a stack of function declarations, where the top-most entry
	// represents the currently active function.
	var funcs []*types.Func

	before := func(n ast.Node) error {
		switch n := n.(type) {
		case *ast.FuncDecl:
			if isDef(n) {
				// push function declaration.
				funcs = append(funcs, n.Type().(*types.Func))
			}
		case *ast.ReturnStmt:
			curFunc := funcs[len(funcs)-1]
			resType, err := typeOf(n.Result)
			if err != nil {
				return errutil.Err(err)
			}
			if !compatible(curFunc.Result, resType) {
				return errutil.Newf("%d: returning %q from a function with incompatible result type %q", n.Result.Start(), resType, curFunc.Result)
			}

		}
		return nil
	}

	after := func(n ast.Node) error {
		switch n := n.(type) {
		case *ast.FuncDecl:
			if isDef(n) {
				// pop function declaration.
				funcs = funcs[:len(funcs)-1]
			}
		}
		return nil
	}

	if err := astutil.WalkBeforeAfter(file, before, after); err != nil {
		return errutil.Err(err)
	}

	return nil
}

func compatible(a, b types.Type) bool {
	// TODO: Implement type compatibility checks.
	return types.Equal(a, b)
}

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
		if !compatible(x, y) {
			return nil, errutil.Newf("invalid operands to binary expression (%q and %q)", x, y)
		}
		// TODO: Implement implicit conversion.
		return x, nil
		panic(fmt.Sprintf("support for type %T not yet implemented.", n))
	case *ast.CallExpr:
		panic(fmt.Sprintf("support for type %T not yet implemented.", n))
	case *ast.Ident:
		// TODO: Make sure that type declarations are handled correctly for
		// keyword types such as "int".
		return n.Decl.Type(), nil
	case *ast.IndexExpr:
		panic(fmt.Sprintf("support for type %T not yet implemented.", n))
	case *ast.ParenExpr:
		panic(fmt.Sprintf("support for type %T not yet implemented.", n))
	case *ast.UnaryExpr:
		panic(fmt.Sprintf("support for type %T not yet implemented.", n))
	default:
		panic(fmt.Sprintf("support for type %T not yet implemented.", n))
	}
}
