// Package astx implements utility functions for generating abstract syntax
// trees.
package astx

import (
	"strconv"

	"github.com/mewkiz/pkg/errutil"
	"github.com/mewmew/uc/ast"
	gocctoken "github.com/mewmew/uc/gocc/token"
	"github.com/mewmew/uc/token"
)

// NewFile returns a new µC source file, based on the following production rule.
//
//    File
//       : DeclStmts
//    ;
func NewFile(decls interface{}) (*ast.File, error) {
	if decls, ok := decls.([]ast.Decl); ok {
		return &ast.File{Decls: decls}, nil
	}
	return nil, errutil.Newf("invalid file declarations type; expected []ast.Decl, got %T", decls)
}

// NewDeclList returns a new declaration list, based on the following production
// rule.
//
//    DeclList
//       : Decl
//    ;
func NewDeclList(decl interface{}) ([]ast.Decl, error) {
	if decl, ok := decl.(ast.Decl); ok {
		return []ast.Decl{decl}, nil
	}
	return nil, errutil.Newf("invalid declaration list declaration type; expected ast.Decl, got %T", decl)
}

// AppendDecl appends decl to the declaration list, based on the following
// production rule.
//
//    DeclList
//       : DeclList Decl
//    ;
func AppendDecl(list, decl interface{}) ([]ast.Decl, error) {
	lst, ok := list.([]ast.Decl)
	if !ok {
		return nil, errutil.Newf("invalid declaration list type; expected []ast.Decl, got %T", list)
	}
	if decl, ok := decl.(ast.Decl); ok {
		return append(lst, decl), nil
	}
	return nil, errutil.Newf("invalid declaration list declaration type; expected ast.Decl, got %T", decl)
}

// NewFuncDecl returns a new function declaration node, based on the following
// production rule.
//
//    FuncDecl
//       : FuncHeader
//    ;
//
//    FuncHeader
//       : BasicType ident "(" Params ")"
//    ;
//
//    Params
//       : empty
//       | ParamList
//    ;
//
func NewFuncDecl(resultType, name, lparen, params, rparen interface{}) (*ast.FuncDecl, error) {
	resType, err := NewType(resultType)
	if err != nil {
		return nil, errutil.Newf("invalid function result type; %v", err)
	}
	ident, err := NewIdent(name)
	if err != nil {
		return nil, errutil.Newf("invalid function name identifier; %v", err)
	}
	lpar, ok := lparen.(*gocctoken.Token)
	if !ok {
		return nil, errutil.Newf("invalid left-parenthesis type; expectd *gocctoken.Token, got %T", lparen)
	}
	pars, ok := params.([]*ast.VarDecl)
	if !ok && params != nil {
		return nil, errutil.Newf("invalid function parameters type; expected []*ast.VarDecl, got %T", params)
	}
	rpar, ok := rparen.(*gocctoken.Token)
	if !ok {
		return nil, errutil.Newf("invalid right-parenthesis type; expectd *gocctoken.Token, got %T", rparen)
	}
	typ := &ast.FuncType{Result: resType, Lparen: lpar.Offset, Params: pars, Rparen: rpar.Offset}
	return &ast.FuncDecl{FuncType: typ, FuncName: ident}, nil
}

// SetFuncBody sets the function body of the given function declaration, based
// on the following production rule.
//
//    FuncDef
//       : FuncHeader BlockStmt
//    ;
func SetFuncBody(f, body interface{}) (*ast.FuncDecl, error) {
	fn, ok := f.(*ast.FuncDecl)
	if !ok {
		return nil, errutil.Newf("invalid function declaration type; expected *ast.FuncDecl, got %T", f)
	}
	if body, ok := body.(*ast.BlockStmt); ok {
		fn.Body = body
		return fn, nil
	}
	return nil, errutil.Newf("invalid function body type; expected *ast.BlockStmt, got %T", body)

}

