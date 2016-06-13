// Package sem implements a set of semantic analysis passes.
package sem

import (
	"github.com/mewkiz/pkg/errutil"
	"github.com/mewmew/uc/ast"
	"github.com/mewmew/uc/sem/semcheck"
	"github.com/mewmew/uc/sem/typecheck"
	"github.com/mewmew/uc/types"
)

// Check performs a static semantic analysis check on the given file.
func Check(file *ast.File) (*Info, error) {
	// Semantic analysis is done in two passes to allow for forward references.
	// Firstly, the global declarations are added to the file-scope. Secondly,
	// the global function declaration bodies are traversed to resolve
	// identifiers and deduce the types of expressions.

	// Identifier resolution.
	info := &Info{
		Types:  make(map[ast.Expr]types.Type),
		Scopes: make(map[ast.Node]*Scope),
	}
	if err := resolve(file, info.Scopes); err != nil {
		return nil, errutil.Err(err)
	}

	// Type-checking.
	if err := typecheck.Check(file, info.Types); err != nil {
		return nil, errutil.Err(err)
	}

	// Semantic analysis.
	if err := semcheck.Check(file); err != nil {
		return nil, errutil.Err(err)
	}

	return info, nil
}

// TODO: Consider to move Info to uc/types.

// Info holds semantic information of a type-checked program.
type Info struct {
	// Types maps expression nodes to types.
	Types map[ast.Expr]types.Type
	// Scopes maps nodes to the scope they define.
	//
	// The following nodes define scopes.
	//
	//    *ast.File
	//    *ast.FuncDecl
	//    *ast.BlockStmt
	Scopes map[ast.Node]*Scope
}
