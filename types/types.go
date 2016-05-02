package types

import (
	"bytes"
	"fmt"
)

// A Type represents a type of µC, and has one of the following underlying
// types.
//
//    *Basic
//    *Array
//    *Func
type Type interface {
	// Equal reports whether t and u are of equal type.
	Equal(u Type) bool
	fmt.Stringer
}

// A Numerical type is numerical if so specified by IsNumerical.
type Numerical interface {
	// IsNumerical reports whether the given type is numerical.
	IsNumerical() bool
}

// Types.
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

func (field *Field) String() string {
	if len(field.Name) > 0 {
		return fmt.Sprintf("%v %v", field.Type, field.Name)
	}
	return field.Type.String()
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

// IsVoid reports whether the given type is a void type.
func IsVoid(t Type) bool {
	if t, ok := t.(*Basic); ok {
		return t.Kind == Void
	}
	return false
}

// IsNumerical reports whether the given type is numerical.
func (t *Basic) IsNumerical() bool {
	switch t.Kind {
	case Int, Char:
		return true
	case Void:
		return false
	default:
		panic(fmt.Sprintf("types.Basic.IsNumerical: unknown basic type (%d)", int(t.Kind)))
	}
}

func (t *Basic) String() string {
	names := map[BasicKind]string{
		Char: "char",
		Int:  "int",
		Void: "void",
	}
	if s, ok := names[t.Kind]; ok {
		return s
	}
	return fmt.Sprintf("unknown basic type (%d)", int(t.Kind))
}

func (t *Array) String() string {
	if t.Len > 0 {
		return fmt.Sprintf("%v[%d]", t.Elem, t.Len)
	}
	return fmt.Sprintf("%v[]", t.Elem)
}

func (t *Func) String() string {
	buf := new(bytes.Buffer)
	fmt.Fprintf(buf, "%v(", t.Result)
	for _, param := range t.Params {
		buf.WriteString(param.String())
	}
	buf.WriteString(")")
	return buf.String()
}

// Verify that the µC types implement the Type interface.
var (
	_ Type = &Basic{}
	_ Type = &Array{}
	_ Type = &Func{}
)