// NewScalarDecl returns a new scalar declaration node, based on the following
// production rule.
//
//    ScalarDecl
//       : TypeName ident
//    ;
func NewScalarDecl(typ, name interface{}) (*ast.VarDecl, error) {
	scalarType, err := NewType(typ)
	if err != nil {
		return nil, errutil.Newf("invalid scalar type; %v", err)
	}
	ident, err := NewIdent(name)
	if err != nil {
		return nil, errutil.Newf("invalid scalar declaration identifier; %v", err)
	}
	return &ast.VarDecl{VarType: scalarType, VarName: ident}, nil
}

// NewArrayDecl returns a new array declaration node, based on the following
// production rule.
//
//    ArrayDecl
//       : BasicType ident "[" int_lit "]"
//    ;
func NewArrayDecl(elem, name, lbracket, length, rbracket interface{}) (*ast.VarDecl, error) {
	typ, err := NewArrayType(elem, lbracket, length, rbracket)
	if err != nil {
		return nil, errutil.Newf("invalid array type; %v", err)
	}
	ident, err := NewIdent(name)
	if err != nil {
		return nil, errutil.Newf("invalid array declaration identifier; %v", err)
	}
	return &ast.VarDecl{VarType: typ, VarName: ident}, nil
}

// NewIntLit returns a new integer, based on the following production rule.
//
//    IntLit
//       : int_lit
//       | char_lit
func NewIntLit(nToken interface{}, kind token.Kind) (int, error) {
	nTok, ok := nToken.(*gocctoken.Token)
	if !ok {
		return 0, errutil.Newf("invalid integer literal type; expectd *gocctoken.Token, got %T", nToken)
	}
	s := string(nTok.Lit)
	switch kind {
	case token.IntLit:
		n, err := strconv.Atoi(s)
		if err != nil {
			return 0, errutil.Err(err)
		}
		return n, nil
	case token.CharLit:
		s, err := strconv.Unquote(s)
		if err != nil {
			return 0, errutil.Newf("unable to unquote character literal; %v", err)
		}
		return int(s[0]), nil
	default:
		return 0, errutil.Newf(`invalid integer literal kind; expected "IntLit" or "CharLit", got %q`, kind)
	}
}

// NewTypeDef returns a new type definition node, based on the following
// production rule.
//
//    TypeDef
//       : "typedef" Type ident
//    ;
func NewTypeDef(typedefTok, typ, name interface{}) (*ast.TypeDef, error) {
	typedef, ok := typedefTok.(*gocctoken.Token)
	if !ok {
		return nil, errutil.Newf("invalid typedef keyword type; expectd *gocctoken.Token, got %T", typedefTok)
	}
	declType, err := NewType(typ)
	if err != nil {
		return nil, errutil.Newf("invalid type definition type; %v", err)
	}
	ident, err := NewIdent(name)
	if err != nil {
		return nil, errutil.Newf("invalid type definition identifier; %v", err)
	}
	return &ast.TypeDef{Typedef: typedef.Offset, DeclType: declType, TypeName: ident}, nil
}

// NewParamList returns a new parameter list, based on the following production
// rule.
//
//    ParamList
//       : Param
//    ;
func NewParamList(param interface{}) ([]*ast.VarDecl, error) {
	if param, ok := param.(*ast.VarDecl); ok {
		return []*ast.VarDecl{param}, nil
	}
	return nil, errutil.Newf("invalid parameter list parameter type; expected *ast.VarDecl, got %T", param)
}

// AppendParam appends parameter to the parameter list, based on the following
// production rule.
//
//    ParamList
//       : ParamList "," Param
//    ;
func AppendParam(list, param interface{}) ([]*ast.VarDecl, error) {
	lst, ok := list.([]*ast.VarDecl)
	if !ok {
		return nil, errutil.Newf("invalid parameter list type; expected []*ast.VarDecl, got %T", list)
	}
	if param, ok := param.(*ast.VarDecl); ok {
		return append(lst, param), nil
	}
	return nil, errutil.Newf("invalid parameter list parameter type; expected *ast.VarDecl, got %T", param)
}

