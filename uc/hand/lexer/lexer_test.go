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
		{
			path: "../../testdata/incorrect/lexer/good.c",
			toks: []token.Token{
				{
					Kind: token.Comment,
					Val:  "\t\t\t\t\t\t\t||\n**\tFile for testing lexical analysis\t\t||\n**\t\t\t\t\t\t\t||\n**\nThis file would confuse a parser, but\n        is 'lexically correct'.\t\t                ||\n*/\n\n/* ** / ** ",
					Pos:  1,
				},
				{
					Kind: token.Comment,
					Val:  " Simple tokens and single characters:",
					Pos:  179,
				},
				{
					Kind: token.Lbrace,
					Val:  "{",
					Pos:  220,
				},
				{
					Kind: token.Rbrace,
					Val:  "}",
					Pos:  222,
				},
				{
					Kind: token.Comment,
					Val:  " until end-of-line comment",
					Pos:  248,
				},
				{
					Kind: token.Ident,
					Val:  "if",
					Pos:  277,
				},
				{
					Kind: token.Ident,
					Val:  "else",
					Pos:  280,
				},
				{
					Kind: token.Ident,
					Val:  "while",
					Pos:  285,
				},
				{
					Kind: token.Comment,
					Val:  " normal comment ",
					Pos:  311,
				},
				{
					Kind: token.Ident,
					Val:  "return",
					Pos:  332,
				},
				{
					Kind: token.Land,
					Val:  "&&",
					Pos:  339,
				},
				{
					Kind: token.Eq,
					Val:  "==",
					Pos:  342,
				},
				{
					Kind: token.Ne,
					Val:  "!=",
					Pos:  345,
				},
				{
					Kind: token.Le,
					Val:  "<=",
					Pos:  348,
				},
				{
					Kind: token.Ge,
					Val:  ">=",
					Pos:  348,
				},
				{
					Kind: token.Ident,
					Val:  "char",
					Pos:  354,
				},
				{
					Kind: token.Ident,
					Val:  "int",
					Pos:  359,
				},
				{
					Kind: token.Ident,
					Val:  "void",
					Pos:  363,
				},
				{
					Kind: token.Add,
					Val:  "+",
					Pos:  371,
				},
				{
					Kind: token.Sub,
					Val:  "-",
					Pos:  373,
				},
				{
					Kind: token.Mul,
					Val:  "*",
					Pos:  375,
				},
				{
					Kind: token.Div,
					Val:  "/",
					Pos:  377,
				},
				{
					Kind: token.Lt,
					Val:  "<",
					Pos:  380,
				},
				{
					Kind: token.Gt,
					Val:  ">",
					Pos:  382,
				},
				{
					Kind: token.Assign,
					Val:  "=",
					Pos:  384,
				},
				{
					Kind: token.Lparen,
					Val:  "(",
					Pos:  387,
				},
				{
					Kind: token.Comma,
					Val:  ",",
					Pos:  388,
				},
				{
					Kind: token.Semicolon,
					Val:  ";",
					Pos:  389,
				},
				{
					Kind: token.Rparen,
					Val:  ")",
					Pos:  390,
				},
				{
					Kind: token.Lbracket,
					Val:  "[",
					Pos:  392,
				},
				{
					Kind: token.Ident,
					Val:  "was",
					Pos:  393,
				},
				{
					Kind: token.Sub,
					Val:  "-",
					Pos:  396,
				},
				{
					Kind: token.Ident,
					Val:  "colon",
					Pos:  397,
				},
				{
					Kind: token.Rbracket,
					Val:  "]",
					Pos:  402,
				},
				{
					Kind: token.Comment,
					Val:  " Comment with bad tokens: _ || | ++ # @ ...  ",
					Pos:  406,
				},
				{
					Kind: token.Comment,
					Val:  " Ditto */ /* : _ || | ++ # @ ...  ",
					Pos:  456,
				},
				{
					Kind: token.Comment,
					Val:  " Identifiers and numbers:",
					Pos:  493,
				},
				{
					Kind: token.IntLit,
					Val:  "17",
					Pos:  522,
				},
				{
					Kind: token.Sub,
					Val:  "-",
					Pos:  525,
				},
				{
					Kind: token.IntLit,
					Val:  "17",
					Pos:  526,
				},
				{
					Kind: token.Comment,
					Val:  " No floats? -17.17e17 -17.17E-17  ",
					Pos:  529,
				},
				{
					Kind: token.Ident,
					Val:  "ponderosa",
					Pos:  568,
				},
				{
					Kind: token.Ident,
					Val:  "Black",
					Pos:  578,
				},
				{
					Kind: token.Ident,
					Val:  "Steel",
					Pos:  584,
				},
				{
					Kind: token.Ident,
					Val:  "PUMPKIN",
					Pos:  590,
				},
				{
					Kind: token.Ident,
					Val:  "AfterMath",
					Pos:  598,
				},
				{
					Kind: token.Ident,
					Val:  "aBBaoN",
					Pos:  608,
				},
				{
					Kind: token.Ident,
					Val:  "faT",
					Pos:  615,
				},
				{
					Kind: token.Ident,
					Val:  "TRacKs",
					Pos:  619,
				},
				{
					Kind: token.Ident,
					Val:  "K9",
					Pos:  628,
				},
				{
					Kind: token.Ident,
					Val:  "K23",
					Pos:  631,
				},
				{
					Kind: token.Ident,
					Val:  "B52",
					Pos:  635,
				},
				{
					Kind: token.Ident,
					Val:  "Track15",
					Pos:  639,
				},
				{
					Kind: token.Ident,
					Val:  "not4money",
					Pos:  647,
				},
				{
					Kind: token.IntLit,
					Val:  "378",
					Pos:  657,
				},
				{
					Kind: token.Ident,
					Val:  "WHOIS666999SIOHM",
					Pos:  661,
				},
				{
					Kind: token.Ident,
					Val:  "was",
					Pos:  687,
				},
				{
					Kind: token.Sub,
					Val:  "-",
					Pos:  690,
				},
				{
					Kind: token.Ident,
					Val:  "floating",
					Pos:  691,
				},
				{
					Kind: token.Sub,
					Val:  "-",
					Pos:  699,
				},
				{
					Kind: token.Ident,
					Val:  "point",
					Pos:  700,
				},
				{
					Kind: token.Sub,
					Val:  "-",
					Pos:  705,
				},
				{
					Kind: token.Ident,
					Val:  "number",
					Pos:  706,
				},
				{
					Kind: token.Comment,
					Val:  " The following 'trap' should be correctly handled:\n\n\t\t* \"2die4U\" consists of the number '2' and the\n\t\t  identifier 'die4U'.\n",
					Pos:  714,
				},
				{
					Kind: token.IntLit,
					Val:  "2",
					Pos:  852,
				},
				{
					Kind: token.Ident,
					Val:  "die4U",
					Pos:  853,
				},
				{
					Kind: token.Comment,
					Val:  "|| The following should all be regarded as identifiers:",
					Pos:  860,
				},
				{
					Kind: token.Ident,
					Val:  "Function",
					Pos:  920,
				},
				{
					Kind: token.Ident,
					Val:  "PrOceDuRE",
					Pos:  929,
				},
				{
					Kind: token.Ident,
					Val:  "begIN",
					Pos:  939,
				},
				{
					Kind: token.Ident,
					Val:  "eNd",
					Pos:  945,
				},
				{
					Kind: token.Ident,
					Val:  "PrinT",
					Pos:  949,
				},
				{
					Kind: token.Ident,
					Val:  "rEad",
					Pos:  955,
				},
				{
					Kind: token.Ident,
					Val:  "iF",
					Pos:  960,
				},
				{
					Kind: token.Ident,
					Val:  "THen",
					Pos:  963,
				},
				{
					Kind: token.Ident,
					Val:  "StaTic",
					Pos:  968,
				},
				{
					Kind: token.Ident,
					Val:  "ElSe",
					Pos:  976,
				},
				{
					Kind: token.Ident,
					Val:  "wHilE",
					Pos:  981,
				},
				{
					Kind: token.Ident,
					Val:  "Do",
					Pos:  987,
				},
				{
					Kind: token.Ident,
					Val:  "reTurN",
					Pos:  990,
				},
				{
					Kind: token.Ident,
					Val:  "noT",
					Pos:  997,
				},
				{
					Kind: token.Ident,
					Val:  "AnD",
					Pos:  1001,
				},
				{
					Kind: token.Ident,
					Val:  "OR",
					Pos:  1005,
				},
				{
					Kind: token.Ident,
					Val:  "TrUE",
					Pos:  1008,
				},
				{
					Kind: token.Ident,
					Val:  "bOOl",
					Pos:  1013,
				},
				{
					Kind: token.Ident,
					Val:  "FalsE",
					Pos:  1018,
				},
				{
					Kind: token.Ident,
					Val:  "sizE",
					Pos:  1024,
				},
				{
					Kind: token.Comment,
					Val:  " It is legal to end the code like this, without an ending newline.",
					Pos:  1031,
				},
				{
					Kind: token.EOF,
					Val:  "",
					// TODO: Figure out how to handle the offset in error cases.
					Pos: 1100,
				},
			},
		},
		// TODO: Add tokens for the following test cases.
		{path: "../../testdata/incorrect/lexer/long-char.c"},
		{path: "../../testdata/incorrect/lexer/ugly.c"},
		{path: "../../testdata/quiet/lexer/l01.c"},
		{path: "../../testdata/quiet/lexer/l02.c"},
		{path: "../../testdata/quiet/lexer/l03.c"},
		{path: "../../testdata/quiet/lexer/l04.c"},
		{path: "../../testdata/quiet/lexer/l05.c"},
		{path: "../../testdata/quiet/lexer/l06.c"},
	}

	for looprun, g := range golden {
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
		if looprun >= 1 {
			break // TODO: Remove this break to test all test cases and not just the first.
		}
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
