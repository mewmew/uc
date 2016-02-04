package lexer_test

import (
	"bytes"
	"log"
	"testing"

	"github.com/mewmew/uc/uc/grammar/lexer"
	"github.com/mewmew/uc/uc/grammar/token"
)

func TestLexer(t *testing.T) {
	var golden = []struct {
		path string
		toks []*token.Token
	}{
		{
			path: "../testdata/incorrect/lexer/bad.c",
			toks: []*token.Token{
				{
					Type: token.TokMap.Type("comment"),
					Lit:  []byte("/*\t\t\t\t\t\t\t||\n**\tFile for testing lexical analysis\t\t||\n**\t\t\t\t\t\t\t||\n**\tThis file is 'lexically incorrect'.\t\t||\n*/"),
					Pos:  token.Pos{Offset: 0},
				},
				{
					Type: token.TokMap.Type("int"),
					Lit:  []byte("int"),
					Pos:  token.Pos{Offset: 114},
				},
				{
					Type: token.TokMap.Type("ident"),
					Lit:  []byte("main"),
					Pos:  token.Pos{Offset: 118},
				},
				{
					Type: token.TokMap.Type("("),
					Lit:  []byte("("),
					Pos:  token.Pos{Offset: 122},
				},
				{
					Type: token.TokMap.Type("void"),
					Lit:  []byte("void"),
					Pos:  token.Pos{Offset: 123},
				},
				{
					Type: token.TokMap.Type(")"),
					Lit:  []byte(")"),
					Pos:  token.Pos{Offset: 127},
				},
				{
					Type: token.TokMap.Type("{"),
					Lit:  []byte("{"),
					Pos:  token.Pos{Offset: 129},
				},
				{
					Type: token.TokMap.Type("int_lit"),
					Lit:  []byte("42"),
					Pos:  token.Pos{Offset: 133},
				},
				{
					Type: token.TokMap.Type(";"),
					Lit:  []byte(";"),
					Pos:  token.Pos{Offset: 135},
				},
				{
					Type: token.TokMap.Type("}"),
					Lit:  []byte("}"),
					Pos:  token.Pos{Offset: 137},
				},
				{
					Type: token.TokMap.Type("INVALID"),
					Lit:  []byte("/*\n\tIt is not legal to end the code like this, \n\twithout a comment terminator\n"),
					Pos:  token.Pos{Offset: 143},
				},
				{
					Type: token.TokMap.Type("$"),
					Lit:  nil,
					Pos:  token.Pos{Offset: 221},
				},
			},
		},
		// TODO: Add tokens for the following test cases.
		{path: "../testdata/incorrect/lexer/good.c"},
		{path: "../testdata/incorrect/lexer/long-char.c"},
		{path: "../testdata/incorrect/lexer/ugly.c"},
		{path: "../testdata/quiet/lexer/l01.c"},
		{path: "../testdata/quiet/lexer/l02.c"},
		{path: "../testdata/quiet/lexer/l03.c"},
		{path: "../testdata/quiet/lexer/l04.c"},
		{path: "../testdata/quiet/lexer/l05.c"},
		{path: "../testdata/quiet/lexer/l06.c"},
	}

	for i, g := range golden {
		log.Println("path:", g.path)
		s, err := lexer.NewLexerFile(g.path)
		if err != nil {
			t.Error(err)
			continue
		}
		for j := 0; ; j++ {
			got := s.Scan()
			if j >= len(g.toks) {
				t.Errorf("i=%d: invalid number of tokens; expected %d tokens, got %d", i, len(g.toks), j)
				continue
			}
			if want := g.toks[j]; !tokenEqual(got, want) {
				t.Errorf("i=%d: token mismatch; expected %#v, got %#v", i, want, got)
			}
			if got.Type == token.EOF {
				if j != len(g.toks)-1 {
					t.Errorf("i=%d: invalid number of tokens; expected %d tokens, got %d", i, len(g.toks), j)
				}
				break
			}
		}
		break // TODO: Remove this break to test all test cases and not just the first.
	}
}

// tokenEqual reports whether the given tokens are equal.
func tokenEqual(t1, t2 *token.Token) bool {
	return bytes.Equal(t1.Lit, t2.Lit) && t1.Type == t2.Type && t1.Pos.Offset == t2.Pos.Offset
}
