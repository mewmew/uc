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
					Val:  "\t\t\t\t\t\t\t||\n**\tFile for testing lexical analysis\t\t||\n**\t\t\t\t\t\t\t||\n**\tThis file would confuse a parser, but\n        is 'lexically correct'.\t\t                ||\n",
					Pos:  1,
				},
				{
					Kind: token.Comment,
					Val:  " ** / ** ",
					Pos:  163,
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
					Kind: token.KwIf,
					Val:  "if",
					Pos:  277,
				},
				{
					Kind: token.KwElse,
					Val:  "else",
					Pos:  280,
				},
				{
					Kind: token.KwWhile,
					Val:  "while",
					Pos:  285,
				},
				{
					Kind: token.Comment,
					Val:  " normal comment ",
					Pos:  311,
				},
				{
					Kind: token.KwReturn,
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
					Pos:  351,
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
					Val:  "R23",
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
		{
			path: "../../testdata/incorrect/lexer/long-char.c",
			toks: []token.Token{
				{
					Kind: token.Ident,
					Val:  "int",
					Pos:  0,
				},
				{
					Kind: token.Ident,
					Val:  "main",
					Pos:  4,
				},
				{
					Kind: token.Lparen,
					Val:  "(",
					Pos:  8,
				},
				{
					Kind: token.Ident,
					Val:  "void",
					Pos:  9,
				},
				{
					Kind: token.Rparen,
					Val:  ")",
					Pos:  13,
				},
				{
					Kind: token.Lbrace,
					Val:  "{",
					Pos:  15,
				},
				{
					Kind: token.Ident,
					Val:  "char",
					Pos:  19,
				},
				{
					Kind: token.Ident,
					Val:  "c",
					Pos:  24,
				},
				{
					Kind: token.Semicolon,
					Val:  ";",
					Pos:  25,
				},
				{
					Kind: token.Ident,
					Val:  "c",
					Pos:  29,
				},
				{
					Kind: token.Assign,
					Val:  "=",
					Pos:  31,
				},
				{
					Kind: token.CharLit,
					Val:  "'c'",
					Pos:  33,
				},
				{
					Kind: token.Semicolon,
					Val:  ";",
					Pos:  36,
				},
				{
					Kind: token.Comment,
					Val:  " OK",
					Pos:  38,
				},
				{
					Kind: token.Ident,
					Val:  "c",
					Pos:  46,
				},
				{
					Kind: token.Assign,
					Val:  "=",
					Pos:  48,
				},
				{
					Kind: token.Error,
					Val:  "unterminated character literal",
					// TODO: Figure out how to handle position of errors.
					Pos: 51,
				},
				{
					Kind: token.Ident,
					Val:  "c",
					Pos:  52,
				},
				{
					Kind: token.Error,
					Val:  "unterminated character literal",
					// TODO: Figure out how to handle position of errors.
					Pos: 54,
				},
				{
					Kind: token.Comment,
					Val:  " Not OK",
					Pos:  56,
				},
				{
					Kind: token.Rbrace,
					Val:  "}",
					Pos:  66,
				},
				{
					Kind: token.EOF,
					Val:  "",
					Pos:  68,
				},
			},
		},
		// TODO: Add tokens for the following test cases.
		//{path: "../../testdata/incorrect/lexer/ugly.c"},
		{
			path: "../../testdata/quiet/lexer/l01.c",
			toks: []token.Token{
				{
					Kind: token.Ident,
					Val:  "int",
					Pos:  0,
				},
				{
					Kind: token.Ident,
					Val:  "main",
					Pos:  4,
				},
				{
					Kind: token.Lparen,
					Val:  "(",
					Pos:  9,
				},
				{
					Kind: token.Ident,
					Val:  "void",
					Pos:  10,
				},
				{
					Kind: token.Rparen,
					Val:  ")",
					Pos:  14,
				},
				{
					Kind: token.Lbrace,
					Val:  "{",
					Pos:  16,
				},
				{
					Kind: token.Semicolon,
					Val:  ";",
					Pos:  18,
				},
				{
					Kind: token.Rbrace,
					Val:  "}",
					Pos:  20,
				},
				{
					Kind: token.EOF,
					Val:  "",
					Pos:  23,
				},
			},
		},
		{
			path: "../../testdata/quiet/lexer/l02.c",
			toks: []token.Token{
				{
					Kind: token.Ident,
					Val:  "int",
					Pos:  0,
				},
				{
					Kind: token.Ident,
					Val:  "foo",
					Pos:  4,
				},
				{
					Kind: token.Semicolon,
					Val:  ";",
					Pos:  7,
				},
				{
					Kind: token.Ident,
					Val:  "int",
					Pos:  10,
				},
				{
					Kind: token.Ident,
					Val:  "BarBara",
					Pos:  14,
				},
				{
					Kind: token.Semicolon,
					Val:  ";",
					Pos:  21,
				},
				{
					Kind: token.Ident,
					Val:  "int",
					Pos:  24,
				},
				{
					Kind: token.Ident,
					Val:  "bar_bara",
					Pos:  28,
				},
				{
					Kind: token.Semicolon,
					Val:  ";",
					Pos:  36,
				},
				{
					Kind: token.Ident,
					Val:  "int",
					Pos:  39,
				},
				{
					Kind: token.Ident,
					Val:  "bar4711",
					Pos:  43,
				},
				{
					Kind: token.Semicolon,
					Val:  ";",
					Pos:  50,
				},
				{
					Kind: token.Ident,
					Val:  "int",
					Pos:  53,
				},
				{
					Kind: token.Ident,
					Val:  "b4rb4r4",
					Pos:  57,
				},
				{
					Kind: token.Semicolon,
					Val:  ";",
					Pos:  64,
				},
				{
					Kind: token.Ident,
					Val:  "int",
					Pos:  67,
				},
				{
					Kind: token.Ident,
					Val:  "abcdefghijklmnopqrstuvwxyz_ABCDEFGHIJKLMNOPQRSTUVWXYZ_0123456789",
					Pos:  71,
				},
				{
					Kind: token.Semicolon,
					Val:  ";",
					Pos:  135,
				},
				{
					Kind: token.Ident,
					Val:  "int",
					Pos:  138,
				},
				{
					Kind: token.Ident,
					Val:  "main",
					Pos:  142,
				},
				{
					Kind: token.Lparen,
					Val:  "(",
					Pos:  146,
				},
				{
					Kind: token.Ident,
					Val:  "void",
					Pos:  147,
				},
				{
					Kind: token.Rparen,
					Val:  ")",
					Pos:  151,
				},
				{
					Kind: token.Lbrace,
					Val:  "{",
					Pos:  153,
				},
				{
					Kind: token.Semicolon,
					Val:  ";",
					Pos:  155,
				},
				{
					Kind: token.Rbrace,
					Val:  "}",
					Pos:  157,
				},
				{
					Kind: token.EOF,
					Val:  "",
					Pos:  159,
				},
			},
		},
		{
			path: "../../testdata/quiet/lexer/l03.c",
			toks: []token.Token{
				{
					Kind: token.Ident,
					Val:  "int",
					Pos:  0,
				},
				{
					Kind: token.Ident,
					Val:  "main",
					Pos:  4,
				},
				{
					Kind: token.Lparen,
					Val:  "(",
					Pos:  8,
				},
				{
					Kind: token.Ident,
					Val:  "void",
					Pos:  9,
				},
				{
					Kind: token.Rparen,
					Val:  ")",
					Pos:  13,
				},
				{
					Kind: token.Lbrace,
					Val:  "{",
					Pos:  15,
				},
				{
					Kind: token.Ident,
					Val:  "int",
					Pos:  19,
				},
				{
					Kind: token.Ident,
					Val:  "i",
					Pos:  23,
				},
				{
					Kind: token.Semicolon,
					Val:  ";",
					Pos:  24,
				},
				{
					Kind: token.Ident,
					Val:  "i",
					Pos:  28,
				},
				{
					Kind: token.Assign,
					Val:  "=",
					Pos:  30,
				},
				{
					Kind: token.IntLit,
					Val:  "123456789",
					Pos:  32,
				},
				{
					Kind: token.Semicolon,
					Val:  ";",
					Pos:  41,
				},
				{
					Kind: token.Comment,
					Val:  " Was:   i = 1234567890;",
					Pos:  43,
				},
				{
					Kind: token.Ident,
					Val:  "i",
					Pos:  71,
				},
				{
					Kind: token.Assign,
					Val:  "=",
					Pos:  73,
				},
				{
					Kind: token.CharLit,
					Val:  "'0'",
					Pos:  75,
				},
				{
					Kind: token.Semicolon,
					Val:  ";",
					Pos:  78,
				},
				{
					Kind: token.Ident,
					Val:  "i",
					Pos:  82,
				},
				{
					Kind: token.Assign,
					Val:  "=",
					Pos:  84,
				},
				{
					Kind: token.CharLit,
					Val:  "'a'",
					Pos:  86,
				},
				{
					Kind: token.Semicolon,
					Val:  ";",
					Pos:  89,
				},
				{
					Kind: token.Ident,
					Val:  "i",
					Pos:  93,
				},
				{
					Kind: token.Assign,
					Val:  "=",
					Pos:  95,
				},
				{
					Kind: token.CharLit,
					Val:  "' '",
					Pos:  97,
				},
				{
					Kind: token.Semicolon,
					Val:  ";",
					Pos:  100,
				},
				{
					Kind: token.Ident,
					Val:  "i",
					Pos:  104,
				},
				{
					Kind: token.Assign,
					Val:  "=",
					Pos:  106,
				},
				{
					Kind: token.CharLit,
					Val:  `'\n'`,
					Pos:  108,
				},
				{
					Kind: token.Semicolon,
					Val:  ";",
					Pos:  112,
				},
				{
					Kind: token.Rbrace,
					Val:  "}",
					Pos:  114,
				},
				{
					Kind: token.EOF,
					Val:  "",
					Pos:  117,
				},
			},
		},
		// TODO: Add tokens for the following test cases.
		//{path: "../../testdata/quiet/lexer/l04.c"},
		//{path: "../../testdata/quiet/lexer/l05.c"},
		//{path: "../../testdata/quiet/lexer/l06.c"},
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
				break
			}
			got := tokens[j]
			if j >= len(g.toks) {
				t.Errorf("%s: invalid number of tokens; expected %d tokens, got %d", g.path, len(g.toks), j)
				break
			}
			if want := g.toks[j]; got != want {
				t.Errorf("%s: token %d mismatch; expected %#v, got %#v", g.path, j, want, got)
			}
			if got.Kind == token.EOF {
				if j != len(g.toks)-1 {
					t.Errorf("%s: invalid number of tokens; expected %d tokens, got %d", g.path, len(g.toks), j)
				}
				break
			}
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
