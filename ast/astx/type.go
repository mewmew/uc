package astx

import (
	"github.com/mewkiz/pkg/errutil"
	"github.com/mewmew/uc/ast"
	gocctoken "github.com/mewmew/uc/gocc/token"
)

// NewType returns a new type of µC.
func NewType(typ interface{}) (ast.Type, error) {
	if typ, ok := typ.(ast.Type); ok {
		return typ, nil
	}
	return nil, errutil.Newf("invalid type; expected ast.Type, got %T", typ)
}

// NewArrayType returns a new array type based on the given element type and
// length.
func NewArrayType(elem, lbracket, length, rbracket interface{}) (*ast.ArrayType, error) {
	len, ok := length.(int)
	if !ok {
		return nil, errutil.Newf("invalid array length type; %T", length)
	}

	var lbrack, rbrack int
	switch lbracket := lbracket.(type) {
	case *gocctoken.Token:
		lbrack = lbracket.Offset
	case int:
		lbrack = lbracket
	default:
		return nil, errutil.Newf("invalid left-bracket type; expectd *gocctoken.Token or int, got %T", lbracket)
	}
	switch rbracket := rbracket.(type) {
	case *gocctoken.Token:
		rbrack = rbracket.Offset
	case int:
		rbrack = rbracket
	default:
		return nil, errutil.Newf("invalid right-bracket type; expectd *gocctoken.Token or int, got %T", rbracket)
	}

	elemType, err := NewType(elem)
	if err != nil {
		return nil, errutil.Newf("invalid array element type; %v", err)
	}
	return &ast.ArrayType{Elem: elemType, Lbracket: lbrack, Len: len, Rbracket: rbrack}, nil
}
