package astx

import (
	"strconv"

	"github.com/mewkiz/pkg/errutil"
	gocctoken "github.com/mewmew/uc/gocc/token"
	"github.com/mewmew/uc/types"
)

// NewType returns a new type of µC.
func NewType(typ interface{}) (types.Type, error) {
	if typ, ok := typ.(types.Type); ok {
		return typ, nil
	}
	return nil, errutil.Newf("invalid type; expected types.Type, got %T", typ)
}

// NewBasicType returns a new basic type of µC, based on the following
// production rules.
//
//    TypeName
//       : "char"
//       | "int"
//       | "void"
//    ;
func NewBasicType(typ interface{}) (*types.Basic, error) {
	s, err := tokenString(typ)
	if err != nil {
		return nil, errutil.Newf("invalid basic type; %v", err)
	}
	switch s {
	case "char":
		return &types.Basic{Kind: types.Char}, nil
	case "int":
		return &types.Basic{Kind: types.Int}, nil
	case "void":
		return &types.Basic{Kind: types.Void}, nil
	default:
		return nil, errutil.Newf(`invalid basic type; expected Char, Int or Void, got %q`, s)
	}
}

// NewArrayType returns a new array type based on the given element type and
// length.
func NewArrayType(elem, length interface{}) (*types.Array, error) {
	var len int
	switch length := length.(type) {
	case *gocctoken.Token:
		s, err := tokenString(length)
		if err != nil {
			return nil, errutil.Newf("invalid array length; %v", err)
		}
		len, err = strconv.Atoi(s)
		if err != nil {
			return nil, errutil.Newf("invalid array length; %v", err)
		}
	case int:
		len = length
	default:
		return nil, errutil.Newf("invalid array length type; %T", length)
	}
	elemType, err := NewType(elem)
	if err != nil {
		return nil, errutil.Newf("invalid array element type; %v", err)
	}
	return &types.Array{Elem: elemType, Len: len}, nil
}