// NewAnonParam returns a new anonymous parameter, based on the following
// production rules.
//
//    Param
//       // BasicType : "void" ;
//       : Type
//    ;
func NewAnonParam(typ interface{}) (*ast.VarDecl, error) {
	if typ, ok := typ.(ast.Type); ok {
		return &ast.VarDecl{VarType: typ}, nil
	}
	return nil, errutil.Newf("invalid anonymous parameter type; expected ast.Type, got %T", typ)
}

// NewExprStmt returns a new expression statement, based on the following
// production rule.
//
//    Stmt
//       : Expr ";"
//    ;
func NewExprStmt(x interface{}) (*ast.ExprStmt, error) {
	if x, ok := x.(ast.Expr); ok {
		return &ast.ExprStmt{X: x}, nil
	}
	return nil, errutil.Newf("invalid expression statement expression type; expected ast.Expr, got %T", x)
}

// NewReturnStmt returns a new return statement, based on the following
// production rule.
//
//    Stmt
//       : "return" Expr ";"
//    ;
func NewReturnStmt(returnToken, result interface{}) (*ast.ReturnStmt, error) {
	retTok, ok := returnToken.(*gocctoken.Token)
	if !ok {
		return nil, errutil.Newf("invalid return keyword type; expected *gocctoken.Token, got %T", returnToken)
	}
	if result == nil {
		return &ast.ReturnStmt{Return: retTok.Offset}, nil
	}
	if result, ok := result.(ast.Expr); ok {
		return &ast.ReturnStmt{Return: retTok.Offset, Result: result}, nil
	}
	return nil, errutil.Newf("invalid return statement result type; expected ast.Expr, got %T", result)
}

// NewWhileStmt returns a new while statement, based on the following production
// rule.
//
//    Stmt
//       : "while" Condition Stmt
//    ;
func NewWhileStmt(whileToken, cond, body interface{}) (*ast.WhileStmt, error) {
	whileTok, ok := whileToken.(*gocctoken.Token)
	if !ok {
		return nil, errutil.Newf("invalid while keyword type; expected *gocctoken.Token, got %T", whileToken)
	}
	condExpr, ok := cond.(ast.Expr)
	if !ok {
		return nil, errutil.Newf("invalid while statement condition type; expected ast.Expr, got %T", cond)
	}
	bodyStmt, ok := body.(ast.Stmt)
	if !ok {
		return nil, errutil.Newf("invalid while statement body type; expected ast.Stmt, got %T", body)
	}
	return &ast.WhileStmt{While: whileTok.Offset, Cond: condExpr, Body: bodyStmt}, nil
}

// NewIfStmt returns a new if statement, based on the following production
// rules.
//
//    Stmt
//       : "if" Condition Stmt ElsePart
//    ;
//
//    ElsePart
//       : empty
//       | "else" Stmt
//    ;
func NewIfStmt(ifToken, cond, trueBranch, falseBranch interface{}) (*ast.IfStmt, error) {
	ifTok, ok := ifToken.(*gocctoken.Token)
	if !ok {
		return nil, errutil.Newf("invalid if keyword type; expected *gocctoken.Token, got %T", ifToken)
	}
	condExpr, ok := cond.(ast.Expr)
	if !ok {
		return nil, errutil.Newf("invalid if statement condition type; expected ast.Expr, got %T", cond)
	}
	bodyStmt, ok := trueBranch.(ast.Stmt)
	if !ok {
		return nil, errutil.Newf("invalid if statement body type; expected ast.Stmt, got %T", trueBranch)
	}
	if falseBranch == nil {
		return &ast.IfStmt{If: ifTok.Offset, Cond: condExpr, Body: bodyStmt}, nil
	}
	if elseStmt, ok := falseBranch.(ast.Stmt); ok {
		return &ast.IfStmt{If: ifTok.Offset, Cond: condExpr, Body: bodyStmt, Else: elseStmt}, nil
	}
	return nil, errutil.Newf("invalid if statement else-body type; expected ast.Stmt, got %T", falseBranch)
}

