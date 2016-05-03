// Package sem implements a set of semantic analysis passes.
package sem

import (
	"github.com/mewkiz/pkg/errutil"
	"github.com/mewmew/uc/ast"
	"github.com/mewmew/uc/sem/semcheck"
	"github.com/mewmew/uc/sem/typecheck"
)

// Check performs a static semantic analysis check on the given file.
func Check(file *ast.File) error {
	// Semantic analysis is done in two passes to allow for forward references.
	// Firstly, the global declarations are added to the file-scope. Secondly,
	// the global function declaration bodies are traversed to resolve
	// identifiers and deduce the types of expressions.

	// Identifier resolution.
	if err := resolve(file); err != nil {
		return errutil.Err(err)
	}

	// Type-checking.
	if err := typecheck.Check(file); err != nil {
		return errutil.Err(err)
	}

	// Semantic analysis.
	if err := semcheck.Check(file); err != nil {
		return errutil.Err(err)
	}

	return nil
}
