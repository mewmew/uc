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
		fmt.Printf("n: %#v\n", n)
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
