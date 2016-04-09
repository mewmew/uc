// Package astx implements utility functions for generating abstract syntax
// trees.
package astx

import (
	"strconv"

	"github.com/mewkiz/pkg/errutil"
	"github.com/mewmew/uc/ast"
	gocctoken "github.com/mewmew/uc/gocc/token"
	"github.com/mewmew/uc/token"
	"github.com/mewmew/uc/types"
)

// NewArrayDecl returns a new array declaration node, based on the following
// production rule.
//
//    TypeName ident "[" int_lit "]"
func NewArrayDecl(elem, name, length interface{}) (ast.Decl, error) {
	typ, err := NewArrayType(elem, length)
	if err != nil {
		return nil, errutil.Newf("invalid array type; %v", err)
	}
	ident, err := NewIdent(name)
	if err != nil {
		return nil, errutil.Newf("invalid declaration identifier; %v", err)
	}
	return &ast.VarDecl{Name: ident, Type: typ}, nil
}

// NewArrayType returns a new array type based on the given element type and
// length.
func NewArrayType(elem, length interface{}) (*types.Array, error) {
	s, err := tokenString(length)
	if err != nil {
		return nil, errutil.Newf("invalid array length; %v", err)
	}
	len, err := strconv.Atoi(s)
	if err != nil {
		return nil, errutil.Newf("invalid array length; %v", err)
	}
	elemType, err := NewType(elem)
	if err != nil {
		return nil, errutil.Newf("invalid array element type; %v", err)
	}
	return &types.Array{Elem: elemType, Len: len}, nil
}

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
//    "char"
//    "int"
//    "void"
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
		return nil, errutil.Newf(`invalid basic type; expected "char", "int" or "void", got %q`, s)
	}
}

// NewExprStmt returns a new expression statement, based on the following
// production rule.
//
//    Expr ";"
func NewExprStmt(x interface{}) (*ast.ExprStmt, error) {
	if x, ok := x.(ast.Expr); ok {
		return &ast.ExprStmt{X: x}, nil
	}
	return nil, errutil.Newf("invalid expression statement expression type; expected ast.Expr, got %T", x)
}

// NewReturnStmt returns a new return statement, based on the following
// production rule.
//
//    "return" Expr ";"
func NewReturnStmt(result interface{}) (*ast.ReturnStmt, error) {
	if result, ok := result.(ast.Expr); ok {
		return &ast.ReturnStmt{Result: result}, nil
	}
	return nil, errutil.Newf("invalid return statement result type; expected ast.Expr, got %T", result)
}

// NewWhileStmt returns a new while statement, based on the following production
// rule.
//
//    "while" Condition Stmt
func NewWhileStmt(cond, body interface{}) (*ast.WhileStmt, error) {
	condExpr, ok := cond.(ast.Expr)
	if !ok {
		return nil, errutil.Newf("invalid while statement condition type; expected ast.Expr, got %T", cond)
	}
	bodyStmt, ok := body.(ast.Stmt)
	if !ok {
		return nil, errutil.Newf("invalid while statement body type; expected ast.Stmt, got %T", body)
	}
	return &ast.WhileStmt{Cond: condExpr, Body: bodyStmt}, nil
}

// NewIfStmt returns a new if statement, based on the following production
// rules.
//
//    "if" Condition Stmt ElsePart
//
//    ElsePart
//       : empty
//       | "else" Stmt
//    ;
func NewIfStmt(cond, trueBranch, falseBranch interface{}) (*ast.IfStmt, error) {
	condExpr, ok := cond.(ast.Expr)
	if !ok {
		return nil, errutil.Newf("invalid if statement condition type; expected ast.Expr, got %T", cond)
	}
	bodyStmt, ok := trueBranch.(ast.Stmt)
	if !ok {
		return nil, errutil.Newf("invalid if statement body type; expected ast.Stmt, got %T", trueBranch)
	}
	// TODO: Verify that the falseBranch != nil logic is correct for 1-way
	// conditionals.
	elseStmt, ok := cond.(ast.Stmt)
	if !ok && falseBranch != nil {
		return nil, errutil.Newf("invalid if statement else-body type; expected ast.Stmt, got %T", falseBranch)
	}
	return &ast.IfStmt{Cond: condExpr, Body: bodyStmt, Else: elseStmt}, nil
}

// NewBlockStmt returns a new block statement, based on the following production
// rule.
//
//    "{" Stmts "}"
func NewBlockStmt(stmts interface{}) (*ast.BlockStmt, error) {
	if stmts, ok := stmts.([]ast.Stmt); ok {
		return &ast.BlockStmt{Stmts: stmts}, nil
	}
	return nil, errutil.Newf("invalid block statements type; expected []ast.Stmt, got %T", stmts)
}

