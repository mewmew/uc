package parser_test

import (
	"fmt"
	"log"
	"reflect"
	"testing"

	"github.com/kr/pretty"
	"github.com/mewmew/uc/ast"
	"github.com/mewmew/uc/gocc/errors"
	"github.com/mewmew/uc/gocc/parser"
	"github.com/mewmew/uc/gocc/scanner"
	"github.com/mewmew/uc/token"
)

func TestParser(t *testing.T) {
	var golden = []struct {
		path string
		want *ast.File
	}{
		{
			path: "../../testdata/quiet/lexer/l05.c",
			want: &ast.File{
				Decls: []ast.Decl{
					&ast.FuncDecl{
						FuncType: &ast.FuncType{
							Result: &ast.Ident{
								NamePos: 0,
								Name:    "int",
							},
							Lparen: 8,
							Params: []*ast.Field{
								{
									Type: &ast.Ident{
										NamePos: 9,
										Name:    "void",
									},
								},
							},
							Rparen: 13,
						},
						FuncName: &ast.Ident{
							NamePos: 4,
							Name:    "main",
						},
						Body: &ast.BlockStmt{
							Lbrace: 15,
							Items: []ast.BlockItem{
								&ast.VarDecl{
									VarType: &ast.Ident{
										NamePos: 19,
										Name:    "int",
									},
									VarName: &ast.Ident{
										NamePos: 23,
										Name:    "i",
									},
								},
								&ast.ExprStmt{
									X: &ast.BinaryExpr{
										X: &ast.BasicLit{
											ValPos: 28,
											Kind:   token.IntLit,
											Val:    "1",
										},
										OpPos: 29,
										Op:    token.Ne,
										Y: &ast.UnaryExpr{
											OpPos: 31,
											Op:    token.Not,
											X: &ast.BasicLit{
												ValPos: 32,
												Kind:   token.IntLit,
												Val:    "3",
											},
										},
									},
								},
								&ast.ExprStmt{
									X: &ast.BinaryExpr{
										X: &ast.BasicLit{
											ValPos: 37,
											Kind:   token.IntLit,
											Val:    "4",
										},
										OpPos: 38,
										Op:    token.Land,
										Y: &ast.ParenExpr{
											Lparen: 40,
											X: &ast.BasicLit{
												ValPos: 41,
												Kind:   token.IntLit,
												Val:    "6",
											},
											Rparen: 42,
										},
									},
								},
								&ast.ExprStmt{
									X: &ast.BinaryExpr{
										X: &ast.BinaryExpr{
											X: &ast.BasicLit{
												ValPos: 47,
												Kind:   token.IntLit,
												Val:    "7",
											},
											OpPos: 48,
											Op:    token.Mul,
											Y: &ast.BasicLit{
												ValPos: 50,
												Kind:   token.IntLit,
												Val:    "8",
											},
										},
										OpPos: 51,
										Op:    token.Add,
										Y: &ast.BasicLit{
											ValPos: 52,
											Kind:   token.IntLit,
											Val:    "10",
										},
									},
								},
								&ast.ExprStmt{
									X: &ast.BinaryExpr{
										X: &ast.ParenExpr{
											Lparen: 58,
											X: &ast.BinaryExpr{
												X: &ast.BasicLit{
													ValPos: 59,
													Kind:   token.IntLit,
													Val:    "11",
												},
												OpPos: 61,
												Op:    token.Sub,
												Y: &ast.BasicLit{
													ValPos: 62,
													Kind:   token.IntLit,
													Val:    "12",
												},
											},
											Rparen: 64,
										},
										OpPos: 65,
										Op:    token.Add,
										Y: &ast.ParenExpr{
											Lparen: 66,
											X: &ast.BinaryExpr{
												X: &ast.BasicLit{
													ValPos: 67,
													Kind:   token.IntLit,
													Val:    "12",
												},
												OpPos: 69,
												Op:    token.Div,
												Y: &ast.BasicLit{
													ValPos: 70,
													Kind:   token.IntLit,
													Val:    "16",
												},
											},
											Rparen: 72,
										},
									},
								},
								&ast.ExprStmt{
									X: &ast.BinaryExpr{
										X: &ast.BinaryExpr{
											X: &ast.BasicLit{
												ValPos: 77,
												Kind:   token.IntLit,
												Val:    "17",
											},
											OpPos: 79,
											Op:    token.Le,
											Y: &ast.BasicLit{
												ValPos: 81,
												Kind:   token.IntLit,
												Val:    "18",
											},
										},
										OpPos: 84,
										Op:    token.Lt,
										Y: &ast.UnaryExpr{
											OpPos: 85,
											Op:    token.Sub,
											X: &ast.BasicLit{
												ValPos: 86,
												Kind:   token.IntLit,
												Val:    "20",
											},
										},
									},
								},
								&ast.ExprStmt{
									X: &ast.BinaryExpr{
										X: &ast.Ident{
											NamePos: 92,
											Name:    "i",
										},
										OpPos: 93,
										Op:    token.Assign,
										Y: &ast.BinaryExpr{
											X: &ast.BasicLit{
												ValPos: 94,
												Kind:   token.IntLit,
												Val:    "21",
											},
											OpPos: 96,
											Op:    token.Eq,
											Y: &ast.BasicLit{
												ValPos: 98,
												Kind:   token.IntLit,
												Val:    "22",
											},
										},
									},
								},
								&ast.ExprStmt{
									X: &ast.BinaryExpr{
										X: &ast.BinaryExpr{
											X: &ast.BasicLit{
												ValPos: 104,
												Kind:   token.IntLit,
												Val:    "25",
											},
											OpPos: 107,
											Op:    token.Ge,
											Y: &ast.BasicLit{
												ValPos: 109,
												Kind:   token.IntLit,
												Val:    "27",
											},
										},
										OpPos: 111,
										Op:    token.Gt,
										Y: &ast.BasicLit{
											ValPos: 112,
											Kind:   token.IntLit,
											Val:    "28",
										},
									},
								},
							},
							Rbrace: 116,
						},
					},
				},
			},
		},

		{
			path: "../../testdata/quiet/parser/p01.c",
			want: &ast.File{
				Decls: []ast.Decl{
					&ast.FuncDecl{
						FuncType: &ast.FuncType{
							Result: &ast.Ident{
								NamePos: 0,
								Name:    "int",
							},
							Lparen: 8,
							Params: []*ast.Field{
								{
									Type: &ast.Ident{
										NamePos: 9,
										Name:    "void",
									},
								},
							},
							Rparen: 13,
						},
						FuncName: &ast.Ident{
							NamePos: 4,
							Name:    "main",
						},
						Body: &ast.BlockStmt{
							Lbrace: 15,
							Items: []ast.BlockItem{
								&ast.VarDecl{
									VarType: &ast.Ident{
										NamePos: 19,
										Name:    "int",
									},
									VarName: &ast.Ident{
										NamePos: 23,
										Name:    "x",
									},
								},
								&ast.VarDecl{
									VarType: &ast.Ident{
										NamePos: 28,
										Name:    "int",
									},
									VarName: &ast.Ident{
										NamePos: 32,
										Name:    "y",
									},
								},
								&ast.ExprStmt{
									X: &ast.BinaryExpr{
										X: &ast.Ident{
											NamePos: 37,
											Name:    "x",
										},
										OpPos: 39,
										Op:    token.Assign,
										Y: &ast.BasicLit{
											ValPos: 41,
											Kind:   token.IntLit,
											Val:    "42",
										},
									},
								},
								&ast.ExprStmt{
									X: &ast.BinaryExpr{
										X: &ast.Ident{
											NamePos: 47,
											Name:    "x",
										},
										OpPos: 48,
										Op:    token.Assign,
										Y: &ast.BinaryExpr{
											X: &ast.Ident{
												NamePos: 49,
												Name:    "y",
											},
											OpPos: 50,
											Op:    token.Assign,
											Y: &ast.BasicLit{
												ValPos: 51,
												Kind:   token.IntLit,
												Val:    "4711",
											},
										},
									},
								},
							},
							Rbrace: 57,
						},
					},
				},
			},
		},

		{
			path: "../../testdata/quiet/parser/p02.c",
			want: &ast.File{
				Decls: []ast.Decl{
					&ast.FuncDecl{
						FuncType: &ast.FuncType{
							Result: &ast.Ident{
								NamePos: 1,
								Name:    "int",
							},
							Lparen: 9,
							Params: []*ast.Field{
								{
									Type: &ast.Ident{
										NamePos: 10,
										Name:    "void",
									},
								},
							},
							Rparen: 14,
						},
						FuncName: &ast.Ident{
							NamePos: 5,
							Name:    "main",
						},
						Body: &ast.BlockStmt{
							Lbrace: 16,
							Items: []ast.BlockItem{
								&ast.VarDecl{
									VarType: &ast.Ident{
										NamePos: 20,
										Name:    "int",
									},
									VarName: &ast.Ident{
										NamePos: 24,
										Name:    "x",
									},
								},
								&ast.EmptyStmt{
									Semicolon: 29,
								},
								&ast.WhileStmt{
									While: 33,
									Cond: &ast.BinaryExpr{
										X: &ast.Ident{
											NamePos: 40,
											Name:    "x",
										},
										OpPos: 41,
										Op:    token.Lt,
										Y: &ast.BasicLit{
											ValPos: 42,
											Kind:   token.IntLit,
											Val:    "10",
										},
									},
									Body: &ast.ExprStmt{
										X: &ast.BinaryExpr{
											X: &ast.Ident{
												NamePos: 46,
												Name:    "x",
											},
											OpPos: 48,
											Op:    token.Assign,
											Y: &ast.BinaryExpr{
												X: &ast.Ident{
													NamePos: 50,
													Name:    "x",
												},
												OpPos: 52,
												Op:    token.Add,
												Y: &ast.BasicLit{
													ValPos: 54,
													Kind:   token.IntLit,
													Val:    "3",
												},
											},
										},
									},
								},
								&ast.IfStmt{
									If: 60,
									Cond: &ast.BasicLit{
										ValPos: 64,
										Kind:   token.IntLit,
										Val:    "1",
									},
									Body: &ast.ExprStmt{
										X: &ast.BinaryExpr{
											X: &ast.Ident{
												NamePos: 67,
												Name:    "x",
											},
											OpPos: 69,
											Op:    token.Assign,
											Y: &ast.BinaryExpr{
												X: &ast.Ident{
													NamePos: 71,
													Name:    "x",
												},
												OpPos: 73,
												Op:    token.Add,
												Y: &ast.BasicLit{
													ValPos: 75,
													Kind:   token.IntLit,
													Val:    "3",
												},
											},
										},
									},
								},
							},
							Rbrace: 78,
						},
					},
				},
			},
		},

		{
			path: "../../testdata/quiet/parser/p03.c",
			want: &ast.File{
				Decls: []ast.Decl{
					&ast.FuncDecl{
						FuncType: &ast.FuncType{
							Result: &ast.Ident{
								NamePos: 1,
								Name:    "int",
							},
							Lparen: 9,
							Params: []*ast.Field{
								{
									Type: &ast.Ident{
										NamePos: 10,
										Name:    "void",
									},
								},
							},
							Rparen: 14,
						},
						FuncName: &ast.Ident{
							NamePos: 5,
							Name:    "main",
						},
						Body: &ast.BlockStmt{
							Lbrace: 16,
							Items: []ast.BlockItem{
								&ast.VarDecl{
									VarType: &ast.Ident{
										NamePos: 20,
										Name:    "int",
									},
									VarName: &ast.Ident{
										NamePos: 24,
										Name:    "x",
									},
								},
								&ast.IfStmt{
									If: 29,
									Cond: &ast.BinaryExpr{
										X: &ast.BasicLit{
											ValPos: 33,
											Kind:   token.IntLit,
											Val:    "1",
										},
										OpPos: 34,
										Op:    token.Lt,
										Y: &ast.BasicLit{
											ValPos: 35,
											Kind:   token.IntLit,
											Val:    "2",
										},
									},
									Body: &ast.ExprStmt{
										X: &ast.BinaryExpr{
											X: &ast.Ident{
												NamePos: 38,
												Name:    "x",
											},
											OpPos: 40,
											Op:    token.Assign,
											Y: &ast.BasicLit{
												ValPos: 42,
												Kind:   token.IntLit,
												Val:    "1",
											},
										},
									},
									Else: &ast.ExprStmt{
										X: &ast.BinaryExpr{
											X: &ast.Ident{
												NamePos: 53,
												Name:    "x",
											},
											OpPos: 55,
											Op:    token.Assign,
											Y: &ast.BasicLit{
												ValPos: 57,
												Kind:   token.IntLit,
												Val:    "2",
											},
										},
									},
								},
							},
							Rbrace: 60,
						},
					},
				},
			},
		},

		{
			path: "../../testdata/quiet/parser/p04.c",
			want: &ast.File{
				Decls: []ast.Decl{
					&ast.FuncDecl{
						FuncType: &ast.FuncType{
							Result: &ast.Ident{
								NamePos: 0,
								Name:    "int",
							},
							Lparen: 8,
							Params: []*ast.Field{
								{
									Type: &ast.Ident{
										NamePos: 9,
										Name:    "void",
									},
								},
							},
							Rparen: 13,
						},
						FuncName: &ast.Ident{
							NamePos: 4,
							Name:    "main",
						},
						Body: &ast.BlockStmt{
							Lbrace: 15,
							Items: []ast.BlockItem{
								&ast.VarDecl{
									VarType: &ast.Ident{
										NamePos: 19,
										Name:    "int",
									},
									VarName: &ast.Ident{
										NamePos: 23,
										Name:    "x",
									},
								},
								&ast.VarDecl{
									VarType: &ast.Ident{
										NamePos: 28,
										Name:    "int",
									},
									VarName: &ast.Ident{
										NamePos: 32,
										Name:    "y",
									},
								},
								&ast.VarDecl{
									VarType: &ast.Ident{
										NamePos: 37,
										Name:    "int",
									},
									VarName: &ast.Ident{
										NamePos: 41,
										Name:    "z",
									},
								},
								&ast.ExprStmt{
									X: &ast.BinaryExpr{
										X: &ast.BinaryExpr{
											X: &ast.BinaryExpr{
												X: &ast.Ident{
													NamePos: 49,
													Name:    "x",
												},
												OpPos: 50,
												Op:    token.Sub,
												Y: &ast.Ident{
													NamePos: 51,
													Name:    "y",
												},
											},
											OpPos: 52,
											Op:    token.Sub,
											Y: &ast.Ident{
												NamePos: 53,
												Name:    "z",
											},
										},
										OpPos: 54,
										Op:    token.Sub,
										Y: &ast.BasicLit{
											ValPos: 55,
											Kind:   token.IntLit,
											Val:    "42",
										},
									},
								},
								&ast.ExprStmt{
									X: &ast.BinaryExpr{
										X: &ast.BinaryExpr{
											X: &ast.BinaryExpr{
												X: &ast.BinaryExpr{
													X: &ast.UnaryExpr{
														OpPos: 90,
														Op:    token.Not,
														X: &ast.Ident{
															NamePos: 91,
															Name:    "x",
														},
													},
													OpPos: 93,
													Op:    token.Mul,
													Y: &ast.Ident{
														NamePos: 95,
														Name:    "y",
													},
												},
												OpPos: 97,
												Op:    token.Add,
												Y: &ast.Ident{
													NamePos: 99,
													Name:    "z",
												},
											},
											OpPos: 101,
											Op:    token.Lt,
											Y: &ast.Ident{
												NamePos: 103,
												Name:    "x",
											},
										},
										OpPos: 105,
										Op:    token.Ne,
										Y: &ast.BinaryExpr{
											X: &ast.BasicLit{
												ValPos: 108,
												Kind:   token.IntLit,
												Val:    "42",
											},
											OpPos: 111,
											Op:    token.Lt,
											Y: &ast.BinaryExpr{
												X: &ast.Ident{
													NamePos: 113,
													Name:    "x",
												},
												OpPos: 115,
												Op:    token.Add,
												Y: &ast.BinaryExpr{
													X: &ast.Ident{
														NamePos: 117,
														Name:    "y",
													},
													OpPos: 119,
													Op:    token.Mul,
													Y: &ast.UnaryExpr{
														OpPos: 121,
														Op:    token.Not,
														X: &ast.Ident{
															NamePos: 122,
															Name:    "x",
														},
													},
												},
											},
										},
									},
								},
							},
							Rbrace: 164,
						},
					},
				},
			},
		},

		{
			path: "../../testdata/quiet/parser/p05.c",
			want: &ast.File{
				Decls: []ast.Decl{
					&ast.VarDecl{
						VarType: &ast.ArrayType{
							Elem: &ast.Ident{
								NamePos: 0,
								Name:    "int",
							},
							Lbracket: 5,
							Len:      10,
							Rbracket: 8,
						},
						VarName: &ast.Ident{
							NamePos: 4,
							Name:    "c",
						},
					},
					&ast.VarDecl{
						VarType: &ast.ArrayType{
							Elem: &ast.Ident{
								NamePos: 11,
								Name:    "char",
							},
							Lbracket: 17,
							Len:      10,
							Rbracket: 20,
						},
						VarName: &ast.Ident{
							NamePos: 16,
							Name:    "d",
						},
					},
					&ast.FuncDecl{
						FuncType: &ast.FuncType{
							Result: &ast.Ident{
								NamePos: 24,
								Name:    "void",
							},
							Lparen: 30,
							Params: []*ast.Field{
								{
									Type: &ast.ArrayType{
										Elem: &ast.Ident{
											NamePos: 31,
											Name:    "int",
										},
										Lbracket: 36,
										Rbracket: 37,
									},
									Name: &ast.Ident{
										NamePos: 35,
										Name:    "h",
									},
								},
								{
									Type: &ast.ArrayType{
										Elem: &ast.Ident{
											NamePos: 40,
											Name:    "char",
										},
										Lbracket: 46,
										Rbracket: 47,
									},
									Name: &ast.Ident{
										NamePos: 45,
										Name:    "i",
									},
								},
							},
							Rparen: 48,
						},
						FuncName: &ast.Ident{
							NamePos: 29,
							Name:    "f",
						},
						Body: &ast.BlockStmt{
							Lbrace: 50,
							Rbrace: 52,
						},
					},
					&ast.FuncDecl{
						FuncType: &ast.FuncType{
							Result: &ast.Ident{
								NamePos: 56,
								Name:    "int",
							},
							Lparen: 64,
							Params: []*ast.Field{
								{
									Type: &ast.Ident{
										NamePos: 65,
										Name:    "void",
									},
								},
							},
							Rparen: 69,
						},
						FuncName: &ast.Ident{
							NamePos: 60,
							Name:    "main",
						},
						Body: &ast.BlockStmt{
							Lbrace: 71,
							Items: []ast.BlockItem{
								&ast.EmptyStmt{
									Semicolon: 75,
								},
							},
							Rbrace: 77,
						},
					},
				},
			},
		},

		{
			path: "../../testdata/quiet/parser/p06.c",
			want: &ast.File{
				Decls: []ast.Decl{
					&ast.FuncDecl{
						FuncType: &ast.FuncType{
							Result: &ast.Ident{
								NamePos: 0,
								Name:    "void",
							},
							Lparen: 6,
							Params: []*ast.Field{
								{
									Type: &ast.Ident{
										NamePos: 7,
										Name:    "void",
									},
								},
							},
							Rparen: 11,
						},
						FuncName: &ast.Ident{
							NamePos: 5,
							Name:    "f",
						},
						Body: &ast.BlockStmt{
							Lbrace: 13,
							Items: []ast.BlockItem{
								&ast.ReturnStmt{
									Return: 17,
								},
							},
							Rbrace: 25,
						},
					},
					&ast.FuncDecl{
						FuncType: &ast.FuncType{
							Result: &ast.Ident{
								NamePos: 28,
								Name:    "int",
							},
							Lparen: 33,
							Params: []*ast.Field{
								{
									Type: &ast.Ident{
										NamePos: 34,
										Name:    "void",
									},
								},
							},
							Rparen: 38,
						},
						FuncName: &ast.Ident{
							NamePos: 32,
							Name:    "g",
						},
						Body: &ast.BlockStmt{
							Lbrace: 40,
							Items: []ast.BlockItem{
								&ast.ReturnStmt{
									Return: 44,
									Result: &ast.BasicLit{
										ValPos: 51,
										Kind:   token.IntLit,
										Val:    "42",
									},
								},
							},
							Rbrace: 55,
						},
					},
					&ast.FuncDecl{
						FuncType: &ast.FuncType{
							Result: &ast.Ident{
								NamePos: 58,
								Name:    "int",
							},
							Lparen: 66,
							Params: []*ast.Field{
								{
									Type: &ast.Ident{
										NamePos: 67,
										Name:    "void",
									},
								},
							},
							Rparen: 71,
						},
						FuncName: &ast.Ident{
							NamePos: 62,
							Name:    "main",
						},
						Body: &ast.BlockStmt{
							Lbrace: 72,
							Items: []ast.BlockItem{
								&ast.ExprStmt{
									X: &ast.CallExpr{
										Name: &ast.Ident{
											NamePos: 76,
											Name:    "f",
										},
										Lparen: 77,
										Rparen: 78,
									},
								},
								&ast.ExprStmt{
									X: &ast.CallExpr{
										Name: &ast.Ident{
											NamePos: 83,
											Name:    "g",
										},
										Lparen: 84,
										Rparen: 85,
									},
								},
							},
							Rbrace: 88,
						},
					},
				},
			},
		},

		{
			path: "../../testdata/quiet/parser/p07.c",
			want: &ast.File{
				Decls: []ast.Decl{
					&ast.FuncDecl{
						FuncType: &ast.FuncType{
							Result: &ast.Ident{
								NamePos: 0,
								Name:    "int",
							},
							Lparen: 8,
							Params: []*ast.Field{
								{
									Type: &ast.Ident{
										NamePos: 9,
										Name:    "void",
									},
								},
							},
							Rparen: 13,
						},
						FuncName: &ast.Ident{
							NamePos: 4,
							Name:    "main",
						},
						Body: &ast.BlockStmt{
							Lbrace: 14,
							Items: []ast.BlockItem{
								&ast.VarDecl{
									VarType: &ast.Ident{
										NamePos: 18,
										Name:    "int",
									},
									VarName: &ast.Ident{
										NamePos: 22,
										Name:    "x",
									},
								},
								&ast.VarDecl{
									VarType: &ast.Ident{
										NamePos: 27,
										Name:    "int",
									},
									VarName: &ast.Ident{
										NamePos: 31,
										Name:    "y",
									},
								},
								&ast.IfStmt{
									If: 37,
									Cond: &ast.Ident{
										NamePos: 40,
										Name:    "x",
									},
									Body: &ast.WhileStmt{
										While: 43,
										Cond: &ast.Ident{
											NamePos: 50,
											Name:    "y",
										},
										Body: &ast.ExprStmt{
											X: &ast.BinaryExpr{
												X: &ast.Ident{
													NamePos: 53,
													Name:    "x",
												},
												OpPos: 54,
												Op:    token.Assign,
												Y: &ast.BasicLit{
													ValPos: 55,
													Kind:   token.IntLit,
													Val:    "42",
												},
											},
										},
									},
								},
								&ast.WhileStmt{
									While: 64,
									Cond: &ast.Ident{
										NamePos: 70,
										Name:    "x",
									},
									Body: &ast.IfStmt{
										If: 73,
										Cond: &ast.Ident{
											NamePos: 76,
											Name:    "y",
										},
										Body: &ast.ExprStmt{
											X: &ast.BinaryExpr{
												X: &ast.Ident{
													NamePos: 79,
													Name:    "x",
												},
												OpPos: 80,
												Op:    token.Assign,
												Y: &ast.BasicLit{
													ValPos: 81,
													Kind:   token.IntLit,
													Val:    "42",
												},
											},
										},
									},
								},
							},
							Rbrace: 85,
						},
					},
				},
			},
		},

		{
			path: "../../testdata/quiet/parser/p08.c",
			want: &ast.File{
				Decls: []ast.Decl{
					&ast.FuncDecl{
						FuncType: &ast.FuncType{
							Result: &ast.Ident{
								NamePos: 70,
								Name:    "int",
							},
							Lparen: 78,
							Params: []*ast.Field{
								{
									Type: &ast.Ident{
										NamePos: 79,
										Name:    "void",
									},
								},
							},
							Rparen: 83,
						},
						FuncName: &ast.Ident{
							NamePos: 74,
							Name:    "main",
						},
						Body: &ast.BlockStmt{
							Lbrace: 84,
							Items: []ast.BlockItem{
								&ast.VarDecl{
									VarType: &ast.Ident{
										NamePos: 88,
										Name:    "int",
									},
									VarName: &ast.Ident{
										NamePos: 92,
										Name:    "x",
									},
								},
								&ast.VarDecl{
									VarType: &ast.Ident{
										NamePos: 97,
										Name:    "int",
									},
									VarName: &ast.Ident{
										NamePos: 101,
										Name:    "y",
									},
								},
								&ast.IfStmt{
									If: 107,
									Cond: &ast.Ident{
										NamePos: 110,
										Name:    "x",
									},
									Body: &ast.IfStmt{
										If: 118,
										Cond: &ast.Ident{
											NamePos: 122,
											Name:    "y",
										},
										Body: &ast.ExprStmt{
											X: &ast.BinaryExpr{
												X: &ast.Ident{
													NamePos: 125,
													Name:    "x",
												},
												OpPos: 127,
												Op:    token.Assign,
												Y: &ast.BasicLit{
													ValPos: 129,
													Kind:   token.IntLit,
													Val:    "4711",
												},
											},
										},
										Else: &ast.ExprStmt{
											X: &ast.BinaryExpr{
												X: &ast.Ident{
													NamePos: 144,
													Name:    "x",
												},
												OpPos: 145,
												Op:    token.Assign,
												Y: &ast.BasicLit{
													ValPos: 146,
													Kind:   token.IntLit,
													Val:    "42",
												},
											},
										},
									},
								},
							},
							Rbrace: 150,
						},
					},
				},
			},
		},
	}

	for _, g := range golden {
		log.Println("path:", g.path)
		s, err := scanner.Open(g.path)
		if err != nil {
			t.Error(err)
			continue
		}
		p := parser.NewParser()
		file, err := p.Parse(s)
		if err != nil {
			t.Error(err)
			continue
		}
		got := file.(*ast.File)
		if !reflect.DeepEqual(got, g.want) {
			t.Errorf("%q: AST mismatch; expected %#v, got %#v", g.path, g.want, got)
			fmt.Println(pretty.Diff(g.want, got))
		} else {
			fmt.Printf("%q: PASS\n", g.path) // TODO: Remove.
		}
	}
}

