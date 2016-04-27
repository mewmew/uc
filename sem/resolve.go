package sem

import (
	"github.com/mewkiz/pkg/errutil"
	"github.com/mewmew/uc/ast"
	"github.com/mewmew/uc/ast/astutil"
)

// resolve performs identifier resolution, mapping identifiers to corresponding
// declarations.
func resolve(file *ast.File) error {
	// TODO: Verify that type keywords cannot be redeclared.

	// Pre-pass, add keyword types and universe declarations.
	universe := NewScope(nil)
	universeDecls := []*ast.TypeDef{
		&ast.TypeDef{DeclType: &ast.Ident{Name: "char"}, TypeName: &ast.Ident{Name: "char"}},
		&ast.TypeDef{DeclType: &ast.Ident{Name: "int"}, TypeName: &ast.Ident{Name: "int"}},
		&ast.TypeDef{DeclType: &ast.Ident{Name: "void"}, TypeName: &ast.Ident{Name: "void"}},
	}
	for _, decl := range universeDecls {
		if err := universe.Insert(decl); err != nil {
			return errutil.Err(err)
		}
	}

	// First pass, add global declarations to file-scope.
	fileScope := NewScope(universe)
	for _, decl := range file.Decls {
		if err := fileScope.Insert(decl); err != nil {
			return errutil.Err(err)
		}
	}

	// TODO: Add local declarations of functions.

	scope := fileScope

	// skip specifies that the block statement body of a function declaration
	// should skip creating a nested scope, as it has already been created by its
	// function declaration, so that function parameters are placed within the
	// correct scope.
	skip := false

	// resolve performs identifier resolution, mapping identifiers to
	// corresponding declarations of the closest lexical scope.
	resolve := func(n ast.Node) error {
		switch n := n.(type) {
		case ast.Decl:
			if err := scope.Insert(n); err != nil {
				return errutil.Err(err)
			}
			if fn, ok := n.(*ast.FuncDecl); ok && astutil.IsDef(fn) {
				skip = true
				scope = NewScope(scope)
			}
		case *ast.BlockStmt:
			if !skip {
				scope = NewScope(scope)
			}
			skip = false
		case *ast.Ident:
			decl, ok := scope.Lookup(n.Name)
			if !ok {
				return errutil.Newf("%d: undeclared identifier %q", n.Start(), n)
			}
			n.Decl = decl
		}
		return nil
	}

	// after reverts to the outer scope after traversing block statements.
	after := func(n ast.Node) error {
		if _, ok := n.(*ast.BlockStmt); ok {
			scope = scope.Outer
		}
		return nil
	}

	if err := astutil.WalkBeforeAfter(file, resolve, after); err != nil {
		return errutil.Err(err)
	}

	return nil
}
