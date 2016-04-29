package lexer_test

import (
	"io/ioutil"
	"log"
	"testing"

	"github.com/mewmew/uc/hand/lexer"
	"github.com/mewmew/uc/token"
)

func TestLexer(t *testing.T) {
	encTokens := []token.Token{
		{
			Kind: token.Comment,
			Val:  "/* 世界您好 */",
			Pos:  0,
		},
		{
			Kind: token.Ident,
			Val:  "int",
			Pos:  19,
		},
		{
			Kind: token.Ident,
			Val:  "a",
			Pos:  23,
		},
		{
			Kind: token.Semicolon,
			Val:  ";",
			Pos:  24,
		},
		{
			Kind: token.Comment,
			Val:  "// Hej världen!",
			Pos:  26,
		},
		{
			Kind: token.Ident,
			Val:  "int",
			Pos:  43,
		},
		{
			Kind: token.Ident,
			Val:  "b",
			Pos:  47,
		},
		{
			Kind: token.Semicolon,
			Val:  ";",
			Pos:  48,
		},
		{
			Kind: token.EOF,
			Val:  "",
			Pos:  50,
		},
	}
	var golden = []struct {
		path string
		toks []token.Token
	}{
		{
			path: "../../testdata/incorrect/lexer/bad.c",
			toks: []token.Token{
				{
					Kind: token.Comment,
					Val:  "/*\t\t\t\t\t\t\t||\n**\tFile for testing lexical analysis\t\t||\n**\t\t\t\t\t\t\t||\n**\tThis file is 'lexically incorrect'.\t\t||\n*/",
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
					Pos:  143,
				},
				{
					Kind: token.EOF,
					Val:  "",
					Pos:  221,
				},
			},
		},

		{
			path: "../../testdata/incorrect/lexer/good.c",
			toks: []token.Token{
				{
					Kind: token.Comment,
					Val:  "/*\t\t\t\t\t\t\t||\n**\tFile for testing lexical analysis\t\t||\n**\t\t\t\t\t\t\t||\n**\tThis file would confuse a parser, but\n        is 'lexically correct'.\t\t                ||\n*/",
					Pos:  1,
				},
				{
					Kind: token.Comment,
					Val:  "/* ** / ** */",
					Pos:  163,
				},
				{
					Kind: token.Comment,
					Val:  "// Simple tokens and single characters:",
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
					Val:  "// until end-of-line comment",
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
					Val:  "/* normal comment */",
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
					Val:  "/* Comment with bad tokens: _ || | ++ # @ ...  */",
					Pos:  406,
				},
				{
					Kind: token.Comment,
					Val:  "// Ditto */ /* : _ || | ++ # @ ...  ",
					Pos:  456,
				},
				{
					Kind: token.Comment,
					Val:  "// Identifiers and numbers:",
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
					Val:  "// No floats? -17.17e17 -17.17E-17  ",
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
					Val:  "/* The following 'trap' should be correctly handled:\n\n\t\t* \"2die4U\" consists of the number '2' and the\n\t\t  identifier 'die4U'.\n*/",
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
					Val:  "//|| The following should all be regarded as identifiers:",
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
					Val:  "// It is legal to end the code like this, without an ending newline.",
					Pos:  1031,
				},
				{
					Kind: token.EOF,
					Val:  "",
					Pos:  1099,
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
					Val:  "// OK",
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
					Pos:  50,
				},
				{
					Kind: token.Ident,
					Val:  "cc",
					Pos:  51,
				},
				{
					Kind: token.Error,
					Val:  "unterminated character literal",
					Pos:  53,
				},
				{
					Kind: token.Semicolon,
					Val:  ";",
					Pos:  54,
				},
				{
					Kind: token.Comment,
					Val:  "// Not OK",
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

		{
			path: "../../testdata/incorrect/lexer/ugly.c",
			toks: []token.Token{
				{
					Kind: token.Error,
					Val:  "unexpected U+007C '|'",
					Pos:  0,
				},
				{
					Kind: token.Add,
					Val:  "+",
					Pos:  1,
				},
				{
					Kind: token.Error,
					Val:  "unexpected U+0001",
					Pos:  2,
				},
				{
					Kind: token.Error,
					Val:  "unexpected U+0002",
					Pos:  3,
				},
				{
					Kind: token.Error,
					Val:  "unexpected U+0003",
					Pos:  4,
				},
				{
					Kind: token.Error,
					Val:  "unexpected U+0004",
					Pos:  5,
				},
				{
					Kind: token.Error,
					Val:  "unexpected U+0005",
					Pos:  6,
				},
				{
					Kind: token.Error,
					Val:  "unexpected U+0006",
					Pos:  7,
				},
				{
					Kind: token.Error,
					Val:  "unexpected U+0007",
					Pos:  8,
				},
				{
					Kind: token.Error,
					Val:  "unexpected U+0008",
					Pos:  9,
				},
				{
					Kind: token.Error,
					Val:  "unexpected U+000E",
					Pos:  15,
				},
				{
					Kind: token.Error,
					Val:  "unexpected U+000F",
					Pos:  16,
				},
				{
					Kind: token.Error,
					Val:  "unexpected U+0010",
					Pos:  17,
				},
				{
					Kind: token.Error,
					Val:  "unexpected U+0011",
					Pos:  18,
				},
				{
					Kind: token.Error,
					Val:  "unexpected U+0012",
					Pos:  19,
				},
				{
					Kind: token.Error,
					Val:  "unexpected U+0013",
					Pos:  20,
				},
				{
					Kind: token.Error,
					Val:  "unexpected U+0014",
					Pos:  21,
				},
				{
					Kind: token.Error,
					Val:  "unexpected U+0015",
					Pos:  22,
				},
				{
					Kind: token.Error,
					Val:  "unexpected U+0016",
					Pos:  23,
				},
				{
					Kind: token.Error,
					Val:  "unexpected U+0017",
					Pos:  24,
				},
				{
					Kind: token.Error,
					Val:  "unexpected U+0018",
					Pos:  25,
				},
				{
					Kind: token.Error,
					Val:  "unexpected U+0019",
					Pos:  26,
				},
				{
					Kind: token.Error,
					Val:  "unexpected U+001A",
					Pos:  27,
				},
				{
					Kind: token.Error,
					Val:  "unexpected U+001B",
					Pos:  28,
				},
				{
					Kind: token.Error,
					Val:  "unexpected U+001C",
					Pos:  29,
				},
				{
					Kind: token.Error,
					Val:  "unexpected U+001D",
					Pos:  30,
				},
				{
					Kind: token.Error,
					Val:  "unexpected U+001E",
					Pos:  31,
				},
				{
					Kind: token.Error,
					Val:  "unexpected U+001F",
					Pos:  32,
				},
				{
					Kind: token.Not,
					Val:  "!",
					Pos:  34,
				},
				{
					Kind: token.Error,
					Val:  `unexpected U+0022 '"'`,
					Pos:  35,
				},
				{
					Kind: token.Error,
					Val:  "unexpected U+0023 '#'",
					Pos:  36,
				},
				{
					Kind: token.Error,
					Val:  "unexpected U+0024 '$'",
					Pos:  37,
				},
				{
					Kind: token.Error,
					Val:  "unexpected U+0025 '%'",
					Pos:  38,
				},
				{
					Kind: token.Error,
					Val:  "expected '&' after '&', got U+0027 '''",
					Pos:  39,
				},
				{
					Kind: token.Error,
					Val:  "unterminated character literal",
					Pos:  40,
				},
				{
					Kind: token.Lparen,
					Val:  "(",
					Pos:  41,
				},
				{
					Kind: token.Rparen,
					Val:  ")",
					Pos:  42,
				},
				{
					Kind: token.Mul,
					Val:  "*",
					Pos:  43,
				},
				{
					Kind: token.Add,
					Val:  "+",
					Pos:  44,
				},
				{
					Kind: token.Comma,
					Val:  ",",
					Pos:  46,
				},
				{
					Kind: token.Sub,
					Val:  "-",
					Pos:  47,
				},
				{
					Kind: token.Error,
					Val:  "unexpected U+002E '.'",
					Pos:  48,
				},
				{
					Kind: token.Div,
					Val:  "/",
					Pos:  49,
				},
				{
					Kind: token.IntLit,
					Val:  "0123456789", // TODO: When octal integer literals have been implemented, fail accordingly.
					Pos:  50,
				},
				{
					Kind: token.Error,
					Val:  "unexpected U+003A ':'",
					Pos:  60,
				},
				{
					Kind: token.Semicolon,
					Val:  ";",
					Pos:  61,
				},
				{
					Kind: token.Le,
					Val:  "<=",
					Pos:  62,
				},
				{
					Kind: token.Gt,
					Val:  ">",
					Pos:  64,
				},
				{
					Kind: token.Error,
					Val:  "unexpected U+003F '?'",
					Pos:  65,
				},
				{
					Kind: token.Error,
					Val:  "unexpected U+0040 '@'",
					Pos:  66,
				},
				{
					Kind: token.Ident,
					Val:  "ABCDEFGHIJKLMNOPQRSTUVWXYZ",
					Pos:  67,
				},
				{
					Kind: token.Lbracket,
					Val:  "[",
					Pos:  93,
				},
				{
					Kind: token.Error,
					Val:  `unexpected U+005C '\'`,
					Pos:  94,
				},
				{
					Kind: token.Rbracket,
					Val:  "]",
					Pos:  95,
				},
				{
					Kind: token.Error,
					Val:  "unexpected U+005E '^'",
					Pos:  96,
				},
				{
					Kind: token.Ident,
					Val:  "_",
					Pos:  97,
				},
				{
					Kind: token.Error,
					Val:  "unexpected U+0060 '`'",
					Pos:  98,
				},
				{
					Kind: token.Ident,
					Val:  "abcdefghijklmnopqrstuvwxyz",
					Pos:  99,
				},
				{
					Kind: token.Lbrace,
					Val:  "{",
					Pos:  125,
				},
				{
					Kind: token.Error,
					Val:  "unexpected U+007C '|'",
					Pos:  126,
				},
				{
					Kind: token.Rbrace,
					Val:  "}",
					Pos:  127,
				},
				{
					Kind: token.Error,
					Val:  "unexpected U+007E '~'",
					Pos:  128,
				},
				{
					Kind: token.Error,
					Val:  "unexpected U+007F",
					Pos:  129,
				},
				{
					Kind: token.Error,
					Val:  "illegal UTF-8 encoding",
					Pos:  130,
				},
				{
					Kind: token.Error,
					Val:  "illegal UTF-8 encoding",
					Pos:  131,
				},
				{
					Kind: token.Error,
					Val:  "illegal UTF-8 encoding",
					Pos:  132,
				},
				{
					Kind: token.Error,
					Val:  "illegal UTF-8 encoding",
					Pos:  133,
				},
				{
					Kind: token.Error,
					Val:  "illegal UTF-8 encoding",
					Pos:  134,
				},
				{
					Kind: token.Error,
					Val:  "illegal UTF-8 encoding",
					Pos:  135,
				},
				{
					Kind: token.Error,
					Val:  "illegal UTF-8 encoding",
					Pos:  136,
				},
				{
					Kind: token.Error,
					Val:  "illegal UTF-8 encoding",
					Pos:  137,
				},
				{
					Kind: token.Error,
					Val:  "illegal UTF-8 encoding",
					Pos:  138,
				},
				{
					Kind: token.Error,
					Val:  "illegal UTF-8 encoding",
					Pos:  139,
				},
				{
					Kind: token.Error,
					Val:  "illegal UTF-8 encoding",
					Pos:  140,
				},
				{
					Kind: token.Error,
					Val:  "illegal UTF-8 encoding",
					Pos:  141,
				},
				{
					Kind: token.Error,
					Val:  "illegal UTF-8 encoding",
					Pos:  142,
				},
				{
					Kind: token.Error,
					Val:  "illegal UTF-8 encoding",
					Pos:  143,
				},
				{
					Kind: token.Error,
					Val:  "illegal UTF-8 encoding",
					Pos:  144,
				},
				{
					Kind: token.Error,
					Val:  "illegal UTF-8 encoding",
					Pos:  145,
				},
				{
					Kind: token.Error,
					Val:  "illegal UTF-8 encoding",
					Pos:  146,
				},
				{
					Kind: token.Error,
					Val:  "illegal UTF-8 encoding",
					Pos:  147,
				},
				{
					Kind: token.Error,
					Val:  "illegal UTF-8 encoding",
					Pos:  148,
				},
				{
					Kind: token.Error,
					Val:  "illegal UTF-8 encoding",
					Pos:  149,
				},
				{
					Kind: token.Error,
					Val:  "illegal UTF-8 encoding",
					Pos:  150,
				},
				{
					Kind: token.Error,
					Val:  "illegal UTF-8 encoding",
					Pos:  151,
				},
				{
					Kind: token.Error,
					Val:  "illegal UTF-8 encoding",
					Pos:  152,
				},
				{
					Kind: token.Error,
					Val:  "illegal UTF-8 encoding",
					Pos:  153,
				},
				{
					Kind: token.Error,
					Val:  "illegal UTF-8 encoding",
					Pos:  154,
				},
				{
					Kind: token.Error,
					Val:  "illegal UTF-8 encoding",
					Pos:  155,
				},
				{
					Kind: token.Error,
					Val:  "illegal UTF-8 encoding",
					Pos:  156,
				},
				{
					Kind: token.Error,
					Val:  "illegal UTF-8 encoding",
					Pos:  157,
				},
				{
					Kind: token.Error,
					Val:  "illegal UTF-8 encoding",
					Pos:  158,
				},
				{
					Kind: token.Error,
					Val:  "illegal UTF-8 encoding",
					Pos:  159,
				},
				{
					Kind: token.Error,
					Val:  "illegal UTF-8 encoding",
					Pos:  160,
				},
				{
					Kind: token.Error,
					Val:  "illegal UTF-8 encoding",
					Pos:  161,
				},
				{
					Kind: token.Error,
					Val:  "illegal UTF-8 encoding",
					Pos:  162,
				},
				{
					Kind: token.Error,
					Val:  "illegal UTF-8 encoding",
					Pos:  163,
				},
				{
					Kind: token.Error,
					Val:  "illegal UTF-8 encoding",
					Pos:  164,
				},
				{
					Kind: token.Error,
					Val:  "illegal UTF-8 encoding",
					Pos:  165,
				},
				{
					Kind: token.Error,
					Val:  "illegal UTF-8 encoding",
					Pos:  166,
				},
				{
					Kind: token.Error,
					Val:  "illegal UTF-8 encoding",
					Pos:  167,
				},
				{
					Kind: token.Error,
					Val:  "illegal UTF-8 encoding",
					Pos:  168,
				},
				{
					Kind: token.Error,
					Val:  "illegal UTF-8 encoding",
					Pos:  169,
				},
				{
					Kind: token.Error,
					Val:  "illegal UTF-8 encoding",
					Pos:  170,
				},
				{
					Kind: token.Error,
					Val:  "illegal UTF-8 encoding",
					Pos:  171,
				},
				{
					Kind: token.Error,
					Val:  "illegal UTF-8 encoding",
					Pos:  172,
				},
				{
					Kind: token.Error,
					Val:  "illegal UTF-8 encoding",
					Pos:  173,
				},
				{
					Kind: token.Error,
					Val:  "illegal UTF-8 encoding",
					Pos:  174,
				},
				{
					Kind: token.Error,
					Val:  "illegal UTF-8 encoding",
					Pos:  175,
				},
				{
					Kind: token.Error,
					Val:  "illegal UTF-8 encoding",
					Pos:  176,
				},
				{
					Kind: token.Error,
					Val:  "illegal UTF-8 encoding",
					Pos:  177,
				},
				{
					Kind: token.Error,
					Val:  "illegal UTF-8 encoding",
					Pos:  178,
				},
				{
					Kind: token.Error,
					Val:  "illegal UTF-8 encoding",
					Pos:  179,
				},
				{
					Kind: token.Error,
					Val:  "illegal UTF-8 encoding",
					Pos:  180,
				},
				{
					Kind: token.Error,
					Val:  "illegal UTF-8 encoding",
					Pos:  181,
				},
				{
					Kind: token.Error,
					Val:  "illegal UTF-8 encoding",
					Pos:  182,
				},
				{
					Kind: token.Error,
					Val:  "illegal UTF-8 encoding",
					Pos:  183,
				},
				{
					Kind: token.Error,
					Val:  "illegal UTF-8 encoding",
					Pos:  184,
				},
				{
					Kind: token.Error,
					Val:  "illegal UTF-8 encoding",
					Pos:  185,
				},
				{
					Kind: token.Error,
					Val:  "illegal UTF-8 encoding",
					Pos:  186,
				},
				{
					Kind: token.Error,
					Val:  "illegal UTF-8 encoding",
					Pos:  187,
				},
				{
					Kind: token.Error,
					Val:  "illegal UTF-8 encoding",
					Pos:  188,
				},
				{
					Kind: token.Error,
					Val:  "illegal UTF-8 encoding",
					Pos:  189,
				},
				{
					Kind: token.Error,
					Val:  "illegal UTF-8 encoding",
					Pos:  190,
				},
				{
					Kind: token.Error,
					Val:  "illegal UTF-8 encoding",
					Pos:  191,
				},
				{
					Kind: token.Error,
					Val:  "illegal UTF-8 encoding",
					Pos:  192,
				},
				{
					Kind: token.Error,
					Val:  "illegal UTF-8 encoding",
					Pos:  193,
				},
				{
					Kind: token.Error,
					Val:  "illegal UTF-8 encoding",
					Pos:  194,
				},
				{
					Kind: token.Error,
					Val:  "illegal UTF-8 encoding",
					Pos:  195,
				},
				{
					Kind: token.Error,
					Val:  "illegal UTF-8 encoding",
					Pos:  196,
				},
				{
					Kind: token.Error,
					Val:  "illegal UTF-8 encoding",
					Pos:  197,
				},
				{
					Kind: token.Error,
					Val:  "illegal UTF-8 encoding",
					Pos:  198,
				},
				{
					Kind: token.Error,
					Val:  "illegal UTF-8 encoding",
					Pos:  199,
				},
				{
					Kind: token.Error,
					Val:  "illegal UTF-8 encoding",
					Pos:  200,
				},
				{
					Kind: token.Error,
					Val:  "illegal UTF-8 encoding",
					Pos:  201,
				},
				{
					Kind: token.Error,
					Val:  "illegal UTF-8 encoding",
					Pos:  202,
				},
				{
					Kind: token.Error,
					Val:  "illegal UTF-8 encoding",
					Pos:  203,
				},
				{
					Kind: token.Error,
					Val:  "illegal UTF-8 encoding",
					Pos:  204,
				},
				{
					Kind: token.Error,
					Val:  "illegal UTF-8 encoding",
					Pos:  205,
				},
				{
					Kind: token.Error,
					Val:  "illegal UTF-8 encoding",
					Pos:  206,
				},
				{
					Kind: token.Error,
					Val:  "illegal UTF-8 encoding",
					Pos:  207,
				},
				{
					Kind: token.Error,
					Val:  "illegal UTF-8 encoding",
					Pos:  208,
				},
				{
					Kind: token.Error,
					Val:  "illegal UTF-8 encoding",
					Pos:  209,
				},
				{
					Kind: token.Error,
					Val:  "illegal UTF-8 encoding",
					Pos:  210,
				},
				{
					Kind: token.Error,
					Val:  "illegal UTF-8 encoding",
					Pos:  211,
				},
				{
					Kind: token.Error,
					Val:  "illegal UTF-8 encoding",
					Pos:  212,
				},
				{
					Kind: token.Error,
					Val:  "illegal UTF-8 encoding",
					Pos:  213,
				},
				{
					Kind: token.Error,
					Val:  "illegal UTF-8 encoding",
					Pos:  214,
				},
				{
					Kind: token.Error,
					Val:  "illegal UTF-8 encoding",
					Pos:  215,
				},
				{
					Kind: token.Error,
					Val:  "illegal UTF-8 encoding",
					Pos:  216,
				},
				{
					Kind: token.Error,
					Val:  "illegal UTF-8 encoding",
					Pos:  217,
				},
				{
					Kind: token.Error,
					Val:  "illegal UTF-8 encoding",
					Pos:  218,
				},
				{
					Kind: token.Error,
					Val:  "illegal UTF-8 encoding",
					Pos:  219,
				},
				{
					Kind: token.Error,
					Val:  "illegal UTF-8 encoding",
					Pos:  220,
				},
				{
					Kind: token.Error,
					Val:  "illegal UTF-8 encoding",
					Pos:  221,
				},
				{
					Kind: token.Error,
					Val:  "illegal UTF-8 encoding",
					Pos:  222,
				},
				{
					Kind: token.Error,
					Val:  "illegal UTF-8 encoding",
					Pos:  223,
				},
				{
					Kind: token.Error,
					Val:  "illegal UTF-8 encoding",
					Pos:  224,
				},
				{
					Kind: token.Error,
					Val:  "illegal UTF-8 encoding",
					Pos:  225,
				},
				{
					Kind: token.Error,
					Val:  "illegal UTF-8 encoding",
					Pos:  226,
				},
				{
					Kind: token.Error,
					Val:  "illegal UTF-8 encoding",
					Pos:  227,
				},
				{
					Kind: token.Error,
					Val:  "illegal UTF-8 encoding",
					Pos:  228,
				},
				{
					Kind: token.Error,
					Val:  "illegal UTF-8 encoding",
					Pos:  229,
				},
				{
					Kind: token.Error,
					Val:  "illegal UTF-8 encoding",
					Pos:  230,
				},
				{
					Kind: token.Error,
					Val:  "illegal UTF-8 encoding",
					Pos:  231,
				},
				{
					Kind: token.Error,
					Val:  "illegal UTF-8 encoding",
					Pos:  232,
				},
				{
					Kind: token.Error,
					Val:  "illegal UTF-8 encoding",
					Pos:  233,
				},
				{
					Kind: token.Error,
					Val:  "illegal UTF-8 encoding",
					Pos:  234,
				},
				{
					Kind: token.Error,
					Val:  "illegal UTF-8 encoding",
					Pos:  235,
				},
				{
					Kind: token.Error,
					Val:  "illegal UTF-8 encoding",
					Pos:  236,
				},
				{
					Kind: token.Error,
					Val:  "illegal UTF-8 encoding",
					Pos:  237,
				},
				{
					Kind: token.Error,
					Val:  "illegal UTF-8 encoding",
					Pos:  238,
				},
				{
					Kind: token.Error,
					Val:  "illegal UTF-8 encoding",
					Pos:  239,
				},
				{
					Kind: token.Error,
					Val:  "illegal UTF-8 encoding",
					Pos:  240,
				},
				{
					Kind: token.Error,
					Val:  "illegal UTF-8 encoding",
					Pos:  241,
				},
				{
					Kind: token.Error,
					Val:  "illegal UTF-8 encoding",
					Pos:  242,
				},
				{
					Kind: token.Error,
					Val:  "illegal UTF-8 encoding",
					Pos:  243,
				},
				{
					Kind: token.Error,
					Val:  "illegal UTF-8 encoding",
					Pos:  244,
				},
				{
					Kind: token.Error,
					Val:  "illegal UTF-8 encoding",
					Pos:  245,
				},
				{
					Kind: token.Error,
					Val:  "illegal UTF-8 encoding",
					Pos:  246,
				},
				{
					Kind: token.Error,
					Val:  "illegal UTF-8 encoding",
					Pos:  247,
				},
				{
					Kind: token.Error,
					Val:  "illegal UTF-8 encoding",
					Pos:  248,
				},
				{
					Kind: token.Error,
					Val:  "illegal UTF-8 encoding",
					Pos:  249,
				},
				{
					Kind: token.Error,
					Val:  "illegal UTF-8 encoding",
					Pos:  250,
				},
				{
					Kind: token.Error,
					Val:  "illegal UTF-8 encoding",
					Pos:  251,
				},
				{
					Kind: token.Error,
					Val:  "illegal UTF-8 encoding",
					Pos:  252,
				},
				{
					Kind: token.Error,
					Val:  "illegal UTF-8 encoding",
					Pos:  253,
				},
				{
					Kind: token.Error,
					Val:  "illegal UTF-8 encoding",
					Pos:  254,
				},
				{
					Kind: token.Error,
					Val:  "illegal UTF-8 encoding",
					Pos:  255,
				},
				{
					Kind: token.Error,
					Val:  "illegal UTF-8 encoding",
					Pos:  256,
				},
				{
					Kind: token.Error,
					Val:  "illegal UTF-8 encoding",
					Pos:  257,
				},
				{
					Kind: token.Add,
					Val:  "+",
					Pos:  258,
				},
				{
					Kind: token.Error,
					Val:  "unexpected U+007C '|'",
					Pos:  259,
				},
				{
					Kind: token.Ident,
					Val:  "int",
					Pos:  261,
				},
				{
					Kind: token.Ident,
					Val:  "a",
					Pos:  265,
				},
				{
					Kind: token.Comma,
					Val:  ",",
					Pos:  267,
				},
				{
					Kind: token.Ident,
					Val:  "b",
					Pos:  269,
				},
				{
					Kind: token.Comma,
					Val:  ",",
					Pos:  271,
				},
				{
					Kind: token.Ident,
					Val:  "c",
					Pos:  273,
				},
				{
					Kind: token.Comma,
					Val:  ",",
					Pos:  275,
				},
				{
					Kind: token.Ident,
					Val:  "d",
					Pos:  277,
				},
				{
					Kind: token.Semicolon,
					Val:  ";",
					Pos:  279,
				},
				{
					Kind: token.Ident,
					Val:  "BEGIN",
					Pos:  281,
				},
				{
					Kind: token.Ident,
					Val:  "a",
					Pos:  287,
				},
				{
					Kind: token.Assign,
					Val:  "=",
					Pos:  289,
				},
				{
					Kind: token.Lparen,
					Val:  "(",
					Pos:  291,
				},
				{
					Kind: token.IntLit,
					Val:  "4711",
					Pos:  293,
				},
				{
					Kind: token.Add,
					Val:  "+",
					Pos:  298,
				},
				{
					Kind: token.IntLit,
					Val:  "17",
					Pos:  300,
				},
				{
					Kind: token.Sub,
					Val:  "-",
					Pos:  303,
				},
				{
					Kind: token.IntLit,
					Val:  "2001",
					Pos:  305,
				},
				{
					Kind: token.Rparen,
					Val:  ")",
					Pos:  310,
				},
				{
					Kind: token.Div,
					Val:  "/",
					Pos:  312,
				},
				{
					Kind: token.IntLit,
					Val:  "100",
					Pos:  314,
				},
				{
					Kind: token.Sub,
					Val:  "-",
					Pos:  318,
				},
				{
					Kind: token.IntLit,
					Val:  "17",
					Pos:  320,
				},
				{
					Kind: token.Semicolon,
					Val:  ";",
					Pos:  323,
				},
				{
					Kind: token.Ident,
					Val:  "b",
					Pos:  325,
				},
				{
					Kind: token.Assign,
					Val:  "=",
					Pos:  327,
				},
				{
					Kind: token.Ident,
					Val:  "a",
					Pos:  329,
				},
				{
					Kind: token.Div,
					Val:  "/",
					Pos:  331,
				},
				{
					Kind: token.IntLit,
					Val:  "5",
					Pos:  333,
				},
				{
					Kind: token.Semicolon,
					Val:  ";",
					Pos:  335,
				},
				{
					Kind: token.Ident,
					Val:  "c",
					Pos:  337,
				},
				{
					Kind: token.Assign,
					Val:  "=",
					Pos:  339,
				},
				{
					Kind: token.Ident,
					Val:  "a",
					Pos:  341,
				},
				{
					Kind: token.Mul,
					Val:  "*",
					Pos:  343,
				},
				{
					Kind: token.Ident,
					Val:  "b",
					Pos:  345,
				},
				{
					Kind: token.Mul,
					Val:  "*",
					Pos:  347,
				},
				{
					Kind: token.Ident,
					Val:  "b",
					Pos:  349,
				},
				{
					Kind: token.Mul,
					Val:  "*",
					Pos:  351,
				},
				{
					Kind: token.Ident,
					Val:  "a",
					Pos:  353,
				},
				{
					Kind: token.Mul,
					Val:  "*",
					Pos:  355,
				},
				{
					Kind: token.Ident,
					Val:  "b",
					Pos:  357,
				},
				{
					Kind: token.Mul,
					Val:  "*",
					Pos:  359,
				},
				{
					Kind: token.Ident,
					Val:  "b",
					Pos:  361,
				},
				{
					Kind: token.Semicolon,
					Val:  ";",
					Pos:  363,
				},
				{
					Kind: token.Ident,
					Val:  "d",
					Pos:  365,
				},
				{
					Kind: token.Assign,
					Val:  "=",
					Pos:  367,
				},
				{
					Kind: token.Ident,
					Val:  "c",
					Pos:  369,
				},
				{
					Kind: token.Add,
					Val:  "+",
					Pos:  371,
				},
				{
					Kind: token.Ident,
					Val:  "c",
					Pos:  373,
				},
				{
					Kind: token.Div,
					Val:  "/",
					Pos:  375,
				},
				{
					Kind: token.Lparen,
					Val:  "(",
					Pos:  377,
				},
				{
					Kind: token.IntLit,
					Val:  "100",
					Pos:  379,
				},
				{
					Kind: token.Div,
					Val:  "/",
					Pos:  383,
				},
				{
					Kind: token.IntLit,
					Val:  "2",
					Pos:  385,
				},
				{
					Kind: token.Rparen,
					Val:  ")",
					Pos:  387,
				},
				{
					Kind: token.Semicolon,
					Val:  ";",
					Pos:  389,
				},
				{
					Kind: token.KwIf,
					Val:  "if",
					Pos:  391,
				},
				{
					Kind: token.Lparen,
					Val:  "(",
					Pos:  394,
				},
				{
					Kind: token.Ident,
					Val:  "a",
					Pos:  396,
				},
				{
					Kind: token.Gt,
					Val:  ">",
					Pos:  398,
				},
				{
					Kind: token.IntLit,
					Val:  "3",
					Pos:  400,
				},
				{
					Kind: token.Mul,
					Val:  "*",
					Pos:  402,
				},
				{
					Kind: token.Ident,
					Val:  "b",
					Pos:  404,
				},
				{
					Kind: token.Ident,
					Val:  "or",
					Pos:  406,
				},
				{
					Kind: token.Ident,
					Val:  "not",
					Pos:  409,
				},
				{
					Kind: token.Ident,
					Val:  "a",
					Pos:  413,
				},
				{
					Kind: token.Eq,
					Val:  "==",
					Pos:  415,
				},
				{
					Kind: token.Ident,
					Val:  "c",
					Pos:  418,
				},
				{
					Kind: token.Div,
					Val:  "/",
					Pos:  420,
				},
				{
					Kind: token.IntLit,
					Val:  "2",
					Pos:  422,
				},
				{
					Kind: token.Ident,
					Val:  "and",
					Pos:  424,
				},
				{
					Kind: token.Ident,
					Val:  "c",
					Pos:  428,
				},
				{
					Kind: token.Gt,
					Val:  ">",
					Pos:  430,
				},
				{
					Kind: token.Lparen,
					Val:  "(",
					Pos:  432,
				},
				{
					Kind: token.Ident,
					Val:  "d",
					Pos:  434,
				},
				{
					Kind: token.Sub,
					Val:  "-",
					Pos:  436,
				},
				{
					Kind: token.IntLit,
					Val:  "3",
					Pos:  438,
				},
				{
					Kind: token.Rparen,
					Val:  ")",
					Pos:  440,
				},
				{
					Kind: token.Mul,
					Val:  "*",
					Pos:  442,
				},
				{
					Kind: token.Lparen,
					Val:  "(",
					Pos:  444,
				},
				{
					Kind: token.Ident,
					Val:  "a",
					Pos:  446,
				},
				{
					Kind: token.Sub,
					Val:  "-",
					Pos:  448,
				},
				{
					Kind: token.Ident,
					Val:  "b",
					Pos:  450,
				},
				{
					Kind: token.Rparen,
					Val:  ")",
					Pos:  452,
				},
				{
					Kind: token.Rparen,
					Val:  ")",
					Pos:  454,
				},
				{
					Kind: token.Ident,
					Val:  "and",
					Pos:  456,
				},
				{
					Kind: token.Ident,
					Val:  "a",
					Pos:  460,
				},
				{
					Kind: token.Lt,
					Val:  "<",
					Pos:  462,
				},
				{
					Kind: token.IntLit,
					Val:  "0",
					Pos:  464,
				},
				{
					Kind: token.Ident,
					Val:  "then",
					Pos:  466,
				},
				{
					Kind: token.Ident,
					Val:  "begin",
					Pos:  471,
				},
				{
					Kind: token.Ident,
					Val:  "a",
					Pos:  477,
				},
				{
					Kind: token.Assign,
					Val:  "=",
					Pos:  479,
				},
				{
					Kind: token.IntLit,
					Val:  "1",
					Pos:  481,
				},
				{
					Kind: token.Semicolon,
					Val:  ";",
					Pos:  483,
				},
				{
					Kind: token.Ident,
					Val:  "read",
					Pos:  485,
				},
				{
					Kind: token.Ident,
					Val:  "b",
					Pos:  490,
				},
				{
					Kind: token.Semicolon,
					Val:  ";",
					Pos:  492,
				},
				{
					Kind: token.Ident,
					Val:  "end",
					Pos:  494,
				},
				{
					Kind: token.KwElse,
					Val:  "else",
					Pos:  498,
				},
				{
					Kind: token.Ident,
					Val:  "begin",
					Pos:  503,
				},
				{
					Kind: token.Ident,
					Val:  "a",
					Pos:  509,
				},
				{
					Kind: token.Assign,
					Val:  "=",
					Pos:  511,
				},
				{
					Kind: token.IntLit,
					Val:  "0",
					Pos:  513,
				},
				{
					Kind: token.Semicolon,
					Val:  ";",
					Pos:  515,
				},
				{
					Kind: token.Ident,
					Val:  "read",
					Pos:  517,
				},
				{
					Kind: token.Ident,
					Val:  "c",
					Pos:  522,
				},
				{
					Kind: token.Semicolon,
					Val:  ";",
					Pos:  524,
				},
				{
					Kind: token.Ident,
					Val:  "end",
					Pos:  526,
				},
				{
					Kind: token.Ident,
					Val:  "END",
					Pos:  530,
				},
				{
					Kind: token.EOF,
					Val:  "",
					Pos:  533,
				},
			},
		},

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
					Val:  "// Was:   i = 1234567890;",
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

		{
			path: "../../testdata/quiet/lexer/l04.c",
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
					Val:  "char",
					Pos:  28,
				},
				{
					Kind: token.Ident,
					Val:  "j",
					Pos:  33,
				},
				{
					Kind: token.Semicolon,
					Val:  ";",
					Pos:  34,
				},
				{
					Kind: token.KwIf,
					Val:  "if",
					Pos:  38,
				},
				{
					Kind: token.Lparen,
					Val:  "(",
					Pos:  41,
				},
				{
					Kind: token.IntLit,
					Val:  "1",
					Pos:  42,
				},
				{
					Kind: token.Eq,
					Val:  "==",
					Pos:  43,
				},
				{
					Kind: token.IntLit,
					Val:  "0",
					Pos:  45,
				},
				{
					Kind: token.Rparen,
					Val:  ")",
					Pos:  46,
				},
				{
					Kind: token.Ident,
					Val:  "i",
					Pos:  48,
				},
				{
					Kind: token.Assign,
					Val:  "=",
					Pos:  50,
				},
				{
					Kind: token.IntLit,
					Val:  "0",
					Pos:  52,
				},
				{
					Kind: token.Semicolon,
					Val:  ";",
					Pos:  53,
				},
				{
					Kind: token.KwElse,
					Val:  "else",
					Pos:  58,
				},
				{
					Kind: token.Ident,
					Val:  "i",
					Pos:  63,
				},
				{
					Kind: token.Assign,
					Val:  "=",
					Pos:  65,
				},
				{
					Kind: token.IntLit,
					Val:  "1",
					Pos:  67,
				},
				{
					Kind: token.Semicolon,
					Val:  ";",
					Pos:  68,
				},
				{
					Kind: token.KwWhile,
					Val:  "while",
					Pos:  72,
				},
				{
					Kind: token.Lparen,
					Val:  "(",
					Pos:  78,
				},
				{
					Kind: token.IntLit,
					Val:  "1",
					Pos:  79,
				},
				{
					Kind: token.Eq,
					Val:  "==",
					Pos:  80,
				},
				{
					Kind: token.IntLit,
					Val:  "0",
					Pos:  82,
				},
				{
					Kind: token.Rparen,
					Val:  ")",
					Pos:  83,
				},
				{
					Kind: token.Ident,
					Val:  "i",
					Pos:  85,
				},
				{
					Kind: token.Assign,
					Val:  "=",
					Pos:  87,
				},
				{
					Kind: token.IntLit,
					Val:  "0",
					Pos:  89,
				},
				{
					Kind: token.Semicolon,
					Val:  ";",
					Pos:  90,
				},
				{
					Kind: token.KwReturn,
					Val:  "return",
					Pos:  94,
				},
				{
					Kind: token.IntLit,
					Val:  "42",
					Pos:  101,
				},
				{
					Kind: token.Semicolon,
					Val:  ";",
					Pos:  103,
				},
				{
					Kind: token.Rbrace,
					Val:  "}",
					Pos:  105,
				},
				{
					Kind: token.EOF,
					Val:  "",
					Pos:  107,
				},
			},
		},

		{
			path: "../../testdata/quiet/lexer/l05.c",
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
					Kind: token.IntLit,
					Val:  "1",
					Pos:  28,
				},
				{
					Kind: token.Ne,
					Val:  "!=",
					Pos:  29,
				},
				{
					Kind: token.Not,
					Val:  "!",
					Pos:  31,
				},
				{
					Kind: token.IntLit,
					Val:  "3",
					Pos:  32,
				},
				{
					Kind: token.Semicolon,
					Val:  ";",
					Pos:  33,
				},
				{
					Kind: token.IntLit,
					Val:  "4",
					Pos:  37,
				},
				{
					Kind: token.Land,
					Val:  "&&",
					Pos:  38,
				},
				{
					Kind: token.Lparen,
					Val:  "(",
					Pos:  40,
				},
				{
					Kind: token.IntLit,
					Val:  "6",
					Pos:  41,
				},
				{
					Kind: token.Rparen,
					Val:  ")",
					Pos:  42,
				},
				{
					Kind: token.Semicolon,
					Val:  ";",
					Pos:  43,
				},
				{
					Kind: token.IntLit,
					Val:  "7",
					Pos:  47,
				},
				{
					Kind: token.Mul,
					Val:  "*",
					Pos:  48,
				},
				{
					Kind: token.IntLit,
					Val:  "8",
					Pos:  50,
				},
				{
					Kind: token.Add,
					Val:  "+",
					Pos:  51,
				},
				{
					Kind: token.IntLit,
					Val:  "10",
					Pos:  52,
				},
				{
					Kind: token.Semicolon,
					Val:  ";",
					Pos:  54,
				},
				{
					Kind: token.Lparen,
					Val:  "(",
					Pos:  58,
				},
				{
					Kind: token.IntLit,
					Val:  "11",
					Pos:  59,
				},
				{
					Kind: token.Sub,
					Val:  "-",
					Pos:  61,
				},
				{
					Kind: token.IntLit,
					Val:  "12",
					Pos:  62,
				},
				{
					Kind: token.Rparen,
					Val:  ")",
					Pos:  64,
				},
				{
					Kind: token.Add,
					Val:  "+",
					Pos:  65,
				},
				{
					Kind: token.Lparen,
					Val:  "(",
					Pos:  66,
				},
				{
					Kind: token.IntLit,
					Val:  "12",
					Pos:  67,
				},
				{
					Kind: token.Div,
					Val:  "/",
					Pos:  69,
				},
				{
					Kind: token.IntLit,
					Val:  "16",
					Pos:  70,
				},
				{
					Kind: token.Rparen,
					Val:  ")",
					Pos:  72,
				},
				{
					Kind: token.Semicolon,
					Val:  ";",
					Pos:  73,
				},
				{
					Kind: token.IntLit,
					Val:  "17",
					Pos:  77,
				},
				{
					Kind: token.Le,
					Val:  "<=",
					Pos:  79,
				},
				{
					Kind: token.IntLit,
					Val:  "18",
					Pos:  81,
				},
				{
					Kind: token.Lt,
					Val:  "<",
					Pos:  84,
				},
				{
					Kind: token.Sub,
					Val:  "-",
					Pos:  85,
				},
				{
					Kind: token.IntLit,
					Val:  "20",
					Pos:  86,
				},
				{
					Kind: token.Semicolon,
					Val:  ";",
					Pos:  88,
				},
				{
					Kind: token.Ident,
					Val:  "i",
					Pos:  92,
				},
				{
					Kind: token.Assign,
					Val:  "=",
					Pos:  93,
				},
				{
					Kind: token.IntLit,
					Val:  "21",
					Pos:  94,
				},
				{
					Kind: token.Eq,
					Val:  "==",
					Pos:  96,
				},
				{
					Kind: token.IntLit,
					Val:  "22",
					Pos:  98,
				},
				{
					Kind: token.Semicolon,
					Val:  ";",
					Pos:  100,
				},
				{
					Kind: token.IntLit,
					Val:  "25",
					Pos:  104,
				},
				{
					Kind: token.Ge,
					Val:  ">=",
					Pos:  107,
				},
				{
					Kind: token.IntLit,
					Val:  "27",
					Pos:  109,
				},
				{
					Kind: token.Gt,
					Val:  ">",
					Pos:  111,
				},
				{
					Kind: token.IntLit,
					Val:  "28",
					Pos:  112,
				},
				{
					Kind: token.Semicolon,
					Val:  ";",
					Pos:  114,
				},
				{
					Kind: token.Rbrace,
					Val:  "}",
					Pos:  116,
				},
				{
					Kind: token.EOF,
					Val:  "",
					Pos:  118,
				},
			},
		},

		{
			path: "../../testdata/quiet/lexer/l06.c",
			toks: []token.Token{
				{
					Kind: token.Comment,
					Val:  "/* This file contains examples of various types of white space and comments. */",
					Pos:  0,
				},
				{
					Kind: token.Comment,
					Val:  "// A blank (32)",
					Pos:  80,
				},
				{
					Kind: token.Comment,
					Val:  "// A new-line (10)",
					Pos:  98,
				},
				{
					Kind: token.Comment,
					Val:  "// a carriage-return (13)",
					Pos:  118,
				},
				{
					Kind: token.Comment,
					Val:  "// form-feed (12)",
					Pos:  146,
				},
				{
					Kind: token.Comment,
					Val:  "// tab(9)",
					Pos:  166,
				},
				{
					Kind: token.Comment,
					Val:  "/* In uC, each file must contain at least one declaration */",
					Pos:  178,
				},
				{
					Kind: token.Comment,
					Val:  "/* a comment... */",
					Pos:  240,
				},
				{
					Kind: token.Ident,
					Val:  "int",
					Pos:  259,
				},
				{
					Kind: token.Ident,
					Val:  "main",
					Pos:  263,
				},
				{
					Kind: token.Lparen,
					Val:  "(",
					Pos:  267,
				},
				{
					Kind: token.Comment,
					Val:  "/* ...and another */",
					Pos:  269,
				},
				{
					Kind: token.Ident,
					Val:  "void",
					Pos:  297,
				},
				{
					Kind: token.Rparen,
					Val:  ")",
					Pos:  302,
				},
				{
					Kind: token.Lbrace,
					Val:  "{",
					Pos:  304,
				},
				{
					Kind: token.Semicolon,
					Val:  ";",
					Pos:  306,
				},
				{
					Kind: token.Rbrace,
					Val:  "}",
					Pos:  308,
				},
				{
					Kind: token.EOF,
					Val:  "",
					Pos:  313,
				},
			},
		},

		// Test encodings.
		{
			path: "../../testdata/encoding/lexer/big5.c",
			toks: []token.Token{
				{
					Kind: token.Comment,
					Val:  "/* \xA5\x40\xAC\xC9\xB1\x7A\xA6\x6E */",
					Pos:  0,
				},
				{
					Kind: token.Ident,
					Val:  "int",
					Pos:  15,
				},
				{
					Kind: token.Ident,
					Val:  "a",
					Pos:  19,
				},
				{
					Kind: token.Semicolon,
					Val:  ";",
					Pos:  20,
				},
				{
					Kind: token.EOF,
					Val:  "",
					Pos:  22,
				},
			},
		},
		{
			path: "../../testdata/encoding/lexer/iso-8859-1.c",
			toks: []token.Token{
				{
					Kind: token.Comment,
					Val:  "// Hej v\xE4rlden!",
					Pos:  0,
				},
				{
					Kind: token.Ident,
					Val:  "int",
					Pos:  16,
				},
				{
					Kind: token.Ident,
					Val:  "b",
					Pos:  20,
				},
				{
					Kind: token.Semicolon,
					Val:  ";",
					Pos:  21,
				},
				{
					Kind: token.EOF,
					Val:  "",
					Pos:  23,
				},
			},
		},
		{
			path: "../../testdata/encoding/lexer/utf-8.c",
			toks: encTokens,
		},
		{
			path: "../../testdata/encoding/lexer/utf-8_bom.c",
			toks: encTokens,
		},
		{
			path: "../../testdata/encoding/lexer/utf-16be_bom.c",
			toks: encTokens,
		},
		{
			path: "../../testdata/encoding/lexer/utf-16le_bom.c",
			toks: encTokens,
		},

		// Extra tests.
		{
			path: "../../testdata/extra/lexer/invalid-esc.c",
			toks: []token.Token{
				{
					Kind: token.Comment,
					Val:  "// Test invalid escape sequence.",
					Pos:  0,
				},
				{
					Kind: token.Ident,
					Val:  "char",
					Pos:  33,
				},
				{
					Kind: token.Ident,
					Val:  "c",
					Pos:  38,
				},
				{
					Kind: token.Assign,
					Val:  "=",
					Pos:  40,
				},
				{
					Kind: token.Error,
					Val:  `unknown escape sequence '\q'`,
					Pos:  43,
				},
				{
					Kind: token.Error,
					Val:  `unexpected U+005C '\'`,
					Pos:  43,
				},
				{
					Kind: token.Ident,
					Val:  "q",
					Pos:  44,
				},
				{
					Kind: token.Error,
					Val:  "unterminated character literal",
					Pos:  45,
				},
				{
					Kind: token.Semicolon,
					Val:  ";",
					Pos:  46,
				},
				{
					Kind: token.EOF,
					Val:  "",
					Pos:  48,
				},
			},
		},
		{
			path: "../../testdata/extra/lexer/multibyte-char-lit.c",
			toks: []token.Token{
				{
					Kind: token.Comment,
					Val:  "// Invalid use of multibyte character literal.",
					Pos:  0,
				},
				{
					Kind: token.Ident,
					Val:  "char",
					Pos:  47,
				},
				{
					Kind: token.Ident,
					Val:  "c",
					Pos:  52,
				},
				{
					Kind: token.Assign,
					Val:  "=",
					Pos:  54,
				},
				{
					Kind: token.Error,
					Val:  "character U+00B5 'µ' too large for enclosing character literal type",
					Pos:  57,
				},
				{
					Kind: token.Error,
					Val:  "unexpected U+00B5 'µ'",
					Pos:  57,
				},
				{
					Kind: token.Error,
					Val:  "unterminated character literal",
					Pos:  59,
				},
				{
					Kind: token.Semicolon,
					Val:  ";",
					Pos:  60,
				},
				{
					Kind: token.EOF,
					Val:  "",
					Pos:  62,
				},
			},
		},
		{
			path: "../../testdata/extra/lexer/illegal-utf8-char-lit.c",
			toks: []token.Token{
				{
					Kind: token.Comment,
					Val:  "// Illegal UTF-8 encoding in character literal.",
					Pos:  0,
				},
				{
					Kind: token.Ident,
					Val:  "char",
					Pos:  48,
				},
				{
					Kind: token.Ident,
					Val:  "c",
					Pos:  53,
				},
				{
					Kind: token.Assign,
					Val:  "=",
					Pos:  55,
				},
				{
					Kind: token.Error,
					Val:  "illegal UTF-8 encoding",
					Pos:  58,
				},
				{
					Kind: token.Error,
					Val:  "illegal UTF-8 encoding",
					Pos:  58,
				},
				{
					Kind: token.Error,
					Val:  "unterminated character literal",
					Pos:  59,
				},
				{
					Kind: token.Semicolon,
					Val:  ";",
					Pos:  60,
				},
				{
					Kind: token.EOF,
					Val:  "",
					Pos:  62,
				},
			},
		},
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
				t.Errorf("%s: invalid number of tokens; expected %d tokens, got %d", g.path, len(g.toks), len(tokens))
				break
			}
			got := tokens[j]
			if j >= len(g.toks) {
				t.Errorf("%s: invalid number of tokens; expected %d tokens, got %d", g.path, len(g.toks), len(tokens))
				break
			}
			if want := g.toks[j]; got != want {
				t.Errorf("%s: token %d mismatch; expected %#v, got %#v", g.path, j, want, got)
			}
			if got.Kind == token.EOF {
				if j != len(g.toks)-1 {
					t.Errorf("%s: invalid number of tokens; expected %d tokens, got %d", g.path, len(g.toks), j+1)
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
	b.SetBytes(int64(len(src)))
	for i := 0; i < b.N; i++ {
		tokens := lexer.ParseString(src)
		for _, tok := range tokens {
			if tok.Kind == token.EOF {
				break
			}
		}
	}
}