// NewBlockStmt returns a new block statement, based on the following production
// rule.
//
//    BlockStmt
//       : "{" BlockItems "}"
//    ;
func NewBlockStmt(lbrace, items, rbrace interface{}) (*ast.BlockStmt, error) {
	lbra, ok := lbrace.(*gocctoken.Token)
	if !ok {
		return nil, errutil.Newf("invalid left-brace type; expectd *gocctoken.Token, got %T", lbrace)
	}
	rbra, ok := rbrace.(*gocctoken.Token)
	if !ok {
		return nil, errutil.Newf("invalid right-brace type; expectd *gocctoken.Token, got %T", rbrace)
	}
	if items == nil {
		return &ast.BlockStmt{Lbrace: lbra.Offset, Rbrace: rbra.Offset}, nil
	}
	if items, ok := items.([]ast.BlockItem); ok {
		return &ast.BlockStmt{Lbrace: lbra.Offset, Items: items, Rbrace: rbra.Offset}, nil
	}
	return nil, errutil.Newf("invalid block statements type; expected []ast.BlockItem, got %T", items)
}

// NewBlockItemList returns a new block item list, based on the following
// production rule.
//
//    BlockItemList
//       : BlockItem
//    ;
func NewBlockItemList(item interface{}) ([]ast.BlockItem, error) {
	if item, ok := item.(ast.BlockItem); ok {
		return []ast.BlockItem{item}, nil
	}
	return nil, errutil.Newf("invalid block item list block item type; expected ast.BlockItem, got %T", item)
}

// AppendBlockItem appends item to the block item list, based on the following
// production rule.
//
//    BlockItemList
//       : BlockItemList BlockItem
//    ;
func AppendBlockItem(list, item interface{}) ([]ast.BlockItem, error) {
	lst, ok := list.([]ast.BlockItem)
	if !ok {
		return nil, errutil.Newf("invalid block item list type; expected []ast.BlockItem, got %T", list)
	}
	if item, ok := item.(ast.BlockItem); ok {
		return append(lst, item), nil
	}
	return nil, errutil.Newf("invalid block item list block item type; expected ast.BlockItem, got %T", item)
}

// NewEmptyStmt returns a new empty statement, based on the following production
// rules.
//
//    Stmt
//       : ";"
//    ;
func NewEmptyStmt(semicolonToken interface{}) (*ast.EmptyStmt, error) {
	semiTok, ok := semicolonToken.(*gocctoken.Token)
	if !ok {
		return nil, errutil.Newf("invalid semicolon type; expected *gocctoken.Token, got %T", semicolonToken)
	}
	return &ast.EmptyStmt{Semicolon: semiTok.Offset}, nil
}

