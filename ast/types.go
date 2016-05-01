package ast

import (
	"fmt"

	"github.com/mewmew/uc/types"
)

// newType returns a new type equivalent to the given type node.
func newType(n Node) types.Type {
	switch n := n.(type) {
	case *ArrayType:
		return &types.Array{Elem: newType(n.Elem), Len: n.Len}
	case *FuncType:
		params := make([]*types.Field, len(n.Params))
		for i := range n.Params {
			params[i] = newField(n.Params[i])
		}
		return &types.Func{Result: newType(n.Result), Params: params}
	case *Ident:
		if n.Decl == nil {
			return newBasic(n)
		}
		return n.Decl.Type()
	default:
		panic(fmt.Sprintf("support for type %T not yet implemented", n))
	}
}

// newField returns a new field type equivalent to the given field node.
func newField(decl *VarDecl) *types.Field {
	typ := &types.Field{Type: decl.Type()}
	if decl.VarName != nil {
		typ.Name = decl.VarName.Name
	}
	return typ
}

// universePos specifies a pseudo-position used for identifiers declared in the
// universe scope.
const universePos = -1

// newBasic returns a new basic type equivalent to the given identifier.
func newBasic(ident *Ident) types.Type {
	// TODO: Check if we may come up with a cleaner solution. At least, this
	// works for now.
	switch ident.Name {
	case "char":
		charIdent := &Ident{NamePos: universePos, Name: "char"}
		charType := &types.Basic{Kind: types.Char}
		charDecl := &TypeDef{DeclType: charIdent, TypeName: charIdent, Val: charType}
		charIdent.Decl = charDecl
		ident.Decl = charDecl
		return charType
	case "int":
		intIdent := &Ident{NamePos: universePos, Name: "int"}
		intType := &types.Basic{Kind: types.Int}
		intDecl := &TypeDef{DeclType: intIdent, TypeName: intIdent, Val: intType}
		intIdent.Decl = intDecl
		ident.Decl = intDecl
		return intType
	case "void":
		voidIdent := &Ident{NamePos: universePos, Name: "void"}
		voidType := &types.Basic{Kind: types.Void}
		voidDecl := &TypeDef{DeclType: voidIdent, TypeName: voidIdent, Val: voidType}
		voidIdent.Decl = voidDecl
		ident.Decl = voidDecl
		return voidType
	default:
		panic(fmt.Sprintf("support for user-defined basic type %q not yet fully supported", ident.Name))
	}
}
