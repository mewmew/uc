// Package astx implements utility functions for generating abstract syntax
// trees.
package astx

import (
	"strconv"

	"github.com/mewkiz/pkg/errutil"
	"github.com/mewmew/uc/ast"
	gocctoken "github.com/mewmew/uc/gocc/token"
	"github.com/mewmew/uc/types"
)

// NewArrayDecl returns a new array declaration node, based on the following
// production rule.
//
//    TypeName ident "[" int_lit "]"
func NewArrayDecl(elem interface{}, name interface{}, length interface{}) (ast.Decl, error) {
	typ, err := NewArrayType(elem, length)
	if err != nil {
		return nil, errutil.Newf("invalid array type; %v", err)
	}
	s, err := NewIdent(name)
	if err != nil {
		return nil, errutil.Newf("invalid declaration identifier; %v", err)
	}
	return &ast.VarDecl{Name: s, Type: typ}, nil
}

// NewArrayType returns a new array type based on the given element type and
// length.
func NewArrayType(elem interface{}, length interface{}) (*types.Array, error) {
	// Parse array length.
	s, err := tokenString(length)
	if err != nil {
		return nil, errutil.Newf("invalid array length; %v", err)
	}
	n, err := strconv.Atoi(s)
	if err != nil {
		return nil, errutil.Newf("invalid array length; %v", err)
	}

	// Validate element type.
	switch elem := elem.(type) {
	case *types.Basic:
		switch elem.Kind {
		case types.Char, types.Int:
			// Valid element type.
		default:
			return nil, errutil.Newf("invalid array element type; %v", elem.Kind)
		}
		return &types.Array{Elem: elem, Len: n}, nil
	default:
		return nil, errutil.Newf("invalid array element type; %v", elem)
	}
}

// NewIdent returns a new identifier experssion node, based on the following
// production rule.
//
//    ident
func NewIdent(name interface{}) (*ast.Ident, error) {
	s, err := tokenString(name)
	if err != nil {
		return nil, errutil.Newf("invalid identifier; %v", err)
	}
	return &ast.Ident{Name: s}, nil
}

// tokenString returns the lexeme of the given token.
func tokenString(tok interface{}) (string, error) {
	switch tok := tok.(type) {
	case *gocctoken.Token:
		return string(tok.Lit), nil
	default:
		return "", errutil.Newf("invalid tok type; expected *gocctoken.Token, got %T", tok)
	}
}
