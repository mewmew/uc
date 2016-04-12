package parser_test

import (
	"log"
	"reflect"
	"testing"

	"github.com/kr/pretty"
	"github.com/mewmew/uc/ast"
	"github.com/mewmew/uc/gocc/lexer"
	"github.com/mewmew/uc/gocc/parser"
	"github.com/mewmew/uc/token"
	"github.com/mewmew/uc/types"
)

func TestParser(t *testing.T) {
	var golden = []struct {
		path string
		want *ast.File
	}{
		{
			path: "../../testdata/quiet/parser/p01.c",
			want: &ast.File{
				Decls: []ast.Decl{
					&ast.FuncDecl{
						Type: &types.Func{
							Result: &types.Basic{
								Kind: types.Int,
							},
						},
						Name: &ast.Ident{
							Name: "main",
						},
						Body: &ast.BlockStmt{
							Items: []ast.BlockItem{
								&ast.VarDecl{
									Type: &types.Basic{
										Kind: types.Int,
									},
									Name: &ast.Ident{
										Name: "x",
									},
								},
								&ast.VarDecl{
									Type: &types.Basic{
										Kind: types.Int,
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
		// {path: "../../testdata/quiet/parser/p02.c"},
		// {path: "../../testdata/quiet/parser/p03.c"},
		// {path: "../../testdata/quiet/parser/p04.c"},
		// {path: "../../testdata/quiet/parser/p05.c"},
		// {path: "../../testdata/quiet/parser/p06.c"},
		// {path: "../../testdata/quiet/parser/p07.c"},
		// {path: "../../testdata/quiet/parser/p08.c"},
		// {path: "../../testdata/incorrect/parser/p01.c"},
		// {path: "../../testdata/incorrect/parser/p02.c"},
		// {path: "../../testdata/incorrect/parser/p03.c"},
		// {path: "../../testdata/incorrect/parser/p04.c"},
		// {path: "../../testdata/incorrect/parser/p05.c"},
		// {path: "../../testdata/incorrect/parser/p06.c"},
		// {path: "../../testdata/incorrect/parser/p07.c"},
		// {path: "../../testdata/incorrect/parser/p08.c"},
		// {path: "../../testdata/incorrect/parser/p09.c"},
		// {path: "../../testdata/incorrect/parser/p10.c"},
		// {path: "../../testdata/incorrect/parser/p11.c"},
		// {path: "../../testdata/incorrect/parser/p12.c"},
		// {path: "../../testdata/incorrect/parser/p13.c"},
		// {path: "../../testdata/incorrect/parser/p14.c"},
	}

	for _, g := range golden {
		log.Println("path:", g.path)
		s, err := lexer.NewLexerFile(g.path)
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
			t.Errorf("%q: ast tree mismatch:\nWant: %v\nGot: %v", g.path, g.want, got)
			pretty.Print(g.want)
			pretty.Print(got)
		}
	}
}

// TODO: add benchmark
