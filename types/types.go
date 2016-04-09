package types

import "fmt"

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
//       | FormalList
//    ;
//
//    FormalDecl
//       : ScalarDecl
//       | TypeName ident "[" "]" // TypeName : "char" | "int" ;
//    ;

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
		// Function argument types.
		Args []Type
		// Return type.
		Result Type
	}
)

// BasicKind describes the kind of basic type.
type BasicKind int

// Basic type.
const (
	Invalid BasicKind = iota // invalid type

	Char // "char"
	Int  // "int"
	Void // "void"
)

func (kind BasicKind) String() string {
	m := map[BasicKind]string{
		Invalid: "invalid kind of basic type",
		Char:    "char",
		Int:     "int",
		Void:    "void",
	}
	if s, ok := m[kind]; ok {
		return s
	}
	return fmt.Sprintf("unknown kind of basic type (%d)", int(kind))
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
