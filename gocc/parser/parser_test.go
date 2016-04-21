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
						Type: &ast.FuncType{
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
						Name: &ast.Ident{
							NamePos: 4,
							Name:    "main",
						},
						Body: &ast.BlockStmt{
							Lbrace: 15,
							Items: []ast.BlockItem{
								&ast.VarDecl{
									Type: &ast.Ident{
										NamePos: 19,
										Name:    "int",
									},
									Name: &ast.Ident{
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
						Type: &ast.FuncType{
							Result: &ast.Ident{
								Name: "int",
							},
							Params: []*ast.Field{
								{
									Type: &ast.Ident{
										Name: "void",
									},
								},
							},
						},
						Name: &ast.Ident{
							Name: "main",
						},
						Body: &ast.BlockStmt{
							Items: []ast.BlockItem{
								&ast.VarDecl{
									Type: &ast.Ident{
										Name: "int",
									},
									Name: &ast.Ident{
										Name: "x",
									},
								},
								&ast.VarDecl{
									Type: &ast.Ident{
										Name: "int",
									},
									Name: &ast.Ident{
										Name: "y",
									},
								},
								&ast.ExprStmt{
									X: &ast.BinaryExpr{
										X: &ast.Ident{
											Name: "x",
										},
										Op: token.Assign,
										Y: &ast.BasicLit{
											Kind: token.IntLit,
											Val:  "42",
										},
									},
								},
								&ast.ExprStmt{
									X: &ast.BinaryExpr{
										X: &ast.Ident{
											Name: "x",
										},
										Op: token.Assign,
										Y: &ast.BinaryExpr{
											X: &ast.Ident{
												Name: "y",
											},
											Op: token.Assign,
											Y: &ast.BasicLit{
												Kind: token.IntLit,
												Val:  "4711",
											},
										},
									},
								},
							},
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
						Type: &ast.FuncType{
							Result: &ast.Ident{
								Name: "int",
							},
							Params: []*ast.Field{
								{
									Type: &ast.Ident{
										Name: "void",
									},
								},
							},
						},
						Name: &ast.Ident{
							Name: "main",
						},
						Body: &ast.BlockStmt{
							Items: []ast.BlockItem{
								&ast.VarDecl{
									Type: &ast.Ident{
										Name: "int",
									},
									Name: &ast.Ident{
										Name: "x",
									},
								},
								&ast.EmptyStmt{},
								&ast.WhileStmt{
									Cond: &ast.BinaryExpr{
										X: &ast.Ident{
											Name: "x",
										},
										Op: token.Lt,
										Y: &ast.BasicLit{
											Kind: token.IntLit,
											Val:  "10",
										},
									},
									Body: &ast.ExprStmt{
										X: &ast.BinaryExpr{
											X: &ast.Ident{
												Name: "x",
											},
											Op: token.Assign,
											Y: &ast.BinaryExpr{
												X: &ast.Ident{
													Name: "x",
												},
												Op: token.Add,
												Y: &ast.BasicLit{
													Kind: token.IntLit,
													Val:  "3",
												},
											},
										},
									},
								},
								&ast.IfStmt{
									Cond: &ast.BasicLit{
										Kind: token.IntLit,
										Val:  "1",
									},
									Body: &ast.ExprStmt{
										X: &ast.BinaryExpr{
											X: &ast.Ident{
												Name: "x",
											},
											Op: token.Assign,
											Y: &ast.BinaryExpr{
												X: &ast.Ident{
													Name: "x",
												},
												Op: token.Add,
												Y: &ast.BasicLit{
													Kind: token.IntLit,
													Val:  "3",
												},
											},
										},
									},
								},
							},
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
						Type: &ast.FuncType{
							Result: &ast.Ident{
								Name: "int",
							},
							Params: []*ast.Field{
								{
									Type: &ast.Ident{
										Name: "void",
									},
								},
							},
						},
						Name: &ast.Ident{
							Name: "main",
						},
						Body: &ast.BlockStmt{
							Items: []ast.BlockItem{
								&ast.VarDecl{
									Type: &ast.Ident{
										Name: "int",
									},
									Name: &ast.Ident{
										Name: "x",
									},
								},
								&ast.IfStmt{
									Cond: &ast.BinaryExpr{
										X: &ast.BasicLit{
											Kind: token.IntLit,
											Val:  "1",
										},
										Op: token.Lt,
										Y: &ast.BasicLit{
											Kind: token.IntLit,
											Val:  "2",
										},
									},
									Body: &ast.ExprStmt{
										X: &ast.BinaryExpr{
											X: &ast.Ident{
												Name: "x",
											},
											Op: token.Assign,
											Y: &ast.BasicLit{
												Kind: token.IntLit,
												Val:  "1",
											},
										},
									},
									Else: &ast.ExprStmt{
										X: &ast.BinaryExpr{
											X: &ast.Ident{
												Name: "x",
											},
											Op: token.Assign,
											Y: &ast.BasicLit{
												Kind: token.IntLit,
												Val:  "2",
											},
										},
									},
								},
							},
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
						Type: &ast.FuncType{
							Result: &ast.Ident{
								Name: "int",
							},
							Params: []*ast.Field{
								{
									Type: &ast.Ident{
										Name: "void",
									},
								},
							},
						},
						Name: &ast.Ident{
							Name: "main",
						},
						Body: &ast.BlockStmt{
							Items: []ast.BlockItem{
								&ast.VarDecl{
									Type: &ast.Ident{
										Name: "int",
									},
									Name: &ast.Ident{
										Name: "x",
									},
								},
								&ast.VarDecl{
									Type: &ast.Ident{
										Name: "int",
									},
									Name: &ast.Ident{
										Name: "y",
									},
								},
								&ast.VarDecl{
									Type: &ast.Ident{
										Name: "int",
									},
									Name: &ast.Ident{
										Name: "z",
									},
								},
								&ast.ExprStmt{
									X: &ast.BinaryExpr{
										X: &ast.BinaryExpr{
											X: &ast.BinaryExpr{
												X: &ast.Ident{
													Name: "x",
												},
												Op: token.Sub,
												Y: &ast.Ident{
													Name: "y",
												},
											},
											Op: token.Sub,
											Y: &ast.Ident{
												Name: "z",
											},
										},
										Op: token.Sub,
										Y: &ast.BasicLit{
											Kind: token.IntLit,
											Val:  "42",
										},
									},
								},
								&ast.ExprStmt{
									X: &ast.BinaryExpr{
										X: &ast.BinaryExpr{
											X: &ast.BinaryExpr{
												X: &ast.BinaryExpr{
													X: &ast.UnaryExpr{
														Op: token.Not,
														X: &ast.Ident{
															Name: "x",
														},
													},
													Op: token.Mul,
													Y: &ast.Ident{
														Name: "y",
													},
												},
												Op: token.Add,
												Y: &ast.Ident{
													Name: "z",
												},
											},
											Op: token.Lt,
											Y: &ast.Ident{
												Name: "x",
											},
										},
										Op: token.Ne,
										Y: &ast.BinaryExpr{
											X: &ast.BasicLit{
												Kind: token.IntLit,
												Val:  "42",
											},
											Op: token.Lt,
											Y: &ast.BinaryExpr{
												X: &ast.Ident{
													Name: "x",
												},
												Op: token.Add,
												Y: &ast.BinaryExpr{
													X: &ast.Ident{
														Name: "y",
													},
													Op: token.Mul,
													Y: &ast.UnaryExpr{
														Op: token.Not,
														X: &ast.Ident{
															Name: "x",
														},
													},
												},
											},
										},
									},
								},
							},
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
						Type: &ast.ArrayType{
							Elem: &ast.Ident{
								Name: "int",
							},
							Len: 10,
						},
						Name: &ast.Ident{
							Name: "c",
						},
					},
					&ast.VarDecl{
						Type: &ast.ArrayType{
							Elem: &ast.Ident{
								Name: "char",
							},
							Len: 10,
						},
						Name: &ast.Ident{
							Name: "d",
						},
					},
					&ast.FuncDecl{
						Type: &ast.FuncType{
							Result: &ast.Ident{
								Name: "void",
							},
							Params: []*ast.Field{
								{
									Type: &ast.ArrayType{
										Elem: &ast.Ident{
											Name: "int",
										},
									},
									Name: &ast.Ident{
										Name: "h",
									},
								},
								{
									Type: &ast.ArrayType{
										Elem: &ast.Ident{
											Name: "char",
										},
									},
									Name: &ast.Ident{
										Name: "i",
									},
								},
							},
						},
						Name: &ast.Ident{
							Name: "f",
						},
						Body: &ast.BlockStmt{},
					},
					&ast.FuncDecl{
						Type: &ast.FuncType{
							Result: &ast.Ident{
								Name: "int",
							},
							Params: []*ast.Field{
								{
									Type: &ast.Ident{
										Name: "void",
									},
								},
							},
						},
						Name: &ast.Ident{
							Name: "main",
						},
						Body: &ast.BlockStmt{
							Items: []ast.BlockItem{
								&ast.EmptyStmt{},
							},
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
						Type: &ast.FuncType{
							Result: &ast.Ident{
								Name: "void",
							},
							Params: []*ast.Field{
								{
									Type: &ast.Ident{
										Name: "void",
									},
								},
							},
						},
						Name: &ast.Ident{
							Name: "f",
						},
						Body: &ast.BlockStmt{
							Items: []ast.BlockItem{
								&ast.ReturnStmt{},
							},
						},
					},
					&ast.FuncDecl{
						Type: &ast.FuncType{
							Result: &ast.Ident{
								Name: "int",
							},
							Params: []*ast.Field{
								{
									Type: &ast.Ident{
										Name: "void",
									},
								},
							},
						},
						Name: &ast.Ident{
							Name: "g",
						},
						Body: &ast.BlockStmt{
							Items: []ast.BlockItem{
								&ast.ReturnStmt{
									Result: &ast.BasicLit{
										Kind: token.IntLit,
										Val:  "42",
									},
								},
							},
						},
					},
					&ast.FuncDecl{
						Type: &ast.FuncType{
							Result: &ast.Ident{
								Name: "int",
							},
							Params: []*ast.Field{
								{
									Type: &ast.Ident{
										Name: "void",
									},
								},
							},
						},
						Name: &ast.Ident{
							Name: "main",
						},
						Body: &ast.BlockStmt{
							Items: []ast.BlockItem{
								&ast.ExprStmt{
									X: &ast.CallExpr{
										Name: &ast.Ident{
											Name: "f",
										},
									},
								},
								&ast.ExprStmt{
									X: &ast.CallExpr{
										Name: &ast.Ident{
											Name: "g",
										},
									},
								},
							},
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
						Type: &ast.FuncType{
							Result: &ast.Ident{
								Name: "int",
							},
							Params: []*ast.Field{
								{
									Type: &ast.Ident{
										Name: "void",
									},
								},
							},
						},
						Name: &ast.Ident{
							Name: "main",
						},
						Body: &ast.BlockStmt{
							Items: []ast.BlockItem{
								&ast.VarDecl{
									Type: &ast.Ident{
										Name: "int",
									},
									Name: &ast.Ident{
										Name: "x",
									},
								},
								&ast.VarDecl{
									Type: &ast.Ident{
										Name: "int",
									},
									Name: &ast.Ident{
										Name: "y",
									},
								},
								&ast.IfStmt{
									Cond: &ast.Ident{
										Name: "x",
									},
									Body: &ast.WhileStmt{
										Cond: &ast.Ident{
											Name: "y",
										},
										Body: &ast.ExprStmt{
											X: &ast.BinaryExpr{
												X: &ast.Ident{
													Name: "x",
												},
												Op: token.Assign,
												Y: &ast.BasicLit{
													Kind: token.IntLit,
													Val:  "42",
												},
											},
										},
									},
								},
								&ast.WhileStmt{
									Cond: &ast.Ident{
										Name: "x",
									},
									Body: &ast.IfStmt{
										Cond: &ast.Ident{
											Name: "y",
										},
										Body: &ast.ExprStmt{
											X: &ast.BinaryExpr{
												X: &ast.Ident{
													Name: "x",
												},
												Op: token.Assign,
												Y: &ast.BasicLit{
													Kind: token.IntLit,
													Val:  "42",
												},
											},
										},
									},
								},
							},
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
						Type: &ast.FuncType{
							Result: &ast.Ident{
								Name: "int",
							},
							Params: []*ast.Field{
								{
									Type: &ast.Ident{
										Name: "void",
									},
								},
							},
						},
						Name: &ast.Ident{
							Name: "main",
						},
						Body: &ast.BlockStmt{
							Items: []ast.BlockItem{
								&ast.VarDecl{
									Type: &ast.Ident{
										Name: "int",
									},
									Name: &ast.Ident{
										Name: "x",
									},
								},
								&ast.VarDecl{
									Type: &ast.Ident{
										Name: "int",
									},
									Name: &ast.Ident{
										Name: "y",
									},
								},
								&ast.IfStmt{
									Cond: &ast.Ident{
										Name: "x",
									},
									Body: &ast.IfStmt{
										Cond: &ast.Ident{
											Name: "y",
										},
										Body: &ast.ExprStmt{
											X: &ast.BinaryExpr{
												X: &ast.Ident{
													Name: "x",
												},
												Op: token.Assign,
												Y: &ast.BasicLit{
													Kind: token.IntLit,
													Val:  "4711",
												},
											},
										},
										Else: &ast.ExprStmt{
											X: &ast.BinaryExpr{
												X: &ast.Ident{
													Name: "x",
												},
												Op: token.Assign,
												Y: &ast.BasicLit{
													Kind: token.IntLit,
													Val:  "42",
												},
											},
										},
									},
								},
							},
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
