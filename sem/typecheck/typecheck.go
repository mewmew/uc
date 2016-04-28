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
			if decl, ok := n.Name.Decl.Type().(*types.Func); ok {
				// Check number of args
				if len(decl.Params) != len(n.Args) {
					return errutil.Newf("%d: calling %q with wrong number of arguments %q", n.Start(), n, decl..String())
				}
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
		case *ast.CallExpr:
			decl, ok := n.Name.Decl.Type().(*types.Func)
			if !ok {
				return errutil.Newf("%d: calling %q from a function with incompatible type %q", n.Start(), n, n.Name.Decl)
			}
			declParams := decl.Params
			callArgs := n.Args
			// Check that callee args types match the declared param types
			for i, declParam := range declParams {
				callArg := callArgs[i]
				if isCallCompatible(declParam.Type, exprType[callArg]) {
					return nil
				}
				return errutil.Newf("%d: calling %q with incompatible type %q instead of %q", callArgs[i].Start(), n, callArgs[i], declParam)
			}
		}
		return nil
	}

	if err := astutil.WalkBeforeAfter(file, before, after); err != nil {
		return errutil.Err(err)
	}

	return nil
}

func isCallCompatible(decl, call types.Type) bool {
	if isCompatible(decl, call) {
		return true
	}
	if decl, ok := decl.(*types.Array); ok {
		if call, ok := call.(*types.Array); ok {
			// TODO: future. Check for other compatibles; pointers,
			// arraynames, strings(gcc)?
			if decl.Len == 0 {
				if declArrayType, ok := decl.Elem.(*types.Basic); ok {
					if callArrayType, ok := call.Elem.(*types.Basic); ok {
						return declArrayType.Kind == callArrayType.Kind
					}
				}
			}
		}
	}
	return false
}
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
