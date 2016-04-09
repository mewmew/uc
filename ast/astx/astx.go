// Package astx implements utility functions for handling abstract syntax trees.
package astx

import (
	"strconv"

	"github.com/mewkiz/pkg/errutil"
	"github.com/mewmew/uc/ast"
	gocctoken "github.com/mewmew/uc/gocc/token"
)

// NewArrayDecl returns a new array declaration node, based on the following
// production rule.
//
//    TypeName ident "[" int_lit "]"
func NewArrayDecl(elem interface{}, name interface{}, length interface{}) (ast.Decl, error) {
	typ, err := NewArrayType(elem, length)
	if err != nil {
		return nil, errutil.Err(err)
	}
	return &ast.VarDecl{Name: NewIdent(name), Type: typ}, nil
}

// NewArrayType returns a new array type based on the given element type and
// length.
func NewArrayType(elem interface{}, length interface{}) (*ast.ArrayType, error) {
	n, err := strconv.Atoi(tokenString(length))
	if err != nil {
		return nil, errutil.Newf("astx.NewArrayType: invalid length; %v", err)
	}

	// Sanity checks.
	switch elem := elem.(type) {
	case *ast.BasicType:
		switch elem.Kind {
		case ast.CharType, ast.IntType:
			// Valid type.
		default:
			return nil, errutil.Newf("astx.NewArrayType: invalid kind of element basic type; %v", elem.Kind)
		}
		return &ast.ArrayType{ElemType: elem, Len: n}, nil
	default:
		return nil, errutil.Newf("astx.NewArrayType: invalid element type; %v", elem)
	}
}

// NewIdent returns a new identifier experssion node, based on the following
// production rule.
//
//    ident
func NewIdent(name interface{}) *ast.Ident {
	return &ast.Ident{Name: tokenString(name)}
}

// tokenString returns the lexeme of the given token.
func tokenString(tok interface{}) string {
	return string(tok.(*gocctoken.Token).Lit)
}
