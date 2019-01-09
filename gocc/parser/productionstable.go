// generated by gocc; DO NOT EDIT.

package parser

import (
	"github.com/mewmew/uc/ast/astx"
	"github.com/mewmew/uc/token"
)

type (
	//TODO: change type and variable names to be consistent with other tables
	ProdTab      [numProductions]ProdTabEntry
	ProdTabEntry struct {
		String     string
		Id         string
		NTType     int
		Index      int
		NumSymbols int
		ReduceFunc func([]Attrib) (Attrib, error)
	}
	Attrib interface {
	}
)

var productionsTable = ProdTab{
	ProdTabEntry{
		String: `S' : File	<<  >>`,
		Id:         "S'",
		NTType:     0,
		Index:      0,
		NumSymbols: 1,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return X[0], nil
		},
	},
	ProdTabEntry{
		String: `File : Decls	<< astx.NewFile(X[0]) >>`,
		Id:         "File",
		NTType:     1,
		Index:      1,
		NumSymbols: 1,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return astx.NewFile(X[0])
		},
	},
	ProdTabEntry{
		String: `Decls : empty	<<  >>`,
		Id:         "Decls",
		NTType:     2,
		Index:      2,
		NumSymbols: 0,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return nil, nil
		},
	},
	ProdTabEntry{
		String: `Decls : DeclList	<<  >>`,
		Id:         "Decls",
		NTType:     2,
		Index:      3,
		NumSymbols: 1,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return X[0], nil
		},
	},
	ProdTabEntry{
		String: `DeclList : Decl	<< astx.NewDeclList(X[0]) >>`,
		Id:         "DeclList",
		NTType:     3,
		Index:      4,
		NumSymbols: 1,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return astx.NewDeclList(X[0])
		},
	},
	ProdTabEntry{
		String: `DeclList : DeclList Decl	<< astx.AppendDecl(X[0], X[1]) >>`,
		Id:         "DeclList",
		NTType:     3,
		Index:      5,
		NumSymbols: 2,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return astx.AppendDecl(X[0], X[1])
		},
	},
	ProdTabEntry{
		String: `Decl : VarDecl ";"	<< X[0], nil >>`,
		Id:         "Decl",
		NTType:     4,
		Index:      6,
		NumSymbols: 2,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return X[0], nil
		},
	},
	ProdTabEntry{
		String: `Decl : FuncDecl ";"	<< X[0], nil >>`,
		Id:         "Decl",
		NTType:     4,
		Index:      7,
		NumSymbols: 2,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return X[0], nil
		},
	},
	ProdTabEntry{
		String: `Decl : FuncDef	<<  >>`,
		Id:         "Decl",
		NTType:     4,
		Index:      8,
		NumSymbols: 1,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return X[0], nil
		},
	},
	ProdTabEntry{
		String: `Decl : TypeDef ";"	<< X[0], nil >>`,
		Id:         "Decl",
		NTType:     4,
		Index:      9,
		NumSymbols: 2,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return X[0], nil
		},
	},
	ProdTabEntry{
		String: `FuncDecl : FuncHeader	<<  >>`,
		Id:         "FuncDecl",
		NTType:     5,
		Index:      10,
		NumSymbols: 1,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return X[0], nil
		},
	},
	ProdTabEntry{
		String: `FuncHeader : BasicType ident "(" Params ")"	<< astx.NewFuncDecl(X[0], X[1], X[2], X[3], X[4]) >>`,
		Id:         "FuncHeader",
		NTType:     6,
		Index:      11,
		NumSymbols: 5,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return astx.NewFuncDecl(X[0], X[1], X[2], X[3], X[4])
		},
	},
	ProdTabEntry{
		String: `FuncDef : FuncHeader BlockStmt	<< astx.SetFuncBody(X[0], X[1]) >>`,
		Id:         "FuncDef",
		NTType:     7,
		Index:      12,
		NumSymbols: 2,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return astx.SetFuncBody(X[0], X[1])
		},
	},
	ProdTabEntry{
		String: `VarDecl : ScalarDecl	<<  >>`,
		Id:         "VarDecl",
		NTType:     8,
		Index:      13,
		NumSymbols: 1,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return X[0], nil
		},
	},
	ProdTabEntry{
		String: `VarDecl : ArrayDecl	<<  >>`,
		Id:         "VarDecl",
		NTType:     8,
		Index:      14,
		NumSymbols: 1,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return X[0], nil
		},
	},
	ProdTabEntry{
		String: `ScalarDecl : BasicType ident	<< astx.NewScalarDecl(X[0], X[1]) >>`,
		Id:         "ScalarDecl",
		NTType:     9,
		Index:      15,
		NumSymbols: 2,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return astx.NewScalarDecl(X[0], X[1])
		},
	},
	ProdTabEntry{
		String: `ArrayDecl : BasicType ident "[" IntLit "]"	<< astx.NewArrayDecl(X[0], X[1], X[2], X[3], X[4]) >>`,
		Id:         "ArrayDecl",
		NTType:     10,
		Index:      16,
		NumSymbols: 5,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return astx.NewArrayDecl(X[0], X[1], X[2], X[3], X[4])
		},
	},
	ProdTabEntry{
		String: `ArrayDecl : BasicType ident "[" "]"	<< astx.NewArrayDecl(X[0], X[1], X[2], 0, X[3]) >>`,
		Id:         "ArrayDecl",
		NTType:     10,
		Index:      17,
		NumSymbols: 4,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return astx.NewArrayDecl(X[0], X[1], X[2], 0, X[3])
		},
	},
	ProdTabEntry{
		String: `IntLit : int_lit	<< astx.NewIntLit(X[0], token.IntLit) >>`,
		Id:         "IntLit",
		NTType:     11,
		Index:      18,
		NumSymbols: 1,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return astx.NewIntLit(X[0], token.IntLit)
		},
	},
	ProdTabEntry{
		String: `IntLit : char_lit	<< astx.NewIntLit(X[0], token.CharLit) >>`,
		Id:         "IntLit",
		NTType:     11,
		Index:      19,
		NumSymbols: 1,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return astx.NewIntLit(X[0], token.CharLit)
		},
	},
	ProdTabEntry{
		String: `TypeDef : "typedef" Type ident	<< astx.NewTypeDef(X[0], X[1], X[2]) >>`,
		Id:         "TypeDef",
		NTType:     12,
		Index:      20,
		NumSymbols: 3,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return astx.NewTypeDef(X[0], X[1], X[2])
		},
	},
	ProdTabEntry{
		String: `BasicType : ident	<< astx.NewIdent(X[0]) >>`,
		Id:         "BasicType",
		NTType:     13,
		Index:      21,
		NumSymbols: 1,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return astx.NewIdent(X[0])
		},
	},
	ProdTabEntry{
		String: `Params : empty	<<  >>`,
		Id:         "Params",
		NTType:     14,
		Index:      22,
		NumSymbols: 0,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return nil, nil
		},
	},
	ProdTabEntry{
		String: `Params : ParamList	<<  >>`,
		Id:         "Params",
		NTType:     14,
		Index:      23,
		NumSymbols: 1,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return X[0], nil
		},
	},
	ProdTabEntry{
		String: `ParamList : Param	<< astx.NewParamList(X[0]) >>`,
		Id:         "ParamList",
		NTType:     15,
		Index:      24,
		NumSymbols: 1,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return astx.NewParamList(X[0])
		},
	},
	ProdTabEntry{
		String: `ParamList : ParamList "," Param	<< astx.AppendParam(X[0], X[2]) >>`,
		Id:         "ParamList",
		NTType:     15,
		Index:      25,
		NumSymbols: 3,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return astx.AppendParam(X[0], X[2])
		},
	},
	ProdTabEntry{
		String: `Param : Type	<< astx.NewAnonParam(X[0]) >>`,
		Id:         "Param",
		NTType:     16,
		Index:      26,
		NumSymbols: 1,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return astx.NewAnonParam(X[0])
		},
	},
	ProdTabEntry{
		String: `Param : VarDecl	<<  >>`,
		Id:         "Param",
		NTType:     16,
		Index:      27,
		NumSymbols: 1,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return X[0], nil
		},
	},
	ProdTabEntry{
		String: `Type : BasicType	<<  >>`,
		Id:         "Type",
		NTType:     17,
		Index:      28,
		NumSymbols: 1,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return X[0], nil
		},
	},
	ProdTabEntry{
		String: `Stmt : MatchedStmt	<<  >>`,
		Id:         "Stmt",
		NTType:     18,
		Index:      29,
		NumSymbols: 1,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return X[0], nil
		},
	},
	ProdTabEntry{
		String: `Stmt : OpenStmt	<<  >>`,
		Id:         "Stmt",
		NTType:     18,
		Index:      30,
		NumSymbols: 1,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return X[0], nil
		},
	},
	ProdTabEntry{
		String: `OtherStmt : Expr ";"	<< astx.NewExprStmt(X[0]) >>`,
		Id:         "OtherStmt",
		NTType:     19,
		Index:      31,
		NumSymbols: 2,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return astx.NewExprStmt(X[0])
		},
	},
	ProdTabEntry{
		String: `OtherStmt : "return" Expr ";"	<< astx.NewReturnStmt(X[0], X[1]) >>`,
		Id:         "OtherStmt",
		NTType:     19,
		Index:      32,
		NumSymbols: 3,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return astx.NewReturnStmt(X[0], X[1])
		},
	},
	ProdTabEntry{
		String: `OtherStmt : "return" ";"	<< astx.NewReturnStmt(X[0], nil) >>`,
		Id:         "OtherStmt",
		NTType:     19,
		Index:      33,
		NumSymbols: 2,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return astx.NewReturnStmt(X[0], nil)
		},
	},
	ProdTabEntry{
		String: `OtherStmt : BlockStmt	<<  >>`,
		Id:         "OtherStmt",
		NTType:     19,
		Index:      34,
		NumSymbols: 1,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return X[0], nil
		},
	},
	ProdTabEntry{
		String: `OtherStmt : ";"	<< astx.NewEmptyStmt(X[0]) >>`,
		Id:         "OtherStmt",
		NTType:     19,
		Index:      35,
		NumSymbols: 1,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return astx.NewEmptyStmt(X[0])
		},
	},
	ProdTabEntry{
		String: `BlockStmt : "{" BlockItems "}"	<< astx.NewBlockStmt(X[0], X[1], X[2]) >>`,
		Id:         "BlockStmt",
		NTType:     20,
		Index:      36,
		NumSymbols: 3,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return astx.NewBlockStmt(X[0], X[1], X[2])
		},
	},
	ProdTabEntry{
		String: `MatchedStmt : "if" Condition MatchedStmt "else" MatchedStmt	<< astx.NewIfStmt(X[0], X[1], X[2], X[4]) >>`,
		Id:         "MatchedStmt",
		NTType:     21,
		Index:      37,
		NumSymbols: 5,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return astx.NewIfStmt(X[0], X[1], X[2], X[4])
		},
	},
	ProdTabEntry{
		String: `MatchedStmt : "while" Condition MatchedStmt	<< astx.NewWhileStmt(X[0], X[1], X[2]) >>`,
		Id:         "MatchedStmt",
		NTType:     21,
		Index:      38,
		NumSymbols: 3,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return astx.NewWhileStmt(X[0], X[1], X[2])
		},
	},
	ProdTabEntry{
		String: `MatchedStmt : OtherStmt	<<  >>`,
		Id:         "MatchedStmt",
		NTType:     21,
		Index:      39,
		NumSymbols: 1,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return X[0], nil
		},
	},
	ProdTabEntry{
		String: `OpenStmt : "if" Condition Stmt	<< astx.NewIfStmt(X[0], X[1], X[2], nil) >>`,
		Id:         "OpenStmt",
		NTType:     22,
		Index:      40,
		NumSymbols: 3,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return astx.NewIfStmt(X[0], X[1], X[2], nil)
		},
	},
	ProdTabEntry{
		String: `OpenStmt : "if" Condition MatchedStmt "else" OpenStmt	<< astx.NewIfStmt(X[0], X[1], X[2], X[4]) >>`,
		Id:         "OpenStmt",
		NTType:     22,
		Index:      41,
		NumSymbols: 5,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return astx.NewIfStmt(X[0], X[1], X[2], X[4])
		},
	},
	ProdTabEntry{
		String: `OpenStmt : "while" Condition OpenStmt	<< astx.NewWhileStmt(X[0], X[1], X[2]) >>`,
		Id:         "OpenStmt",
		NTType:     22,
		Index:      42,
		NumSymbols: 3,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return astx.NewWhileStmt(X[0], X[1], X[2])
		},
	},
	ProdTabEntry{
		String: `Condition : "(" Expr ")"	<< X[1], nil >>`,
		Id:         "Condition",
		NTType:     23,
		Index:      43,
		NumSymbols: 3,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return X[1], nil
		},
	},
	ProdTabEntry{
		String: `BlockItems : empty	<<  >>`,
		Id:         "BlockItems",
		NTType:     24,
		Index:      44,
		NumSymbols: 0,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return nil, nil
		},
	},
	ProdTabEntry{
		String: `BlockItems : BlockItemList	<<  >>`,
		Id:         "BlockItems",
		NTType:     24,
		Index:      45,
		NumSymbols: 1,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return X[0], nil
		},
	},
	ProdTabEntry{
		String: `BlockItemList : BlockItem	<< astx.NewBlockItemList(X[0]) >>`,
		Id:         "BlockItemList",
		NTType:     25,
		Index:      46,
		NumSymbols: 1,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return astx.NewBlockItemList(X[0])
		},
	},
	ProdTabEntry{
		String: `BlockItemList : BlockItemList BlockItem	<< astx.AppendBlockItem(X[0], X[1]) >>`,
		Id:         "BlockItemList",
		NTType:     25,
		Index:      47,
		NumSymbols: 2,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return astx.AppendBlockItem(X[0], X[1])
		},
	},
	ProdTabEntry{
		String: `BlockItem : Decl	<<  >>`,
		Id:         "BlockItem",
		NTType:     26,
		Index:      48,
		NumSymbols: 1,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return X[0], nil
		},
	},
	ProdTabEntry{
		String: `BlockItem : Stmt	<<  >>`,
		Id:         "BlockItem",
		NTType:     26,
		Index:      49,
		NumSymbols: 1,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return X[0], nil
		},
	},
	ProdTabEntry{
		String: `Expr : Expr2R	<<  >>`,
		Id:         "Expr",
		NTType:     27,
		Index:      50,
		NumSymbols: 1,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return X[0], nil
		},
	},
	ProdTabEntry{
		String: `Expr2R : Expr5L	<<  >>`,
		Id:         "Expr2R",
		NTType:     28,
		Index:      51,
		NumSymbols: 1,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return X[0], nil
		},
	},
	ProdTabEntry{
		String: `Expr2R : Expr5L "=" Expr2R	<< astx.NewBinaryExpr(X[0], X[1], X[2]) >>`,
		Id:         "Expr2R",
		NTType:     28,
		Index:      52,
		NumSymbols: 3,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return astx.NewBinaryExpr(X[0], X[1], X[2])
		},
	},
	ProdTabEntry{
		String: `Expr5L : Expr9L	<<  >>`,
		Id:         "Expr5L",
		NTType:     29,
		Index:      53,
		NumSymbols: 1,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return X[0], nil
		},
	},
	ProdTabEntry{
		String: `Expr5L : Expr5L "&&" Expr9L	<< astx.NewBinaryExpr(X[0], X[1], X[2]) >>`,
		Id:         "Expr5L",
		NTType:     29,
		Index:      54,
		NumSymbols: 3,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return astx.NewBinaryExpr(X[0], X[1], X[2])
		},
	},
	ProdTabEntry{
		String: `Expr9L : Expr10L	<<  >>`,
		Id:         "Expr9L",
		NTType:     30,
		Index:      55,
		NumSymbols: 1,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return X[0], nil
		},
	},
	ProdTabEntry{
		String: `Expr9L : Expr9L "==" Expr10L	<< astx.NewBinaryExpr(X[0], X[1], X[2]) >>`,
		Id:         "Expr9L",
		NTType:     30,
		Index:      56,
		NumSymbols: 3,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return astx.NewBinaryExpr(X[0], X[1], X[2])
		},
	},
	ProdTabEntry{
		String: `Expr9L : Expr9L "!=" Expr10L	<< astx.NewBinaryExpr(X[0], X[1], X[2]) >>`,
		Id:         "Expr9L",
		NTType:     30,
		Index:      57,
		NumSymbols: 3,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return astx.NewBinaryExpr(X[0], X[1], X[2])
		},
	},
	ProdTabEntry{
		String: `Expr10L : Expr12L	<<  >>`,
		Id:         "Expr10L",
		NTType:     31,
		Index:      58,
		NumSymbols: 1,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return X[0], nil
		},
	},
	ProdTabEntry{
		String: `Expr10L : Expr10L "<" Expr12L	<< astx.NewBinaryExpr(X[0], X[1], X[2]) >>`,
		Id:         "Expr10L",
		NTType:     31,
		Index:      59,
		NumSymbols: 3,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return astx.NewBinaryExpr(X[0], X[1], X[2])
		},
	},
	ProdTabEntry{
		String: `Expr10L : Expr10L ">" Expr12L	<< astx.NewBinaryExpr(X[0], X[1], X[2]) >>`,
		Id:         "Expr10L",
		NTType:     31,
		Index:      60,
		NumSymbols: 3,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return astx.NewBinaryExpr(X[0], X[1], X[2])
		},
	},
	ProdTabEntry{
		String: `Expr10L : Expr10L "<=" Expr12L	<< astx.NewBinaryExpr(X[0], X[1], X[2]) >>`,
		Id:         "Expr10L",
		NTType:     31,
		Index:      61,
		NumSymbols: 3,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return astx.NewBinaryExpr(X[0], X[1], X[2])
		},
	},
	ProdTabEntry{
		String: `Expr10L : Expr10L ">=" Expr12L	<< astx.NewBinaryExpr(X[0], X[1], X[2]) >>`,
		Id:         "Expr10L",
		NTType:     31,
		Index:      62,
		NumSymbols: 3,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return astx.NewBinaryExpr(X[0], X[1], X[2])
		},
	},
	ProdTabEntry{
		String: `Expr12L : Expr13L	<<  >>`,
		Id:         "Expr12L",
		NTType:     32,
		Index:      63,
		NumSymbols: 1,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return X[0], nil
		},
	},
	ProdTabEntry{
		String: `Expr12L : Expr12L "+" Expr13L	<< astx.NewBinaryExpr(X[0], X[1], X[2]) >>`,
		Id:         "Expr12L",
		NTType:     32,
		Index:      64,
		NumSymbols: 3,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return astx.NewBinaryExpr(X[0], X[1], X[2])
		},
	},
	ProdTabEntry{
		String: `Expr12L : Expr12L "-" Expr13L	<< astx.NewBinaryExpr(X[0], X[1], X[2]) >>`,
		Id:         "Expr12L",
		NTType:     32,
		Index:      65,
		NumSymbols: 3,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return astx.NewBinaryExpr(X[0], X[1], X[2])
		},
	},
	ProdTabEntry{
		String: `Expr13L : Expr14	<<  >>`,
		Id:         "Expr13L",
		NTType:     33,
		Index:      66,
		NumSymbols: 1,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return X[0], nil
		},
	},
	ProdTabEntry{
		String: `Expr13L : Expr13L "*" Expr14	<< astx.NewBinaryExpr(X[0], X[1], X[2]) >>`,
		Id:         "Expr13L",
		NTType:     33,
		Index:      67,
		NumSymbols: 3,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return astx.NewBinaryExpr(X[0], X[1], X[2])
		},
	},
	ProdTabEntry{
		String: `Expr13L : Expr13L "/" Expr14	<< astx.NewBinaryExpr(X[0], X[1], X[2]) >>`,
		Id:         "Expr13L",
		NTType:     33,
		Index:      68,
		NumSymbols: 3,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return astx.NewBinaryExpr(X[0], X[1], X[2])
		},
	},
	ProdTabEntry{
		String: `Expr14 : Expr15	<<  >>`,
		Id:         "Expr14",
		NTType:     34,
		Index:      69,
		NumSymbols: 1,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return X[0], nil
		},
	},
	ProdTabEntry{
		String: `Expr14 : "-" Expr14	<< astx.NewUnaryExpr(X[0], X[1]) >>`,
		Id:         "Expr14",
		NTType:     34,
		Index:      70,
		NumSymbols: 2,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return astx.NewUnaryExpr(X[0], X[1])
		},
	},
	ProdTabEntry{
		String: `Expr14 : "!" Expr14	<< astx.NewUnaryExpr(X[0], X[1]) >>`,
		Id:         "Expr14",
		NTType:     34,
		Index:      71,
		NumSymbols: 2,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return astx.NewUnaryExpr(X[0], X[1])
		},
	},
	ProdTabEntry{
		String: `Expr15 : PrimaryExpr	<<  >>`,
		Id:         "Expr15",
		NTType:     35,
		Index:      72,
		NumSymbols: 1,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return X[0], nil
		},
	},
	ProdTabEntry{
		String: `Expr15 : ident "[" Expr "]"	<< astx.NewIndexExpr(X[0], X[1], X[2], X[3]) >>`,
		Id:         "Expr15",
		NTType:     35,
		Index:      73,
		NumSymbols: 4,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return astx.NewIndexExpr(X[0], X[1], X[2], X[3])
		},
	},
	ProdTabEntry{
		String: `Expr15 : ident "(" Args ")"	<< astx.NewCallExpr(X[0], X[1], X[2], X[3]) >>`,
		Id:         "Expr15",
		NTType:     35,
		Index:      74,
		NumSymbols: 4,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return astx.NewCallExpr(X[0], X[1], X[2], X[3])
		},
	},
	ProdTabEntry{
		String: `PrimaryExpr : int_lit	<< astx.NewBasicLit(X[0], token.IntLit) >>`,
		Id:         "PrimaryExpr",
		NTType:     36,
		Index:      75,
		NumSymbols: 1,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return astx.NewBasicLit(X[0], token.IntLit)
		},
	},
	ProdTabEntry{
		String: `PrimaryExpr : char_lit	<< astx.NewBasicLit(X[0], token.CharLit) >>`,
		Id:         "PrimaryExpr",
		NTType:     36,
		Index:      76,
		NumSymbols: 1,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return astx.NewBasicLit(X[0], token.CharLit)
		},
	},
	ProdTabEntry{
		String: `PrimaryExpr : ident	<< astx.NewIdent(X[0]) >>`,
		Id:         "PrimaryExpr",
		NTType:     36,
		Index:      77,
		NumSymbols: 1,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return astx.NewIdent(X[0])
		},
	},
	ProdTabEntry{
		String: `PrimaryExpr : ParenExpr	<<  >>`,
		Id:         "PrimaryExpr",
		NTType:     36,
		Index:      78,
		NumSymbols: 1,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return X[0], nil
		},
	},
	ProdTabEntry{
		String: `ParenExpr : "(" Expr ")"	<< astx.NewParenExpr(X[0], X[1], X[2]) >>`,
		Id:         "ParenExpr",
		NTType:     37,
		Index:      79,
		NumSymbols: 3,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return astx.NewParenExpr(X[0], X[1], X[2])
		},
	},
	ProdTabEntry{
		String: `Args : empty	<<  >>`,
		Id:         "Args",
		NTType:     38,
		Index:      80,
		NumSymbols: 0,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return nil, nil
		},
	},
	ProdTabEntry{
		String: `Args : ExprList	<<  >>`,
		Id:         "Args",
		NTType:     38,
		Index:      81,
		NumSymbols: 1,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return X[0], nil
		},
	},
	ProdTabEntry{
		String: `ExprList : Expr	<< astx.NewExprList(X[0]) >>`,
		Id:         "ExprList",
		NTType:     39,
		Index:      82,
		NumSymbols: 1,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return astx.NewExprList(X[0])
		},
	},
	ProdTabEntry{
		String: `ExprList : ExprList "," Expr	<< astx.AppendExpr(X[0], X[2]) >>`,
		Id:         "ExprList",
		NTType:     39,
		Index:      83,
		NumSymbols: 3,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return astx.AppendExpr(X[0], X[2])
		},
	},
}