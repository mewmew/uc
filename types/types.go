package types

// TODO: Implement type checking which ensures correct uses of "void".
// Relevant sections of the uC BNF grammar have been included below.
//
//    TopLevelDecl
//       : VarDecl ";"
//       | TypeName ident "(" Params ")" FuncBody // TypeName : "char" | "int" | "void" ;
//    ;
//
//    ScalarDecl
//       : TypeName ident // TypeName : "char" | "int" ;
//    ;
//
//    ArrayDecl
//       : TypeName ident "[" int_lit "]" // TypeName : "char" | "int" ;
//    ;
//
//    Params
//       : TypeName   // TypeName : "void" ;
//       | ParamList
//    ;
//
//    ParamDecl
//       : ScalarDecl
//       | TypeName ident "[" "]" // TypeName : "char" | "int" ;
//    ;

// TODO: Make sure that array declarations (e.g. `int x[5]`) may only be used
// within declaration statements, and array parameter declarations (e.g. `int
// x[]`) may only be used within function signatures.

// A Type represents a type of µC, and has one of the following underlying
// types.
//
//    *Basic
//    *Array
//    *Func
type Type interface {
	// Equal reports whether t and u are of equal type.
	Equal(u Type) bool
}

type (
	// A Basic represents a basic type.
	//
	// Examples.
	//
	//    char
	//    int
	Basic struct {
		// Kind of basic type.
		Kind BasicKind
	}

	// An Array represents an array type.
	//
	// Examples.
	//
	//    int[]
	//    char[128]
	Array struct {
		// Element type.
		Elem Type
		// Array length.
		Len int
	}

	// A Func represents a function signature.
	//
	// Examples.
	//
	//    int(void)
	//    int(int a, int b)
	Func struct {
		// Return type.
		Result Type
		// Function parameter types; or nil if void parameter.
		Params []*Field
	}
)

//go:generate stringer -type BasicKind
//go:generate gorename -from basickind_string.go::i -to kind

// BasicKind describes the kind of basic type.
type BasicKind int

// Basic type.
const (
	Invalid BasicKind = iota // invalid type

	Char // "char"
	Int  // "int"
	Void // "void"
)

// A Field represents a field declaration in a struct type, or a parameter
// declaration in a function signature.
//
// Examples.
//
//    char
//    int a
type Field struct {
	// Field type.
	Type Type
	// Field name; or empty.
	Name string
}

// Equal reports whether t and u are of equal type.
func (t *Basic) Equal(u Type) bool {
	if u, ok := u.(*Basic); ok {
		return t.Kind == u.Kind
	}
	return false
}

// Equal reports whether t and u are of equal type.
func (t *Array) Equal(u Type) bool {
	if u, ok := u.(*Array); ok {
		return t.Len == u.Len && Equal(t.Elem, u.Elem)
	}
	return false
}

// Equal reports whether t and u are of equal type.
func (t *Func) Equal(u Type) bool {
	if u, ok := u.(*Func); ok {
		if !Equal(t.Result, u.Result) {
			return false
		}
		if len(u.Params) != len(t.Params) {
			return false
		}
		for i := range t.Params {
			if !Equal(t.Params[i].Type, u.Params[i].Type) {
				return false
			}
		}
		return true
	}
	return false
}

// Equal reports whether t and u are of equal type.
func Equal(t, u Type) bool {
	return t.Equal(u)
}

// Verify that the µC types implement the Type interface.
var (
	_ Type = &Basic{}
	_ Type = &Array{}
	_ Type = &Func{}
)
