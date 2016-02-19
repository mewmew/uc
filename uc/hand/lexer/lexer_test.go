package lexer_test

import (
	"io/ioutil"
	"log"
	"testing"

	"github.com/mewmew/uc/uc/hand/lexer"
	"github.com/mewmew/uc/uc/hand/token"
)

func TestLexer(t *testing.T) {
	var golden = []struct {
		path string
		toks []token.Token
	}{
		{
			path: "../../testdata/incorrect/lexer/bad.c",
			toks: []token.Token{
				{
					Kind: token.Comment,
					Val:  "\t\t\t\t\t\t\t||\n**\tFile for testing lexical analysis\t\t||\n**\t\t\t\t\t\t\t||\n**\tThis file is 'lexically incorrect'.\t\t||\n",
					Pos:  0,
				},
				{
					Kind: token.Ident,
					Val:  "int",
					Pos:  114,
				},
				{
					Kind: token.Ident,
					Val:  "main",
					Pos:  118,
				},
				{
					Kind: token.Lparen,
					Val:  "(",
					Pos:  122,
				},
				{
					Kind: token.Ident,
					Val:  "void",
					Pos:  123,
				},
				{
					Kind: token.Rparen,
					Val:  ")",
					Pos:  127,
				},
				{
					Kind: token.Lbrace,
					Val:  "{",
					Pos:  129,
				},
				{
					Kind: token.IntLit,
					Val:  "42",
					Pos:  133,
				},
				{
					Kind: token.Semicolon,
					Val:  ";",
					Pos:  135,
				},
				{
					Kind: token.Rbrace,
					Val:  "}",
					Pos:  137,
				},
				{
					Kind: token.Error,
					Val:  "unexpected eof in block comment",
					// TODO: Figure out how to handle the offset in error cases.
					Pos: 220,
				},
				{
					Kind: token.EOF,
					Val:  "",
					// TODO: Figure out how to handle the offset in error cases.
					Pos: 221,
				},
			},
		},
		// TODO: Add tokens for the following test cases.
		{path: "../../testdata/incorrect/lexer/good.c"},
		{path: "../../testdata/incorrect/lexer/long-char.c"},
		{path: "../../testdata/incorrect/lexer/ugly.c"},
		{path: "../../testdata/quiet/lexer/l01.c"},
		{path: "../../testdata/quiet/lexer/l02.c"},
		{path: "../../testdata/quiet/lexer/l03.c"},
		{path: "../../testdata/quiet/lexer/l04.c"},
		{path: "../../testdata/quiet/lexer/l05.c"},
		{path: "../../testdata/quiet/lexer/l06.c"},
	}

	for _, g := range golden {
		log.Println("path:", g.path)
		tokens, err := lexer.ParseFile(g.path)
		if err != nil {
			t.Error(err)
			continue
		}
		for j := 0; ; j++ {
			if j >= len(tokens) {
				t.Errorf("%s: invalid number of tokens; expected %d tokens, got %d", g.path, len(tokens), j)
				continue
			}
			got := tokens[j]
			if j >= len(g.toks) {
				t.Errorf("%s: invalid number of tokens; expected %d tokens, got %d", g.path, len(g.toks), j)
				continue
			}
			if want := g.toks[j]; got != want {
				t.Errorf("%s: token mismatch; expected %#v, got %#v", g.path, want, got)
			}
			if got.Kind == token.EOF {
				if j != len(g.toks)-1 {
					t.Errorf("%s: invalid number of tokens; expected %d tokens, got %d", g.path, len(g.toks), j)
				}
				break
			}
		}
		break // TODO: Remove this break to test all test cases and not just the first.
	}
}

func BenchmarkLexer(b *testing.B) {
	buf, err := ioutil.ReadFile("../../testdata/noisy/advanced/eval.c")
	if err != nil {
		b.Fatal(err)
	}
	src := string(buf)
	for i := 0; i < b.N; i++ {
		tokens := lexer.ParseString(src)
		for _, tok := range tokens {
			if tok.Kind == token.EOF {
				break
			}
		}
	}
}
