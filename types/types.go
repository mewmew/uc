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
//
// All type nodes implement the ast.Node interface.
type Type interface {
	// Start returns the start position of the node within the input stream.
	Start() int
	// End returns the first character immediately after the node within the
	// input stream.
	End() int
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
type Field struct {
	// Field type.
	Type Type
	// Field name.
	Name string
}

// Start returns the start position of the node within the input stream.
func (n *Basic) Start() int { panic("ast.Basic.Start: not yet implemented") }

// Start returns the start position of the node within the input stream.
func (n *Array) Start() int { panic("ast.Array.Start: not yet implemented") }

// Start returns the start position of the node within the input stream.
func (n *Func) Start() int { panic("ast.Func.Start: not yet implemented") }

// Start returns the start position of the node within the input stream.
func (n *Field) Start() int { panic("ast.Field.Start: not yet implemented") }

// End returns the first character immediately after the node within the input
// stream.
func (n *Basic) End() int { panic("ast.Basic.End: not yet implemented") }

// End returns the first character immediately after the node within the input
// stream.
func (n *Array) End() int { panic("ast.Array.End: not yet implemented") }

// End returns the first character immediately after the node within the input
// stream.
func (n *Func) End() int { panic("ast.Func.End: not yet implemented") }

// End returns the first character immediately after the node within the input
// stream.
func (n *Field) End() int { panic("ast.Field.End: not yet implemented") }

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
