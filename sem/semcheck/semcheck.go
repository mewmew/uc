// Package semcheck implements a static semantic analysis checker for µC.
package semcheck

import (
	"github.com/mewkiz/pkg/errutil"
	"github.com/mewmew/uc/ast"
	"github.com/mewmew/uc/ast/astutil"
	"github.com/mewmew/uc/sem/errors"
)

// NoNestedFunctions disables the checking for nested functions
var NoNestedFunctions = false

// Check performs static semantic analysis on the given file.
func Check(file *ast.File) error {
	for _, decl := range file.Decls {
		switch decl := decl.(type) {
		case *ast.FuncDecl:
			// Check for nested functions.
			if NoNestedFunctions {
				if err := checkNestedFunctions(decl); err != nil {
					return errutil.Err(err)
				}
			}
		}
	}
	return nil
}

// checkNestedFunctions reports an error if the given function contains any
// nested function definitions.
func checkNestedFunctions(fn *ast.FuncDecl) error {
	if !astutil.IsDef(fn) {
		return nil
	}
	check := func(n ast.Node) error {
		if n, ok := n.(*ast.FuncDecl); ok {
			return errors.Newf(n.FuncName.Start(), "nested functions not allowed")
		}
		return nil
	}
	nop := func(ast.Node) error { return nil }
	if err := astutil.WalkBeforeAfter(fn.Body, check, nop); err != nil {
		return errutil.Err(err)
	}
	return nil
}

// TODO: Verify that all declarations occur at the beginning of the function
// body, and after the first non-declaration statement, no other declarations
// should be allowed to occur. Note, this pass should only be enabled for older
// versions of C, as newer ones allow declarations to occur throughout the
// function (albeit with other restrictions, e.g. goto may not jump over
// declarations).

// TODO: Add semantic analysis pass which verifies that declaration statements
// precedes any other statements in the body of function block.
//
//    // first specifies the first non-declaration statement within the
//    // statements of the block.
//    first := -1
//    for i, stmt := f.Body.Stmts {
//       if _, ok := stmt.(*DeclStmt); ok {
//          if first != -1 {
//             return errutil.Newf("declaration statement %v occurs after first non-declaration statement %v in function body", stmt, f.Body.Stmts[first])
//          }
//       } else if first == -1 {
//          first = i
//       }
//    }