// NewBinaryExpr returns a new binary experssion node, based on the following
// production rules.
//
//    Expr2R
//       : Expr2R "=" Expr5L
//    ;
//
//    Expr5L
//       : Expr5L "&&" Expr9L
//    ;
//
//    Expr9L
//       : Expr9L "==" Expr10L
//       | Expr9L "!=" Expr10L
//    ;
//
//    Expr10L
//       : Expr10L "<" Expr12L
//       | Expr10L ">" Expr12L
//       | Expr10L "<=" Expr12L
//       | Expr10L ">=" Expr12L
//    ;
//
//    Expr12L
//       : Expr12L "+" Expr13L
//       | Expr12L "-" Expr13L
//    ;
//
//    Expr13L
//       : Expr13L "*" Expr14
//       | Expr13L "/" Expr14
//    ;
func NewBinaryExpr(x, opToken, y interface{}) (*ast.BinaryExpr, error) {
	opTok, ok := opToken.(*gocctoken.Token)
	if !ok {
		return nil, errutil.Newf("invalid binary operator type; expectd *gocctoken.Token, got %T", opToken)
	}
	var op token.Kind
	switch lit := string(opTok.Lit); lit {
	case "=":
		op = token.Assign
	case "&&":
		op = token.Land
	case "==":
		op = token.Eq
	case "!=":
		op = token.Ne
	case "<":
		op = token.Lt
	case ">":
		op = token.Gt
	case "<=":
		op = token.Le
	case ">=":
		op = token.Ge
	case "+":
		op = token.Add
	case "-":
		op = token.Sub
	case "*":
		op = token.Mul
	case "/":
		op = token.Div
	default:
		return nil, errutil.Newf(`invalid binary operator; expected "=", "&&", "==", "!=", "<", ">", "<=", ">=", "+", "-", "*" or "/", got %q`, lit)
	}

	arg0, ok := x.(ast.Expr)
	if !ok {
		return nil, errutil.Newf("invalid first binary operand type; expected ast.Expr, got %T", x)
	}
	arg1, ok := y.(ast.Expr)
	if !ok {
		return nil, errutil.Newf("invalid second binary operand type; expected ast.Expr, got %T", y)
	}
	return &ast.BinaryExpr{X: arg0, OpPos: opTok.Offset, Op: op, Y: arg1}, nil
}

// NewUnaryExpr returns a new unary experssion node, based on the following
// production rules.
//
//    Expr14
//       : "-" Expr15
//       | "!" Expr15
//    ;
func NewUnaryExpr(opToken, x interface{}) (*ast.UnaryExpr, error) {
	opTok, ok := opToken.(*gocctoken.Token)
	if !ok {
		return nil, errutil.Newf("invalid unary operator type; expectd *gocctoken.Token, got %T", opToken)
	}
	var op token.Kind
	switch lit := string(opTok.Lit); lit {
	case "-":
		op = token.Sub
	case "!":
		op = token.Not
	default:
		return nil, errutil.Newf(`invalid unary operator; expected "-" or "!", got %q`, lit)
	}
	if x, ok := x.(ast.Expr); ok {
		return &ast.UnaryExpr{OpPos: opTok.Offset, Op: op, X: x}, nil
	}
	return nil, errutil.Newf("invalid unary operand type; expected ast.Expr, got %T", x)
}

// NewBasicLit returns a new basic literal experssion node of the given kind,
// based on the following production rule.
//
//    Expr15
//       : int_lit
//       | char_lit
//    ;
func NewBasicLit(valToken interface{}, kind token.Kind) (*ast.BasicLit, error) {
	valTok, ok := valToken.(*gocctoken.Token)
	if !ok {
		return nil, errutil.Newf("invalid basic literal type; expected *gocctoken.Token, got %T", valToken)
	}
	switch kind {
	case token.CharLit, token.IntLit:
		// Valid kind.
	default:
		return nil, errutil.Newf("invalid basic literal kind; expected CharLit or IntLit, got %v", kind)
	}
	return &ast.BasicLit{ValPos: valTok.Offset, Kind: kind, Val: string(valTok.Lit)}, nil
}

// NewIdent returns a new identifier experssion node, based on the following
// production rule.
//
//    Expr15
//       : ident
//    ;
func NewIdent(nameToken interface{}) (*ast.Ident, error) {
	nameTok, ok := nameToken.(*gocctoken.Token)
	if !ok {
		return nil, errutil.Newf("invalid identifier type; expected *gocctoken.Token, got %T", nameToken)
	}
	return &ast.Ident{NamePos: nameTok.Offset, Name: string(nameTok.Lit)}, nil
}

