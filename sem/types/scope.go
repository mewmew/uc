package types

import (
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
	name := decl.Name().String()
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

	// First pass, add global declarations to file-scope.
	fileScope := NewScope(nil)
	for _, decl := range file.Decls {
		fileScope.Insert(decl)
	}

	// TODO: Add local declarations of functions.

	// resolve performs identifier resolution, mapping identifiers to
	// corresponding declarations of the closest lexical scope.
	resolve := func(n ast.Node) error {
		ident, ok := n.(*ast.Ident)
		if !ok {
			return nil
		}
		decl, ok := fileScope.Lookup(ident.Name)
		if !ok {
			// TODO: Remove hack and implement proper support for ident types
			switch ident.Name {
			case "char", "int", "void":
				return nil
			}

			// TODO: Report undeclared identifier.
			// TODO: Figure out how to handle basic type identifiers.
			return errutil.Newf("undeclared identifier %q", ident)
		}
		ident.Decl = decl
		return nil
	}

	if err := astutil.Walk(file, resolve); err != nil {
		return errutil.Err(err)
	}

	// 1) Identifier resolution.

	// 2) Type deduction of expressions.

	return nil
}

// isDef reports whether the given declaration is a definition.
func isDef(decl ast.Decl) bool {
	return decl.Value() != nil
}
