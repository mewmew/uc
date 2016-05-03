// Package astutil implements utility functions for handling parse trees.
package astutil

import "github.com/mewmew/uc/ast"

// IsDef reports whether the given declaration is a definition.
func IsDef(decl ast.Decl) bool {
	if _, ok := decl.(*ast.VarDecl); ok {
		return true
	}
	return decl.Value() != nil
}