// NewIndexExpr returns a new index expression, based on the following
// production rule.
//
//    Expr15
//       : ident "[" Expr "]"
//    ;
func NewIndexExpr(name, lbracket, index, rbracket interface{}) (*ast.IndexExpr, error) {
	ident, err := NewIdent(name)
	if err != nil {
		return nil, errutil.Newf("invalid array name; %v", err)
	}
	lbrack, ok := lbracket.(*gocctoken.Token)
	if !ok {
		return nil, errutil.Newf("invalid left-bracket type; expectd *gocctoken.Token, got %T", lbracket)
	}
	rbrack, ok := rbracket.(*gocctoken.Token)
	if !ok {
		return nil, errutil.Newf("invalid right-bracket type; expectd *gocctoken.Token, got %T", rbracket)
	}
	if index, ok := index.(ast.Expr); ok {
		return &ast.IndexExpr{Name: ident, Lbracket: lbrack.Offset, Index: index, Rbracket: rbrack.Offset}, nil
	}
	return nil, errutil.Newf("invalid index expression type; expected ast.Expr, got %T", index)
}

// NewCallExpr returns a new call expression, based on the following production
// rule.
//
//    Expr15
//       : ident "(" Args ")"
//    ;
func NewCallExpr(name, lparen, args, rparen interface{}) (*ast.CallExpr, error) {
	ident, err := NewIdent(name)
	if err != nil {
		return nil, errutil.Newf("invalid function name; %v", err)
	}
	lpar, ok := lparen.(*gocctoken.Token)
	if !ok {
		return nil, errutil.Newf("invalid left-parenthesis type; expectd *gocctoken.Token, got %T", lparen)
	}
	rpar, ok := rparen.(*gocctoken.Token)
	if !ok {
		return nil, errutil.Newf("invalid right-parenthesis type; expectd *gocctoken.Token, got %T", rparen)
	}
	if args == nil {
		return &ast.CallExpr{Name: ident, Lparen: lpar.Offset, Rparen: rpar.Offset}, nil
	}
	if args, ok := args.([]ast.Expr); ok {
		return &ast.CallExpr{Name: ident, Lparen: lpar.Offset, Args: args, Rparen: rpar.Offset}, nil
	}
	return nil, errutil.Newf("invalid function arguments type; expected []ast.Expr, got %T", args)
}

// NewParenExpr returns a new parenthesized expression, based on the following
// production rule.
//
//    ParenExpr
//       : "(" Expr ")"
//    ;
func NewParenExpr(lparen, x, rparen interface{}) (*ast.ParenExpr, error) {
	lpar, ok := lparen.(*gocctoken.Token)
	if !ok {
		return nil, errutil.Newf("invalid left-parenthesis type; expectd *gocctoken.Token, got %T", lparen)
	}
	rpar, ok := rparen.(*gocctoken.Token)
	if !ok {
		return nil, errutil.Newf("invalid right-parenthesis type; expectd *gocctoken.Token, got %T", rparen)
	}
	if x, ok := x.(ast.Expr); ok {
		return &ast.ParenExpr{Lparen: lpar.Offset, X: x, Rparen: rpar.Offset}, nil
	}
	return nil, errutil.Newf("invalid parenthesized expression type; expected ast.Expr, got %T", x)
}

// NewExprList returns a new expression list, based on the following production
// rule.
//
//    ExprList
//       : Expr
//    ;
func NewExprList(x interface{}) ([]ast.Expr, error) {
	if x, ok := x.(ast.Expr); ok {
		return []ast.Expr{x}, nil
	}
	return nil, errutil.Newf("invalid expression list expression type; expected ast.Expr, got %T", x)
}

// AppendExpr appends x to the expression list, based on the following
// production rule.
//
//    ExprList
//       : ExprList "," Expr
//    ;
func AppendExpr(list, x interface{}) ([]ast.Expr, error) {
	lst, ok := list.([]ast.Expr)
	if !ok {
		return nil, errutil.Newf("invalid expression list type; expected []ast.Expr, got %T", list)
	}
	if x, ok := x.(ast.Expr); ok {
		return append(lst, x), nil
	}
	return nil, errutil.Newf("invalid expression list expression type; expected ast.Expr, got %T", x)
}