// NewBinaryExpr returns a new binary experssion node, based on the following
// production rules.
//
//    Expr2R "=" Expr5L
//    Expr5L "&&" Expr9L
//    Expr9L "==" Expr10L
//    Expr9L "!=" Expr10L
//    Expr10L "<" Expr12L
//    Expr10L ">" Expr12L
//    Expr10L "<=" Expr12L
//    Expr10L ">=" Expr12L
//    Expr12L "+" Expr13L
//    Expr12L "-" Expr13L
//    Expr13L "*" Expr14
//    Expr13L "/" Expr14
func NewBinaryExpr(x interface{}, op token.Kind, y interface{}) (*ast.BinaryExpr, error) {
	switch op {
	case token.Assign,
		token.Land,
		token.Eq, token.Ne,
		token.Lt, token.Gt, token.Le, token.Ge,
		token.Add, token.Sub,
		token.Mul, token.Div:
		// Valid binary operator.
	default:
		return nil, errutil.Newf("invalid binary operator; expected Assign, Land, Eq, Ne, Lt, Gt, Le, Ge, Add, Sub, Mul or Div, got %v", op)
	}
	arg0, ok := x.(ast.Expr)
	if !ok {
		return nil, errutil.Newf("invalid first binary operand type; expected ast.Expr, got %T", x)
	}
	arg1, ok := y.(ast.Expr)
	if !ok {
		return nil, errutil.Newf("invalid second binary operand type; expected ast.Expr, got %T", y)
	}
	return &ast.BinaryExpr{X: arg0, Op: op, Y: arg1}, nil
}

// NewUnaryExpr returns a new unary experssion node, based on the following
// production rules.
//
//    "-" Expr15
//    "!" Expr15
func NewUnaryExpr(op token.Kind, x interface{}) (*ast.UnaryExpr, error) {
	switch op {
	case token.Sub, token.Not:
		// Valid unary operator.
	default:
		return nil, errutil.Newf("invalid unary operator; expected Sub or Not, got %v", op)
	}
	if x, ok := x.(ast.Expr); ok {
		return &ast.UnaryExpr{Op: op, X: x}, nil
	}
	return nil, errutil.Newf("invalid unary operand type; expected ast.Expr, got %T", x)
}

// TODO: Add char_lit production rule to NewBasicLit doc comment once handled
// explicitly in uc.bnf.

// NewBasicLit returns a new basic literal experssion node of the given kind,
// based on the following production rule.
//
//    int_lit
func NewBasicLit(val interface{}, kind token.Kind) (*ast.BasicLit, error) {
	s, err := tokenString(val)
	if err != nil {
		return nil, errutil.Newf("invalid basic literal value; %v", err)
	}
	switch kind {
	case token.CharLit, token.IntLit:
		// Valid kind.
	default:
		return nil, errutil.Newf("invalid basic literal kind; expected CharLit or IntLit, got %v", kind)
	}
	return &ast.BasicLit{Kind: kind, Val: s}, nil
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

// NewIndexExpr returns a new index expression, based on the following
// production rule.
//
//    ident "[" Expr "]"
func NewIndexExpr(name, index interface{}) (*ast.IndexExpr, error) {
	ident, err := NewIdent(name)
	if err != nil {
		return nil, errutil.Newf("invalid array name; %v", err)
	}
	if index, ok := index.(ast.Expr); ok {
		return &ast.IndexExpr{Name: ident, Index: index}, nil
	}
	return nil, errutil.Newf("invalid index expression type; expected ast.Expr, got %T", index)
}

// NewCallExpr returns a new call expression, based on the following production
// rule.
//
//    ident "(" Actuals ")"
func NewCallExpr(name, args interface{}) (*ast.CallExpr, error) {
	ident, err := NewIdent(name)
	if err != nil {
		return nil, errutil.Newf("invalid function name; %v", err)
	}
	if args, ok := args.([]ast.Expr); ok {
		return &ast.CallExpr{Name: ident, Args: args}, nil
	}
	return nil, errutil.Newf("invalid function arguments type; expected []ast.Expr, got %T", args)
}

// NewParenExpr returns a new parenthesized expression, based on the following
// production rule.
//
//    "(" Expr ")"
func NewParenExpr(x interface{}) (*ast.ParenExpr, error) {
	if x, ok := x.(ast.Expr); ok {
		return &ast.ParenExpr{X: x}, nil
	}
	return nil, errutil.Newf("invalid parenthesized expression type; expected ast.Expr, got %T", x)
}

// tokenString returns the lexeme of the given token.
func tokenString(tok interface{}) (string, error) {
	if tok, ok := tok.(*gocctoken.Token); ok {
		return string(tok.Lit), nil
	}
	return "", errutil.Newf("invalid tok type; expected *gocctoken.Token, got %T", tok)
}
