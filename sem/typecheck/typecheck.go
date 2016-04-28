package typecheck

import (
	"github.com/mewkiz/pkg/errutil"
	"github.com/mewmew/uc/ast"
	"github.com/mewmew/uc/ast/astutil"
	"github.com/mewmew/uc/types"
)

// Check type-checks the given file.
func Check(file *ast.File) error {
	// Deduce the types of expressions.
	exprType, err := deduce(file)
	if err != nil {
		return errutil.Err(err)
	}

	// Type-check file.
	if err := check(file, exprType); err != nil {
		return errutil.Err(err)
	}

	return nil
}

// check type-checks the given file.
func check(file *ast.File, exprType map[ast.Expr]types.Type) error {
	// funcs is a stack of function declarations, where the top-most entry
	// represents the currently active function.
	var funcs []*types.Func

	before := func(n ast.Node) error {
		switch n := n.(type) {
		case *ast.FuncDecl:
			if astutil.IsDef(n) {
				// push function declaration.
				funcs = append(funcs, n.Type().(*types.Func))
			}
		case *ast.ReturnStmt:
			curFunc := funcs[len(funcs)-1]
			var resType types.Type
			resType = &types.Basic{Kind: types.Void}
			if n.Result != nil {
				resType = exprType[n.Result]
			}
			if !compatible(resType, curFunc.Result) {
				resPos := n.Start()
				if n.Result != nil {
					resPos = n.Result.Start()
				}
				return errutil.Newf("%d: returning %q from a function with incompatible result type %q", resPos, resType, curFunc.Result)
			}
		default:
			// TODO: Implement type-checking for remaining node types.
			//log.Printf("not type-checked: %T\n", n)
		}
		return nil
	}

	after := func(n ast.Node) error {
		switch n := n.(type) {
		case *ast.FuncDecl:
			if astutil.IsDef(n) {
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
	if types.Equal(a, b) {
		return true
	}
	if a, ok := a.(types.Numerical); ok {
		if b, ok := b.(types.Numerical); ok {
			// TODO: Check for other compatibles; pointers, arraynames, strings(gcc)?
			return a.IsNumerical() && b.IsNumerical()
		}
	}
	return false
}
