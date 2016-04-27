package typecheck

import (
	"github.com/mewkiz/pkg/errutil"
	"github.com/mewmew/uc/ast"
)

// Check type-checks the given file.
func Check(file *ast.File) error {
	// Type-checking is done in two passes to allow for forward references.
	// Firstly, the global declarations are added to the file-scope. Secondly,
	// the global function declaration bodies are traversed to resolve
	// identifiers and deduce the types of expressions.

	// 2) Type deduction of expressions.
	if err := deduce(file); err != nil {
		return errutil.Err(err)
	}

	// TODO: Remove debug output.

	return nil
}
