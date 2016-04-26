// Package sem implements a set of semantic analysis passes.
package sem

import (
	"github.com/mewkiz/pkg/errutil"
	"github.com/mewmew/uc/ast"
	"github.com/mewmew/uc/sem/semcheck"
	"github.com/mewmew/uc/sem/typecheck"
)

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

// Check performs a static semantic analysis check on the given file.
func Check(file *ast.File) error {
	if err := typecheck.Check(file); err != nil {
		return errutil.Err(err)
	}
	if err := semcheck.Check(file); err != nil {
		return errutil.Err(err)
	}
	return nil
}
