// Package typecheck implements type-checking of parse trees.
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
			if !isCompatible(resType, curFunc.Result) {
				resPos := n.Start()
				if n.Result != nil {
					resPos = n.Result.Start()
				}
				return errutil.Newf("%d: returning %q from a function with incompatible result type %q", resPos, resType, curFunc.Result)
			}
		case *ast.CallExpr:
			fn, ok := n.Name.Decl.Type().(*types.Func)
			if !ok {
				return errutil.Newf("%d: cannot call non-function %q of type %q", n.Lparen, n.Name, fn)
			}
			// Check number of arguments.
			if len(n.Args) < len(fn.Params) {
				return errutil.Newf("%d: calling %q with too few arguments; expected %d, got %d", n.Lparen, n.Name, len(fn.Params), len(n.Args))
			} else if len(n.Args) > len(fn.Params) {
				return errutil.Newf("%d: calling %q with too many arguments; expected %d, got %d", n.Lparen, n.Name, len(fn.Params), len(n.Args))
			}
			// Check that callee argument types match the function parameter types.
			for i, param := range fn.Params {
				arg := n.Args[i]
				paramType := param.Type
				argType := exprType[arg]
				if !isCompatibleArg(paramType, argType) {
					return errutil.Newf("%d: calling %q with incompatible argument type %q to parameter of type %q", arg.Start(), n.Name, argType, paramType)
				}
			}
			return nil
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

// isCompatibleArg reports whether the given argument and function parameter
// types are compatible.
func isCompatibleArg(param, arg types.Type) bool {
	if isCompatible(param, arg) {
		return true
	}
	if param, ok := param.(*types.Array); ok {
		if arg, ok := arg.(*types.Array); ok {
			// TODO: future. Check for other compatibles; pointers,
			// arraynames, strings(gcc)?
			if param.Len == 0 {
				if paramElem, ok := param.Elem.(*types.Basic); ok {
					if argElem, ok := arg.Elem.(*types.Basic); ok {
						return paramElem.Kind == argElem.Kind
					}
				}
			}
		}
	}
	return false
}

// isCompatible reports whether a and b are of compatible types.
func isCompatible(a, b types.Type) bool {
	if types.Equal(a, b) {
		return true
	}
	if a, ok := a.(types.Numerical); ok {
		if b, ok := b.(types.Numerical); ok {
			// TODO: future. Check for other compatibles; pointers,
			// arraynames, strings(gcc)?
			return a.IsNumerical() && b.IsNumerical()
		}
	}
	return false
}
