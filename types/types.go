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
	// isType ensures that only µC types can be assigned to the Type interface.
	isType()
}

type (

	// A Basic represents a basic type.
	Basic struct {
		// Kind of basic type.
		Kind BasicKind
	}

	// An Array represents an array type.
	Array struct {
		// Element type.
		Elem Type
		// Array length.
		Len int
	}

	// A Func represents a function signature.
	Func struct {
		// Function parameter types; or nil if void parameter.
		Params []*Field
		// Return type.
		Result Type
	}
)

//go:generate stringer -type BasicKind

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
type Field struct {
	// Field type.
	Type Type
	// Field name.
	Name string
}

// isType ensures that only µC types can be assigned to the Type interface.
func (n *Basic) isType() {}
func (n *Array) isType() {}
func (n *Func) isType()  {}

// Verify that the µC types implement the Type interface.
var (
	_ Type = &Basic{}
	_ Type = &Array{}
	_ Type = &Func{}
)