func TestParserError(t *testing.T) {
	var golden = []struct {
		path string
		want string
	}{
		{
			path: "../../testdata/incorrect/parser/pe01.c",
			want: `102: unexpected ")", expected ["!" "(" "-" "ident" "int_lit"]`,
		},
		{
			path: "../../testdata/incorrect/parser/pe02.c",
			want: `112: unexpected "}", expected ["!=" "&&" "*" "+" "-" "/" ";" "<" "<=" "=" "==" ">" ">="]`,
		},
		{
			path: "../../testdata/incorrect/parser/pe03.c",
			want: `129: unexpected "}", expected ["!" "(" "-" ";" "ident" "if" "int_lit" "return" "while" "{"]`,
		},
		{
			path: "../../testdata/incorrect/parser/pe04.c",
			want: `111: unexpected "a", expected ["!=" "&&" "(" "*" "+" "-" "/" ";" "<" "<=" "=" "==" ">" ">=" "["]`,
		},
		{
			path: "../../testdata/incorrect/parser/pe05.c",
			want: `71: unexpected "else", expected ["ident"]`,
		},
		{
			path: "../../testdata/incorrect/parser/pe06.c",
			want: `73: unexpected "b", expected ["(" ";" "["]`,
		},
		{
			path: "../../testdata/incorrect/parser/pe07.c",
			want: `72: unexpected ",", expected ["(" ";" "["]`,
		},
		{
			path: "../../testdata/incorrect/parser/pe08.c",
			want: `86: unexpected "42", expected [";" "{"]`,
		},
		{
			// TODO: The ';' at offset 80 in pe09.c shuold probably be a '{', as
			// indicated by the comment "// '}' missing "
			//
			// Update this test case if the test file is fixed.
			path: "../../testdata/incorrect/parser/pe09.c",
			want: `87: unexpected ";", expected ["$" "ident"]`,
		},
		{
			path: "../../testdata/incorrect/parser/pe10.c",
			want: `135: unexpected ")", expected ["!" "(" "-" "ident" "int_lit"]`,
		},
		{
			path: "../../testdata/incorrect/parser/pe11.c",
			want: `70: unexpected "(", expected ["ident"]`,
		},
		{
			path: "../../testdata/incorrect/parser/pe12.c",
			want: `77: unexpected "{", expected ["(" ";" "["]`,
		},
		{
			path: "../../testdata/incorrect/parser/pe13.c",
			want: `126: unexpected ")", expected ["ident"]`,
		},
		{
			path: "../../testdata/incorrect/parser/pe14.c",
			// Note, nested procedures is explicitly allowed by the parser, as the
			// validation is postponed to the semantic analysis checker.
			//
			// References.
			//    https://github.com/mewmew/uc/issues/38
			//    https://github.com/mewmew/uc/issues/43
			want: "",
		},
	}

	for _, g := range golden {
		log.Println("path:", g.path)
		s, err := scanner.Open(g.path)
		if err != nil {
			t.Error(err)
			continue
		}
		p := parser.NewParser()
		_, err = p.Parse(s)
		got := ""
		if err != nil {
			if e, ok := err.(*errors.Error); ok {
				// Unwrap Gocc error.
				err = parser.NewError(e)
			}
			got = err.Error()
		}
		if got != g.want {
			t.Errorf("%q: error mismatch; expected `%v`, got `%v`", g.path, g.want, got)
		}
	}
}

// TODO: add benchmark
