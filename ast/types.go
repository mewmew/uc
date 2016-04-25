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
		return newBasicType(n.Name)
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

// newBasicType returns a new basic type equivalent to the given identifier.
func newBasicType(name string) *types.Basic {
	var kind types.BasicKind
	switch name {
	case "char":
		kind = types.Char
	case "int":
		kind = types.Int
	case "void":
		kind = types.Void
	default:
		// TODO: Implement support for user-defined types.
		panic(fmt.Sprintf("support for basic type %q not yet implemented", name))
	}
	return &types.Basic{Kind: kind}
}
