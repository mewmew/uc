// Package typecheck implements type-checking of parse trees.
package typecheck

import (
	"github.com/mewkiz/pkg/errutil"
	"github.com/mewmew/uc/ast"
	"github.com/mewmew/uc/ast/astutil"
	"github.com/mewmew/uc/sem/errors"
	"github.com/mewmew/uc/types"
)

// Check type-checks the given file.
func Check(file *ast.File) error {
	// Deduce the types of expressions.
	exprTypes, err := deduce(file)
	if err != nil {
		return errutil.Err(err)
	}

	// Type-check file.
	if err := check(file, exprTypes); err != nil {
		return errutil.Err(err)
	}

	return nil
}

// check type-checks the given file.
func check(file *ast.File, exprTypes map[ast.Expr]types.Type) error {
	// funcs is a stack of function declarations, where the top-most entry
	// represents the currently active function.
	var funcs []*types.Func

	// check type-checks the given node.
	check := func(n ast.Node) error {
		switch n := n.(type) {
		case *ast.FuncDecl:
			if astutil.IsDef(n) {
				// push function declaration.
				funcs = append(funcs, n.Type().(*types.Func))

				// Verify that parameter names are not obmitted in function
				// definitions.
				for _, param := range n.FuncType.Params {
					if !types.IsVoid(param.Type()) && param.VarName == nil {
						return errors.Newf(param.VarType.Start(), "parameter name obmitted")
					}
				}

			}
		case *ast.ReturnStmt:
			// Verify result type.
			curFunc := funcs[len(funcs)-1]
			var resultType types.Type
			resultType = &types.Basic{Kind: types.Void}
			if n.Result != nil {
				resultType = exprTypes[n.Result]
			}
			if !isCompatible(resultType, curFunc.Result) {
				resultPos := n.Start()
				if n.Result != nil {
					resultPos = n.Result.Start()
				}
				return errors.Newf(resultPos, "returning %q from a function with incompatible result type %q", resultType, curFunc.Result)
			}
		case *ast.CallExpr:
			funcType, ok := n.Name.Decl.Type().(*types.Func)
			if !ok {
				return errors.Newf(n.Lparen, "cannot call non-function %q of type %q", n.Name, funcType)
			}
			// TODO: Implement support for functions with variable arguments (i.e.
			// ellipsis).

			// Verify call without arguments.
			if len(n.Args) == 0 && len(funcType.Params) == 1 && types.IsVoid(funcType.Params[0].Type) {
				return nil
			}

			// Check number of arguments.
			if len(n.Args) < len(funcType.Params) {
				return errors.Newf(n.Lparen, "calling %q with too few arguments; expected %d, got %d", n.Name, len(funcType.Params), len(n.Args))
			}
			if len(n.Args) > len(funcType.Params) {
				return errors.Newf(n.Lparen, "calling %q with too many arguments; expected %d, got %d", n.Name, len(funcType.Params), len(n.Args))
			}

			// Check that call argument types match the function parameter types.
			for i, param := range funcType.Params {
				arg := n.Args[i]
				argType := exprTypes[arg]
				paramType := param.Type
				if !isCompatibleArg(argType, paramType) {
					return errors.Newf(arg.Start(), "calling %q with incompatible argument type %q to parameter of type %q", n.Name, argType, paramType)
				}
			}
		case *ast.FuncType:
			for _, param := range n.Params {
				paramType := param.Type()
				if len(n.Params) > 1 && types.IsVoid(paramType) {
					return errors.Newf(n.Lparen, `"void" must be the only parameter`)
				}
			}
		default:
			// TODO: Implement type-checking for remaining node types.
			//log.Printf("not type-checked: %T\n", n)
		}
		return nil
	}

	// after reverts to the outer function after traversing function definitions.
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

	// Walk the AST of the given file to perform type-checking.
	if err := astutil.WalkBeforeAfter(file, check, after); err != nil {
		return errutil.Err(err)
	}

	return nil
}

// isCompatibleArg reports whether the given call argument and function
// parameter types are compatible.
func isCompatibleArg(arg, param types.Type) bool {
	if isCompatible(arg, param) {
		return true
	}
	if arg, ok := arg.(*types.Array); ok {
		if param, ok := param.(*types.Array); ok {
			// TODO: Check for other compatible types (e.g. pointers, array names,
			// strings).
			if param.Len != 0 {
				return false
			}
			return types.Equal(arg.Elem, param.Elem)
		}
	}
	return false
}

// isCompatible reports whether t and u are of compatible types.
func isCompatible(t, u types.Type) bool {
	if types.Equal(t, u) {
		return true
	}
	if t, ok := t.(types.Numerical); ok {
		if u, ok := u.(types.Numerical); ok {
			// TODO: Check for other compatible types (e.g. pointers, array names,
			// strings).
			return t.IsNumerical() && u.IsNumerical()
		}
	}
	return false
}
