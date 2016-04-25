package types

import (
	"fmt"

	"github.com/kr/pretty"
	"github.com/mewkiz/pkg/errutil"
	"github.com/mewmew/uc/ast"
	"github.com/mewmew/uc/ast/astutil"
	"github.com/mewmew/uc/types"
)

// A Scope maintains the set of named language entities declared in the lexical
// scope and a link to the immediately surrounding outer scope.
type Scope struct {
	// Immediately surrounding outer scope; or nil if file-scope.
	Outer *Scope
	// Identifiers declared within the current scope.
	Decls map[string]ast.Decl
}

// NewScope returns a new lexical scope immediately surrouded by the given outer
// scope.
func NewScope(outer *Scope) *Scope {
	return &Scope{Outer: outer, Decls: make(map[string]ast.Decl)}
}

// Insert inserts the given declaration into the current scope.
func (s *Scope) Insert(decl ast.Decl) error {
	// Early return for first-time declarations.
	ident := decl.Name()
	if ident == nil {
		// Anonymous function parameter declaration declaration.
		return nil
	}
	name := ident.String()
	prev, ok := s.Decls[name]
	if !ok {
		s.Decls[name] = decl
		return nil
	}

	// Previously declared.
	if !types.Equal(prev.Type(), decl.Type()) {
		return errutil.Newf("%d: redefinition of %q with different type: %q vs %q", decl.Start(), name, decl.Type(), prev.Type())
	}

	// The last tentative definition becomes the definition, unless defined
	// explicitly (e.g. having an initializer or function body).
	if !isDef(prev) {
		s.Decls[name] = decl
		return nil
	}

	// Definition already present in scope.
	if isDef(decl) {
		// TODO: Split into two error messages once support for "NOTE:" error
		// messages and list of error messages have been added.
		return errutil.Newf("%d: redefinition of %q; previously declared at %d", decl.Start(), name, prev.Start())
	}

	// Declaration of previously declared identifier.
	return nil
}

// Lookup returns the declaration of name in the innermost scope of s. The
// returned boolean variable reports whether a declaration of name was located.
func (s *Scope) Lookup(name string) (ast.Decl, bool) {
	if decl, ok := s.Decls[name]; ok {
		return decl, true
	}
	if s.Outer == nil {
		return nil, false
	}
	return s.Outer.Lookup(name)
}

// Check type-checks the given file.
func Check(file *ast.File) error {
	// Type-checking is done in two passes to allow for forward references.
	// Firstly, the global declarations are added to the file-scope. Secondly,
	// the global function declaration bodies are traversed to resolve
	// identifiers and deduce the types of expressions.

	// TODO: Add keyword type declarations to universe scope. Remove special
	// cases for basic types (i.e. identifiers).

	// Pre-pass, add keyword types and universe declarations.
	//universe := NewScope(nil)
	//for _, decl := range universeDecls {
	//	universe.Insert(decl)
	//}

	// First pass, add global declarations to file-scope.
	fileScope := NewScope(nil)
	for _, decl := range file.Decls {
		fileScope.Insert(decl)
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
			scope.Insert(n)
			if fn, ok := n.(*ast.FuncDecl); ok && isDef(fn) {
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
				// TODO: Remove hack and implement proper support for ident types
				switch n.Name {
				case "char", "int", "void":
					return nil
				}

				return errutil.Newf("undeclared identifier %q", n)
			}
			n.Decl = decl
		}
		return nil
	}

	// after reverts to the outer scope after traversing block statements.
	after := func(n ast.Node) error {
		if _, ok := n.(*ast.BlockStmt); ok {
			pretty.Print(scope)
			fmt.Println("===========================================")
			scope = scope.Outer
		}
		return nil
	}

	if err := astutil.WalkBeforeAfter(file, resolve, after); err != nil {
		return errutil.Err(err)
	}

	// 1) Identifier resolution.

	// 2) Type deduction of expressions.

	// TODO: Remove debug output.

	return nil
}

// isDef reports whether the given declaration is a definition.
func isDef(decl ast.Decl) bool {
	return decl.Value() != nil
}
