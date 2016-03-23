package lexer_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"testing"

	"github.com/mewmew/uc/uc/gocc/lexer"
	"github.com/mewmew/uc/uc/gocc/token"
)

func TestLexer(t *testing.T) {
	var golden = []struct {
		path string
		toks []*token.Token
	}{
		{
			path: "../../testdata/incorrect/lexer/bad.c",
			toks: []*token.Token{
				{
					Type: token.TokMap.Type("comment"),
					Lit:  []byte("/*\t\t\t\t\t\t\t||\n**\tFile for testing lexical analysis\t\t||\n**\t\t\t\t\t\t\t||\n**\tThis file is 'lexically incorrect'.\t\t||\n*/"),
					Pos:  token.Pos{Offset: 0},
				},
				{
					Type: token.TokMap.Type("ident"),
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
					Type: token.TokMap.Type("ident"),
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
					Lit:  []byte(""),
					Pos:  token.Pos{Offset: 221},
				},
			},
		},

		{
			path: "../../testdata/incorrect/lexer/good.c",
			toks: []*token.Token{
				{
					Type: token.TokMap.Type("comment"),
					Lit:  []byte("/*\t\t\t\t\t\t\t||\n**\tFile for testing lexical analysis\t\t||\n**\t\t\t\t\t\t\t||\n**\tThis file would confuse a parser, but\n        is 'lexically correct'.\t\t                ||\n*/"),
					Pos:  token.Pos{Offset: 1},
				},
				{
					Type: token.TokMap.Type("comment"),
					Lit:  []byte("/* ** / ** */"),
					Pos:  token.Pos{Offset: 163},
				},
				{
					Type: token.TokMap.Type("comment"),
					Lit:  []byte("// Simple tokens and single characters:\n"),
					Pos:  token.Pos{Offset: 179},
				},
				{
					Type: token.TokMap.Type("{"),
					Lit:  []byte("{"),
					Pos:  token.Pos{Offset: 220},
				},
				{
					Type: token.TokMap.Type("}"),
					Lit:  []byte("}"),
					Pos:  token.Pos{Offset: 222},
				},
				{
					Type: token.TokMap.Type("comment"),
					Lit:  []byte("// until end-of-line comment\n"),
					Pos:  token.Pos{Offset: 248},
				},
				{
					Type: token.TokMap.Type("if"),
					Lit:  []byte("if"),
					Pos:  token.Pos{Offset: 277},
				},
				{
					Type: token.TokMap.Type("else"),
					Lit:  []byte("else"),
					Pos:  token.Pos{Offset: 280},
				},
				{
					Type: token.TokMap.Type("while"),
					Lit:  []byte("while"),
					Pos:  token.Pos{Offset: 285},
				},
				{
					Type: token.TokMap.Type("comment"),
					Lit:  []byte("/* normal comment */"),
					Pos:  token.Pos{Offset: 311},
				},
				{
					Type: token.TokMap.Type("return"),
					Lit:  []byte("return"),
					Pos:  token.Pos{Offset: 332},
				},
				{
					Type: token.TokMap.Type("&&"),
					Lit:  []byte("&&"),
					Pos:  token.Pos{Offset: 339},
				},
				{
					Type: token.TokMap.Type("=="),
					Lit:  []byte("=="),
					Pos:  token.Pos{Offset: 342},
				},
				{
					Type: token.TokMap.Type("!="),
					Lit:  []byte("!="),
					Pos:  token.Pos{Offset: 345},
				},
				{
					Type: token.TokMap.Type("<="),
					Lit:  []byte("<="),
					Pos:  token.Pos{Offset: 348},
				},
				{
					Type: token.TokMap.Type(">="),
					Lit:  []byte(">="),
					Pos:  token.Pos{Offset: 351},
				},
				{
					Type: token.TokMap.Type("ident"),
					Lit:  []byte("char"),
					Pos:  token.Pos{Offset: 354},
				},
				{
					Type: token.TokMap.Type("ident"),
					Lit:  []byte("int"),
					Pos:  token.Pos{Offset: 359},
				},
				{
					Type: token.TokMap.Type("ident"),
					Lit:  []byte("void"),
					Pos:  token.Pos{Offset: 363},
				},
				{
					Type: token.TokMap.Type("+"),
					Lit:  []byte("+"),
					Pos:  token.Pos{Offset: 371},
				},
				{
					Type: token.TokMap.Type("-"),
					Lit:  []byte("-"),
					Pos:  token.Pos{Offset: 373},
				},
				{
					Type: token.TokMap.Type("*"),
					Lit:  []byte("*"),
					Pos:  token.Pos{Offset: 375},
				},
				{
					Type: token.TokMap.Type("/"),
					Lit:  []byte("/"),
					Pos:  token.Pos{Offset: 377},
				},
				{
					Type: token.TokMap.Type("<"),
					Lit:  []byte("<"),
					Pos:  token.Pos{Offset: 380},
				},
				{
					Type: token.TokMap.Type(">"),
					Lit:  []byte(">"),
					Pos:  token.Pos{Offset: 382},
				},
				{
					Type: token.TokMap.Type("="),
					Lit:  []byte("="),
					Pos:  token.Pos{Offset: 384},
				},
				{
					Type: token.TokMap.Type("("),
					Lit:  []byte("("),
					Pos:  token.Pos{Offset: 387},
				},
				{
					Type: token.TokMap.Type(","),
					Lit:  []byte(","),
					Pos:  token.Pos{Offset: 388},
				},
				{
					Type: token.TokMap.Type(";"),
					Lit:  []byte(";"),
					Pos:  token.Pos{Offset: 389},
				},
				{
					Type: token.TokMap.Type(")"),
					Lit:  []byte(")"),
					Pos:  token.Pos{Offset: 390},
				},
				{
					Type: token.TokMap.Type("["),
					Lit:  []byte("["),
					Pos:  token.Pos{Offset: 392},
				},
				{
					Type: token.TokMap.Type("ident"),
					Lit:  []byte("was"),
					Pos:  token.Pos{Offset: 393},
				},
				{
					Type: token.TokMap.Type("-"),
					Lit:  []byte("-"),
					Pos:  token.Pos{Offset: 396},
				},
				{
					Type: token.TokMap.Type("ident"),
					Lit:  []byte("colon"),
					Pos:  token.Pos{Offset: 397},
				},
				{
					Type: token.TokMap.Type("]"),
					Lit:  []byte("]"),
					Pos:  token.Pos{Offset: 402},
				},
				{
					Type: token.TokMap.Type("comment"),
					Lit:  []byte("/* Comment with bad tokens: _ || | ++ # @ ...  */"),
					Pos:  token.Pos{Offset: 406},
				},
				{
					Type: token.TokMap.Type("comment"),
					Lit:  []byte("// Ditto */ /* : _ || | ++ # @ ...  \n"),
					Pos:  token.Pos{Offset: 456},
				},
				{
					Type: token.TokMap.Type("comment"),
					Lit:  []byte("// Identifiers and numbers:\n"),
					Pos:  token.Pos{Offset: 493},
				},
				{
					Type: token.TokMap.Type("int_lit"),
					Lit:  []byte("17"),
					Pos:  token.Pos{Offset: 522},
				},
				{
					Type: token.TokMap.Type("-"),
					Lit:  []byte("-"),
					Pos:  token.Pos{Offset: 525},
				},
				{
					Type: token.TokMap.Type("int_lit"),
					Lit:  []byte("17"),
					Pos:  token.Pos{Offset: 526},
				},
				{
					Type: token.TokMap.Type("comment"),
					Lit:  []byte("// No floats? -17.17e17 -17.17E-17  \n"),
					Pos:  token.Pos{Offset: 529},
				},
				{
					Type: token.TokMap.Type("ident"),
					Lit:  []byte("ponderosa"),
					Pos:  token.Pos{Offset: 568},
				},
				{
					Type: token.TokMap.Type("ident"),
					Lit:  []byte("Black"),
					Pos:  token.Pos{Offset: 578},
				},
				{
					Type: token.TokMap.Type("ident"),
					Lit:  []byte("Steel"),
					Pos:  token.Pos{Offset: 584},
				},
				{
					Type: token.TokMap.Type("ident"),
					Lit:  []byte("PUMPKIN"),
					Pos:  token.Pos{Offset: 590},
				},
				{
					Type: token.TokMap.Type("ident"),
					Lit:  []byte("AfterMath"),
					Pos:  token.Pos{Offset: 598},
				},
				{
					Type: token.TokMap.Type("ident"),
					Lit:  []byte("aBBaoN"),
					Pos:  token.Pos{Offset: 608},
				},
				{
					Type: token.TokMap.Type("ident"),
					Lit:  []byte("faT"),
					Pos:  token.Pos{Offset: 615},
				},
				{
					Type: token.TokMap.Type("ident"),
					Lit:  []byte("TRacKs"),
					Pos:  token.Pos{Offset: 619},
				},
				{
					Type: token.TokMap.Type("ident"),
					Lit:  []byte("K9"),
					Pos:  token.Pos{Offset: 628},
				},
				{
					Type: token.TokMap.Type("ident"),
					Lit:  []byte("R23"),
					Pos:  token.Pos{Offset: 631},
				},
				{
					Type: token.TokMap.Type("ident"),
					Lit:  []byte("B52"),
					Pos:  token.Pos{Offset: 635},
				},
				{
					Type: token.TokMap.Type("ident"),
					Lit:  []byte("Track15"),
					Pos:  token.Pos{Offset: 639},
				},
				{
					Type: token.TokMap.Type("ident"),
					Lit:  []byte("not4money"),
					Pos:  token.Pos{Offset: 647},
				},
				{
					Type: token.TokMap.Type("int_lit"),
					Lit:  []byte("378"),
					Pos:  token.Pos{Offset: 657},
				},
				{
					Type: token.TokMap.Type("ident"),
					Lit:  []byte("WHOIS666999SIOHM"),
					Pos:  token.Pos{Offset: 661},
				},
				{
					Type: token.TokMap.Type("ident"),
					Lit:  []byte("was"),
					Pos:  token.Pos{Offset: 687},
				},
				{
					Type: token.TokMap.Type("-"),
					Lit:  []byte("-"),
					Pos:  token.Pos{Offset: 690},
				},
				{
					Type: token.TokMap.Type("ident"),
					Lit:  []byte("floating"),
					Pos:  token.Pos{Offset: 691},
				},
				{
					Type: token.TokMap.Type("-"),
					Lit:  []byte("-"),
					Pos:  token.Pos{Offset: 699},
				},
				{
					Type: token.TokMap.Type("ident"),
					Lit:  []byte("point"),
					Pos:  token.Pos{Offset: 700},
				},
				{
					Type: token.TokMap.Type("-"),
					Lit:  []byte("-"),
					Pos:  token.Pos{Offset: 705},
				},
				{
					Type: token.TokMap.Type("ident"),
					Lit:  []byte("number"),
					Pos:  token.Pos{Offset: 706},
				},
				{
					Type: token.TokMap.Type("comment"),
					Lit:  []byte("/* The following 'trap' should be correctly handled:\n\n\t\t* \"2die4U\" consists of the number '2' and the\n\t\t  identifier 'die4U'.\n*/"),
					Pos:  token.Pos{Offset: 714},
				},
				{
					Type: token.TokMap.Type("int_lit"),
					Lit:  []byte("2"),
					Pos:  token.Pos{Offset: 852},
				},
				{
					Type: token.TokMap.Type("ident"),
					Lit:  []byte("die4U"),
					Pos:  token.Pos{Offset: 853},
				},
				{
					Type: token.TokMap.Type("comment"),
					Lit:  []byte("//|| The following should all be regarded as identifiers:\n"),
					Pos:  token.Pos{Offset: 860},
				},
				{
					Type: token.TokMap.Type("ident"),
					Lit:  []byte("Function"),
					Pos:  token.Pos{Offset: 920},
				},
				{
					Type: token.TokMap.Type("ident"),
					Lit:  []byte("PrOceDuRE"),
					Pos:  token.Pos{Offset: 929},
				},
				{
					Type: token.TokMap.Type("ident"),
					Lit:  []byte("begIN"),
					Pos:  token.Pos{Offset: 939},
				},
				{
					Type: token.TokMap.Type("ident"),
					Lit:  []byte("eNd"),
					Pos:  token.Pos{Offset: 945},
				},
				{
					Type: token.TokMap.Type("ident"),
					Lit:  []byte("PrinT"),
					Pos:  token.Pos{Offset: 949},
				},
				{
					Type: token.TokMap.Type("ident"),
					Lit:  []byte("rEad"),
					Pos:  token.Pos{Offset: 955},
				},
				{
					Type: token.TokMap.Type("ident"),
					Lit:  []byte("iF"),
					Pos:  token.Pos{Offset: 960},
				},
				{
					Type: token.TokMap.Type("ident"),
					Lit:  []byte("THen"),
					Pos:  token.Pos{Offset: 963},
				},
				{
					Type: token.TokMap.Type("ident"),
					Lit:  []byte("StaTic"),
					Pos:  token.Pos{Offset: 968},
				},
				{
					Type: token.TokMap.Type("ident"),
					Lit:  []byte("ElSe"),
					Pos:  token.Pos{Offset: 976},
				},
				{
					Type: token.TokMap.Type("ident"),
					Lit:  []byte("wHilE"),
					Pos:  token.Pos{Offset: 981},
				},
				{
					Type: token.TokMap.Type("ident"),
					Lit:  []byte("Do"),
					Pos:  token.Pos{Offset: 987},
				},
				{
					Type: token.TokMap.Type("ident"),
					Lit:  []byte("reTurN"),
					Pos:  token.Pos{Offset: 990},
				},
				{
					Type: token.TokMap.Type("ident"),
					Lit:  []byte("noT"),
					Pos:  token.Pos{Offset: 997},
				},
				{
					Type: token.TokMap.Type("ident"),
					Lit:  []byte("AnD"),
					Pos:  token.Pos{Offset: 1001},
				},
				{
					Type: token.TokMap.Type("ident"),
					Lit:  []byte("OR"),
					Pos:  token.Pos{Offset: 1005},
				},
				{
					Type: token.TokMap.Type("ident"),
					Lit:  []byte("TrUE"),
					Pos:  token.Pos{Offset: 1008},
				},
				{
					Type: token.TokMap.Type("ident"),
					Lit:  []byte("bOOl"),
					Pos:  token.Pos{Offset: 1013},
				},
				{
					Type: token.TokMap.Type("ident"),
					Lit:  []byte("FalsE"),
					Pos:  token.Pos{Offset: 1018},
				},
				{
					Type: token.TokMap.Type("ident"),
					Lit:  []byte("sizE"),
					Pos:  token.Pos{Offset: 1024},
				},
				{
					Type: token.TokMap.Type("comment"),
					Lit:  []byte("// It is legal to end the code like this, without an ending newline."), // TODO: Figure out how to handle line comments, ending with EOF.
					Pos:  token.Pos{Offset: 1031},
				},
				{
					Type: token.TokMap.Type("$"),
					Lit:  []byte(""),
					Pos:  token.Pos{Offset: 1099},
				},
			},
		},

		{
			path: "../../testdata/incorrect/lexer/long-char.c",
			toks: []*token.Token{
				{
					Type: token.TokMap.Type("ident"),
					Lit:  []byte("int"),
					Pos:  token.Pos{Offset: 0},
				},
				{
					Type: token.TokMap.Type("ident"),
					Lit:  []byte("main"),
					Pos:  token.Pos{Offset: 4},
				},
				{
					Type: token.TokMap.Type("("),
					Lit:  []byte("("),
					Pos:  token.Pos{Offset: 8},
				},
				{
					Type: token.TokMap.Type("ident"),
					Lit:  []byte("void"),
					Pos:  token.Pos{Offset: 9},
				},
				{
					Type: token.TokMap.Type(")"),
					Lit:  []byte(")"),
					Pos:  token.Pos{Offset: 13},
				},
				{
					Type: token.TokMap.Type("{"),
					Lit:  []byte("{"),
					Pos:  token.Pos{Offset: 15},
				},
				{
					Type: token.TokMap.Type("ident"),
					Lit:  []byte("char"),
					Pos:  token.Pos{Offset: 19},
				},
				{
					Type: token.TokMap.Type("ident"),
					Lit:  []byte("c"),
					Pos:  token.Pos{Offset: 24},
				},
				{
					Type: token.TokMap.Type(";"),
					Lit:  []byte(";"),
					Pos:  token.Pos{Offset: 25},
				},
				{
					Type: token.TokMap.Type("ident"),
					Lit:  []byte("c"),
					Pos:  token.Pos{Offset: 29},
				},
				{
					Type: token.TokMap.Type("="),
					Lit:  []byte("="),
					Pos:  token.Pos{Offset: 31},
				},
				{
					Type: token.TokMap.Type("int_lit"), // TODO: Figure out if char_lit and int_lit should be separate tokens.
					Lit:  []byte("'c'"),
					Pos:  token.Pos{Offset: 33},
				},
				{
					Type: token.TokMap.Type(";"),
					Lit:  []byte(";"),
					Pos:  token.Pos{Offset: 36},
				},
				{
					Type: token.TokMap.Type("comment"),
					Lit:  []byte("// OK\n"),
					Pos:  token.Pos{Offset: 38},
				},
				{
					Type: token.TokMap.Type("ident"),
					Lit:  []byte("c"),
					Pos:  token.Pos{Offset: 46},
				},
				{
					Type: token.TokMap.Type("="),
					Lit:  []byte("="),
					Pos:  token.Pos{Offset: 48},
				},
				{
					Type: token.TokMap.Type("INVALID"),
					Lit:  []byte("'cc"),
					Pos:  token.Pos{Offset: 50},
				},
				{
					Type: token.TokMap.Type("INVALID"),
					Lit:  []byte("'; "),
					Pos:  token.Pos{Offset: 53},
				},
				{
					Type: token.TokMap.Type("comment"),
					Lit:  []byte("// Not OK\n"),
					Pos:  token.Pos{Offset: 56},
				},
				{
					Type: token.TokMap.Type("}"),
					Lit:  []byte("}"),
					Pos:  token.Pos{Offset: 66},
				},
				{
					Type: token.TokMap.Type("$"),
					Lit:  []byte(""),
					Pos:  token.Pos{Offset: 68},
				},
			},
		},

		{
			path: "../../testdata/incorrect/lexer/ugly.c",
			toks: []*token.Token{
				{
					Type: token.TokMap.Type("INVALID"),
					Lit:  []byte("|"),
					Pos:  token.Pos{Offset: 0},
				},
				{
					Type: token.TokMap.Type("+"),
					Lit:  []byte("+"),
					Pos:  token.Pos{Offset: 1},
				},
				{
					Type: token.TokMap.Type("INVALID"),
					Lit:  []byte("\x01"),
					Pos:  token.Pos{Offset: 2},
				},
				{
					Type: token.TokMap.Type("INVALID"),
					Lit:  []byte("\x02"),
					Pos:  token.Pos{Offset: 3},
				},
				{
					Type: token.TokMap.Type("INVALID"),
					Lit:  []byte("\x03"),
					Pos:  token.Pos{Offset: 4},
				},
				{
					Type: token.TokMap.Type("INVALID"),
					Lit:  []byte("\x04"),
					Pos:  token.Pos{Offset: 5},
				},
				{
					Type: token.TokMap.Type("INVALID"),
					Lit:  []byte("\x05"),
					Pos:  token.Pos{Offset: 6},
				},
				{
					Type: token.TokMap.Type("INVALID"),
					Lit:  []byte("\x06"),
					Pos:  token.Pos{Offset: 7},
				},
				{
					Type: token.TokMap.Type("INVALID"),
					Lit:  []byte("\x07"),
					Pos:  token.Pos{Offset: 8},
				},
				{
					Type: token.TokMap.Type("INVALID"),
					Lit:  []byte("\x08"),
					Pos:  token.Pos{Offset: 9},
				},
				{
					Type: token.TokMap.Type("INVALID"),
					Lit:  []byte("\x0E"),
					Pos:  token.Pos{Offset: 15},
				},
				{
					Type: token.TokMap.Type("INVALID"),
					Lit:  []byte("\x0F"),
					Pos:  token.Pos{Offset: 16},
				},
				{
					Type: token.TokMap.Type("INVALID"),
					Lit:  []byte("\x10"),
					Pos:  token.Pos{Offset: 17},
				},
				{
					Type: token.TokMap.Type("INVALID"),
					Lit:  []byte("\x11"),
					Pos:  token.Pos{Offset: 18},
				},
				{
					Type: token.TokMap.Type("INVALID"),
					Lit:  []byte("\x12"),
					Pos:  token.Pos{Offset: 19},
				},
				{
					Type: token.TokMap.Type("INVALID"),
					Lit:  []byte("\x13"),
					Pos:  token.Pos{Offset: 20},
				},
				{
					Type: token.TokMap.Type("INVALID"),
					Lit:  []byte("\x14"),
					Pos:  token.Pos{Offset: 21},
				},
				{
					Type: token.TokMap.Type("INVALID"),
					Lit:  []byte("\x15"),
					Pos:  token.Pos{Offset: 22},
				},
				{
					Type: token.TokMap.Type("INVALID"),
					Lit:  []byte("\x16"),
					Pos:  token.Pos{Offset: 23},
				},
				{
					Type: token.TokMap.Type("INVALID"),
					Lit:  []byte("\x17"),
					Pos:  token.Pos{Offset: 24},
				},
				{
					Type: token.TokMap.Type("INVALID"),
					Lit:  []byte("\x18"),
					Pos:  token.Pos{Offset: 25},
				},
				{
					Type: token.TokMap.Type("INVALID"),
					Lit:  []byte("\x19"),
					Pos:  token.Pos{Offset: 26},
				},
				{
					Type: token.TokMap.Type("INVALID"),
					Lit:  []byte("\x1A"),
					Pos:  token.Pos{Offset: 27},
				},
				{
					Type: token.TokMap.Type("INVALID"),
					Lit:  []byte("\x1B"),
					Pos:  token.Pos{Offset: 28},
				},
				{
					Type: token.TokMap.Type("INVALID"),
					Lit:  []byte("\x1C"),
					Pos:  token.Pos{Offset: 29},
				},
				{
					Type: token.TokMap.Type("INVALID"),
					Lit:  []byte("\x1D"),
					Pos:  token.Pos{Offset: 30},
				},
				{
					Type: token.TokMap.Type("INVALID"),
					Lit:  []byte("\x1E"),
					Pos:  token.Pos{Offset: 31},
				},
				{
					Type: token.TokMap.Type("INVALID"),
					Lit:  []byte("\x1F"),
					Pos:  token.Pos{Offset: 32},
				},
				{
					Type: token.TokMap.Type("!"),
					Lit:  []byte("!"),
					Pos:  token.Pos{Offset: 34},
				},
				{
					Type: token.TokMap.Type("INVALID"),
					Lit:  []byte(`"`),
					Pos:  token.Pos{Offset: 35},
				},
				{
					Type: token.TokMap.Type("INVALID"),
					Lit:  []byte("#"),
					Pos:  token.Pos{Offset: 36},
				},
				{
					Type: token.TokMap.Type("INVALID"),
					Lit:  []byte("$"),
					Pos:  token.Pos{Offset: 37},
				},
				{
					Type: token.TokMap.Type("INVALID"),
					Lit:  []byte("%"),
					Pos:  token.Pos{Offset: 38},
				},
				{
					Type: token.TokMap.Type("INVALID"),
					Lit:  []byte("&'"),
					Pos:  token.Pos{Offset: 39},
				},
				{
					Type: token.TokMap.Type("("),
					Lit:  []byte("("),
					Pos:  token.Pos{Offset: 41},
				},
				{
					Type: token.TokMap.Type(")"),
					Lit:  []byte(")"),
					Pos:  token.Pos{Offset: 42},
				},
				{
					Type: token.TokMap.Type("*"),
					Lit:  []byte("*"),
					Pos:  token.Pos{Offset: 43},
				},
				{
					Type: token.TokMap.Type("+"),
					Lit:  []byte("+"),
					Pos:  token.Pos{Offset: 44},
				},
				{
					Type: token.TokMap.Type(","),
					Lit:  []byte(","),
					Pos:  token.Pos{Offset: 46},
				},
				{
					Type: token.TokMap.Type("-"),
					Lit:  []byte("-"),
					Pos:  token.Pos{Offset: 47},
				},
				{
					Type: token.TokMap.Type("INVALID"),
					Lit:  []byte("."),
					Pos:  token.Pos{Offset: 48},
				},
				{
					Type: token.TokMap.Type("/"),
					Lit:  []byte("/"),
					Pos:  token.Pos{Offset: 49},
				},
				{
					Type: token.TokMap.Type("int_lit"),
					Lit:  []byte("0123456789"), // TODO: When octal integer literals have been implemented, fail accordingly.
					Pos:  token.Pos{Offset: 50},
				},
				{
					Type: token.TokMap.Type("INVALID"),
					Lit:  []byte(":"),
					Pos:  token.Pos{Offset: 60},
				},
				{
					Type: token.TokMap.Type(";"),
					Lit:  []byte(";"),
					Pos:  token.Pos{Offset: 61},
				},
				{
					Type: token.TokMap.Type("<="),
					Lit:  []byte("<="),
					Pos:  token.Pos{Offset: 62},
				},
				{
					Type: token.TokMap.Type(">"),
					Lit:  []byte(">"),
					Pos:  token.Pos{Offset: 64},
				},
				{
					Type: token.TokMap.Type("INVALID"),
					Lit:  []byte("?"),
					Pos:  token.Pos{Offset: 65},
				},
				{
					Type: token.TokMap.Type("INVALID"),
					Lit:  []byte("@"),
					Pos:  token.Pos{Offset: 66},
				},
				{
					Type: token.TokMap.Type("ident"),
					Lit:  []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ"),
					Pos:  token.Pos{Offset: 67},
				},
				{
					Type: token.TokMap.Type("["),
					Lit:  []byte("["),
					Pos:  token.Pos{Offset: 93},
				},
				{
					Type: token.TokMap.Type("INVALID"),
					Lit:  []byte(`\`),
					Pos:  token.Pos{Offset: 94},
				},
				{
					Type: token.TokMap.Type("]"),
					Lit:  []byte("]"),
					Pos:  token.Pos{Offset: 95},
				},
				{
					Type: token.TokMap.Type("INVALID"),
					Lit:  []byte("^"),
					Pos:  token.Pos{Offset: 96},
				},
				{
					Type: token.TokMap.Type("ident"),
					Lit:  []byte("_"),
					Pos:  token.Pos{Offset: 97},
				},
				{
					Type: token.TokMap.Type("INVALID"),
					Lit:  []byte("`"),
					Pos:  token.Pos{Offset: 98},
				},
				{
					Type: token.TokMap.Type("ident"),
					Lit:  []byte("abcdefghijklmnopqrstuvwxyz"),
					Pos:  token.Pos{Offset: 99},
				},
				{
					Type: token.TokMap.Type("{"),
					Lit:  []byte("{"),
					Pos:  token.Pos{Offset: 125},
				},
				{
					Type: token.TokMap.Type("INVALID"),
					Lit:  []byte("|"),
					Pos:  token.Pos{Offset: 126},
				},
				{
					Type: token.TokMap.Type("}"),
					Lit:  []byte("}"),
					Pos:  token.Pos{Offset: 127},
				},
				{
					Type: token.TokMap.Type("INVALID"),
					Lit:  []byte("~"),
					Pos:  token.Pos{Offset: 128},
				},
				{
					Type: token.TokMap.Type("INVALID"),
					Lit:  []byte("\x7F"),
					Pos:  token.Pos{Offset: 129},
				},
				{
					Type: token.TokMap.Type("INVALID"),
					Lit:  []byte("\x80"),
					Pos:  token.Pos{Offset: 130},
				},
				{
					Type: token.TokMap.Type("INVALID"),
					Lit:  []byte("\x81"),
					Pos:  token.Pos{Offset: 131},
				},
				{
					Type: token.TokMap.Type("INVALID"),
					Lit:  []byte("\x82"),
					Pos:  token.Pos{Offset: 132},
				},
				{
					Type: token.TokMap.Type("INVALID"),
					Lit:  []byte("\x83"),
					Pos:  token.Pos{Offset: 133},
				},
				{
					Type: token.TokMap.Type("INVALID"),
					Lit:  []byte("\x84"),
					Pos:  token.Pos{Offset: 134},
				},
				{
					Type: token.TokMap.Type("INVALID"),
					Lit:  []byte("\x85"),
					Pos:  token.Pos{Offset: 135},
				},
				{
					Type: token.TokMap.Type("INVALID"),
					Lit:  []byte("\x86"),
					Pos:  token.Pos{Offset: 136},
				},
				{
					Type: token.TokMap.Type("INVALID"),
					Lit:  []byte("\x87"),
					Pos:  token.Pos{Offset: 137},
				},
				{
					Type: token.TokMap.Type("INVALID"),
					Lit:  []byte("\x88"),
					Pos:  token.Pos{Offset: 138},
				},
				{
					Type: token.TokMap.Type("INVALID"),
					Lit:  []byte("\x89"),
					Pos:  token.Pos{Offset: 139},
				},
				{
					Type: token.TokMap.Type("INVALID"),
					Lit:  []byte("\x8A"),
					Pos:  token.Pos{Offset: 140},
				},
				{
					Type: token.TokMap.Type("INVALID"),
					Lit:  []byte("\x8B"),
					Pos:  token.Pos{Offset: 141},
				},
				{
					Type: token.TokMap.Type("INVALID"),
					Lit:  []byte("\x8C"),
					Pos:  token.Pos{Offset: 142},
				},
				{
					Type: token.TokMap.Type("INVALID"),
					Lit:  []byte("\x8D"),
					Pos:  token.Pos{Offset: 143},
				},
				{
					Type: token.TokMap.Type("INVALID"),
					Lit:  []byte("\x8E"),
					Pos:  token.Pos{Offset: 144},
				},
				{
					Type: token.TokMap.Type("INVALID"),
					Lit:  []byte("\x8F"),
					Pos:  token.Pos{Offset: 145},
				},
				{
					Type: token.TokMap.Type("INVALID"),
					Lit:  []byte("\x90"),
					Pos:  token.Pos{Offset: 146},
				},
				{
					Type: token.TokMap.Type("INVALID"),
					Lit:  []byte("\x91"),
					Pos:  token.Pos{Offset: 147},
				},
				{
					Type: token.TokMap.Type("INVALID"),
					Lit:  []byte("\x92"),
					Pos:  token.Pos{Offset: 148},
				},
				{
					Type: token.TokMap.Type("INVALID"),
					Lit:  []byte("\x93"),
					Pos:  token.Pos{Offset: 149},
				},
				{
					Type: token.TokMap.Type("INVALID"),
					Lit:  []byte("\x94"),
					Pos:  token.Pos{Offset: 150},
				},
				{
					Type: token.TokMap.Type("INVALID"),
					Lit:  []byte("\x95"),
					Pos:  token.Pos{Offset: 151},
				},
				{
					Type: token.TokMap.Type("INVALID"),
					Lit:  []byte("\x96"),
					Pos:  token.Pos{Offset: 152},
				},
				{
					Type: token.TokMap.Type("INVALID"),
					Lit:  []byte("\x97"),
					Pos:  token.Pos{Offset: 153},
				},
				{
					Type: token.TokMap.Type("INVALID"),
					Lit:  []byte("\x98"),
					Pos:  token.Pos{Offset: 154},
				},
				{
					Type: token.TokMap.Type("INVALID"),
					Lit:  []byte("\x99"),
					Pos:  token.Pos{Offset: 155},
				},
				{
					Type: token.TokMap.Type("INVALID"),
					Lit:  []byte("\x9A"),
					Pos:  token.Pos{Offset: 156},
				},
				{
					Type: token.TokMap.Type("INVALID"),
					Lit:  []byte("\x9B"),
					Pos:  token.Pos{Offset: 157},
				},
				{
					Type: token.TokMap.Type("INVALID"),
					Lit:  []byte("\x9C"),
					Pos:  token.Pos{Offset: 158},
				},
				{
					Type: token.TokMap.Type("INVALID"),
					Lit:  []byte("\x9D"),
					Pos:  token.Pos{Offset: 159},
				},
				{
					Type: token.TokMap.Type("INVALID"),
					Lit:  []byte("\x9E"),
					Pos:  token.Pos{Offset: 160},
				},
				{
					Type: token.TokMap.Type("INVALID"),
					Lit:  []byte("\x9F"),
					Pos:  token.Pos{Offset: 161},
				},
				{
					Type: token.TokMap.Type("INVALID"),
					Lit:  []byte("\xA0"),
					Pos:  token.Pos{Offset: 162},
				},
				{
					Type: token.TokMap.Type("INVALID"),
					Lit:  []byte("\xA1"),
					Pos:  token.Pos{Offset: 163},
				},
				{
					Type: token.TokMap.Type("INVALID"),
					Lit:  []byte("\xA2"),
					Pos:  token.Pos{Offset: 164},
				},
				{
					Type: token.TokMap.Type("INVALID"),
					Lit:  []byte("\xA3"),
					Pos:  token.Pos{Offset: 165},
				},
				{
					Type: token.TokMap.Type("INVALID"),
					Lit:  []byte("\xA4"),
					Pos:  token.Pos{Offset: 166},
				},
				{
					Type: token.TokMap.Type("INVALID"),
					Lit:  []byte("\xA5"),
					Pos:  token.Pos{Offset: 167},
				},
				{
					Type: token.TokMap.Type("INVALID"),
					Lit:  []byte("\xA6"),
					Pos:  token.Pos{Offset: 168},
				},
				{
					Type: token.TokMap.Type("INVALID"),
					Lit:  []byte("\xA7"),
					Pos:  token.Pos{Offset: 169},
				},
				{
					Type: token.TokMap.Type("INVALID"),
					Lit:  []byte("\xA8"),
					Pos:  token.Pos{Offset: 170},
				},
				{
					Type: token.TokMap.Type("INVALID"),
					Lit:  []byte("\xA9"),
					Pos:  token.Pos{Offset: 171},
				},
				{
					Type: token.TokMap.Type("INVALID"),
					Lit:  []byte("\xAA"),
					Pos:  token.Pos{Offset: 172},
				},
				{
					Type: token.TokMap.Type("INVALID"),
					Lit:  []byte("\xAB"),
					Pos:  token.Pos{Offset: 173},
				},
				{
					Type: token.TokMap.Type("INVALID"),
					Lit:  []byte("\xAC"),
					Pos:  token.Pos{Offset: 174},
				},
				{
					Type: token.TokMap.Type("INVALID"),
					Lit:  []byte("\xAD"),
					Pos:  token.Pos{Offset: 175},
				},
				{
					Type: token.TokMap.Type("INVALID"),
					Lit:  []byte("\xAE"),
					Pos:  token.Pos{Offset: 176},
				},
				{
					Type: token.TokMap.Type("INVALID"),
					Lit:  []byte("\xAF"),
					Pos:  token.Pos{Offset: 177},
				},
				{
					Type: token.TokMap.Type("INVALID"),
					Lit:  []byte("\xB0"),
					Pos:  token.Pos{Offset: 178},
				},
				{
					Type: token.TokMap.Type("INVALID"),
					Lit:  []byte("\xB1"),
					Pos:  token.Pos{Offset: 179},
				},
				{
					Type: token.TokMap.Type("INVALID"),
					Lit:  []byte("\xB2"),
					Pos:  token.Pos{Offset: 180},
				},
				{
					Type: token.TokMap.Type("INVALID"),
					Lit:  []byte("\xB3"),
					Pos:  token.Pos{Offset: 181},
				},
				{
					Type: token.TokMap.Type("INVALID"),
					Lit:  []byte("\xB4"),
					Pos:  token.Pos{Offset: 182},
				},
				{
					Type: token.TokMap.Type("INVALID"),
					Lit:  []byte("\xB5"),
					Pos:  token.Pos{Offset: 183},
				},
				{
					Type: token.TokMap.Type("INVALID"),
					Lit:  []byte("\xB6"),
					Pos:  token.Pos{Offset: 184},
				},
				{
					Type: token.TokMap.Type("INVALID"),
					Lit:  []byte("\xB7"),
					Pos:  token.Pos{Offset: 185},
				},
				{
					Type: token.TokMap.Type("INVALID"),
					Lit:  []byte("\xB8"),
					Pos:  token.Pos{Offset: 186},
				},
				{
					Type: token.TokMap.Type("INVALID"),
					Lit:  []byte("\xB9"),
					Pos:  token.Pos{Offset: 187},
				},
				{
					Type: token.TokMap.Type("INVALID"),
					Lit:  []byte("\xBA"),
					Pos:  token.Pos{Offset: 188},
				},
				{
					Type: token.TokMap.Type("INVALID"),
					Lit:  []byte("\xBB"),
					Pos:  token.Pos{Offset: 189},
				},
				{
					Type: token.TokMap.Type("INVALID"),
					Lit:  []byte("\xBC"),
					Pos:  token.Pos{Offset: 190},
				},
				{
					Type: token.TokMap.Type("INVALID"),
					Lit:  []byte("\xBD"),
					Pos:  token.Pos{Offset: 191},
				},
				{
					Type: token.TokMap.Type("INVALID"),
					Lit:  []byte("\xBE"),
					Pos:  token.Pos{Offset: 192},
				},
				{
					Type: token.TokMap.Type("INVALID"),
					Lit:  []byte("\xBF"),
					Pos:  token.Pos{Offset: 193},
				},
				{
					Type: token.TokMap.Type("INVALID"),
					Lit:  []byte("\xC0"),
					Pos:  token.Pos{Offset: 194},
				},
				{
					Type: token.TokMap.Type("INVALID"),
					Lit:  []byte("\xC1"),
					Pos:  token.Pos{Offset: 195},
				},
				{
					Type: token.TokMap.Type("INVALID"),
					Lit:  []byte("\xC2"),
					Pos:  token.Pos{Offset: 196},
				},
				{
					Type: token.TokMap.Type("INVALID"),
					Lit:  []byte("\xC3"),
					Pos:  token.Pos{Offset: 197},
				},
				{
					Type: token.TokMap.Type("INVALID"),
					Lit:  []byte("\xC4"),
					Pos:  token.Pos{Offset: 198},
				},
				{
					Type: token.TokMap.Type("INVALID"),
					Lit:  []byte("\xC5"),
					Pos:  token.Pos{Offset: 199},
				},
				{
					Type: token.TokMap.Type("INVALID"),
					Lit:  []byte("\xC6"),
					Pos:  token.Pos{Offset: 200},
				},
				{
					Type: token.TokMap.Type("INVALID"),
					Lit:  []byte("\xC7"),
					Pos:  token.Pos{Offset: 201},
				},
				{
					Type: token.TokMap.Type("INVALID"),
					Lit:  []byte("\xC8"),
					Pos:  token.Pos{Offset: 202},
				},
				{
					Type: token.TokMap.Type("INVALID"),
					Lit:  []byte("\xC9"),
					Pos:  token.Pos{Offset: 203},
				},
				{
					Type: token.TokMap.Type("INVALID"),
					Lit:  []byte("\xCA"),
					Pos:  token.Pos{Offset: 204},
				},
				{
					Type: token.TokMap.Type("INVALID"),
					Lit:  []byte("\xCB"),
					Pos:  token.Pos{Offset: 205},
				},
				{
					Type: token.TokMap.Type("INVALID"),
					Lit:  []byte("\xCC"),
					Pos:  token.Pos{Offset: 206},
				},
				{
					Type: token.TokMap.Type("INVALID"),
					Lit:  []byte("\xCD"),
					Pos:  token.Pos{Offset: 207},
				},
				{
					Type: token.TokMap.Type("INVALID"),
					Lit:  []byte("\xCE"),
					Pos:  token.Pos{Offset: 208},
				},
				{
					Type: token.TokMap.Type("INVALID"),
					Lit:  []byte("\xCF"),
					Pos:  token.Pos{Offset: 209},
				},
				{
					Type: token.TokMap.Type("INVALID"),
					Lit:  []byte("\xD0"),
					Pos:  token.Pos{Offset: 210},
				},
				{
					Type: token.TokMap.Type("INVALID"),
					Lit:  []byte("\xD1"),
					Pos:  token.Pos{Offset: 211},
				},
				{
					Type: token.TokMap.Type("INVALID"),
					Lit:  []byte("\xD2"),
					Pos:  token.Pos{Offset: 212},
				},
				{
					Type: token.TokMap.Type("INVALID"),
					Lit:  []byte("\xD3"),
					Pos:  token.Pos{Offset: 213},
				},
				{
					Type: token.TokMap.Type("INVALID"),
					Lit:  []byte("\xD4"),
					Pos:  token.Pos{Offset: 214},
				},
				{
					Type: token.TokMap.Type("INVALID"),
					Lit:  []byte("\xD5"),
					Pos:  token.Pos{Offset: 215},
				},
				{
					Type: token.TokMap.Type("INVALID"),
					Lit:  []byte("\xD6"),
					Pos:  token.Pos{Offset: 216},
				},
				{
					Type: token.TokMap.Type("INVALID"),
					Lit:  []byte("\xD7"),
					Pos:  token.Pos{Offset: 217},
				},
				{
					Type: token.TokMap.Type("INVALID"),
					Lit:  []byte("\xD8"),
					Pos:  token.Pos{Offset: 218},
				},
				{
					Type: token.TokMap.Type("INVALID"),
					Lit:  []byte("\xD9"),
					Pos:  token.Pos{Offset: 219},
				},
				{
					Type: token.TokMap.Type("INVALID"),
					Lit:  []byte("\xDA"),
					Pos:  token.Pos{Offset: 220},
				},
				{
					Type: token.TokMap.Type("INVALID"),
					Lit:  []byte("\xDB"),
					Pos:  token.Pos{Offset: 221},
				},
				{
					Type: token.TokMap.Type("INVALID"),
					Lit:  []byte("\xDC"),
					Pos:  token.Pos{Offset: 222},
				},
				{
					Type: token.TokMap.Type("INVALID"),
					Lit:  []byte("\xDD"),
					Pos:  token.Pos{Offset: 223},
				},
				{
					Type: token.TokMap.Type("INVALID"),
					Lit:  []byte("\xDE"),
					Pos:  token.Pos{Offset: 224},
				},
				{
					Type: token.TokMap.Type("INVALID"),
					Lit:  []byte("\xDF"),
					Pos:  token.Pos{Offset: 225},
				},
				{
					Type: token.TokMap.Type("INVALID"),
					Lit:  []byte("\xE0"),
					Pos:  token.Pos{Offset: 226},
				},
				{
					Type: token.TokMap.Type("INVALID"),
					Lit:  []byte("\xE1"),
					Pos:  token.Pos{Offset: 227},
				},
				{
					Type: token.TokMap.Type("INVALID"),
					Lit:  []byte("\xE2"),
					Pos:  token.Pos{Offset: 228},
				},
				{
					Type: token.TokMap.Type("INVALID"),
					Lit:  []byte("\xE3"),
					Pos:  token.Pos{Offset: 229},
				},
				{
					Type: token.TokMap.Type("INVALID"),
					Lit:  []byte("\xE4"),
					Pos:  token.Pos{Offset: 230},
				},
				{
					Type: token.TokMap.Type("INVALID"),
					Lit:  []byte("\xE5"),
					Pos:  token.Pos{Offset: 231},
				},
				{
					Type: token.TokMap.Type("INVALID"),
					Lit:  []byte("\xE6"),
					Pos:  token.Pos{Offset: 232},
				},
				{
					Type: token.TokMap.Type("INVALID"),
					Lit:  []byte("\xE7"),
					Pos:  token.Pos{Offset: 233},
				},
				{
					Type: token.TokMap.Type("INVALID"),
					Lit:  []byte("\xE8"),
					Pos:  token.Pos{Offset: 234},
				},
				{
					Type: token.TokMap.Type("INVALID"),
					Lit:  []byte("\xE9"),
					Pos:  token.Pos{Offset: 235},
				},
				{
					Type: token.TokMap.Type("INVALID"),
					Lit:  []byte("\xEA"),
					Pos:  token.Pos{Offset: 236},
				},
				{
					Type: token.TokMap.Type("INVALID"),
					Lit:  []byte("\xEB"),
					Pos:  token.Pos{Offset: 237},
				},
				{
					Type: token.TokMap.Type("INVALID"),
					Lit:  []byte("\xEC"),
					Pos:  token.Pos{Offset: 238},
				},
				{
					Type: token.TokMap.Type("INVALID"),
					Lit:  []byte("\xED"),
					Pos:  token.Pos{Offset: 239},
				},
				{
					Type: token.TokMap.Type("INVALID"),
					Lit:  []byte("\xEE"),
					Pos:  token.Pos{Offset: 240},
				},
				{
					Type: token.TokMap.Type("INVALID"),
					Lit:  []byte("\xEF"),
					Pos:  token.Pos{Offset: 241},
				},
				{
					Type: token.TokMap.Type("INVALID"),
					Lit:  []byte("\xF0"),
					Pos:  token.Pos{Offset: 242},
				},
				{
					Type: token.TokMap.Type("INVALID"),
					Lit:  []byte("\xF1"),
					Pos:  token.Pos{Offset: 243},
				},
				{
					Type: token.TokMap.Type("INVALID"),
					Lit:  []byte("\xF2"),
					Pos:  token.Pos{Offset: 244},
				},
				{
					Type: token.TokMap.Type("INVALID"),
					Lit:  []byte("\xF3"),
					Pos:  token.Pos{Offset: 245},
				},
				{
					Type: token.TokMap.Type("INVALID"),
					Lit:  []byte("\xF4"),
					Pos:  token.Pos{Offset: 246},
				},
				{
					Type: token.TokMap.Type("INVALID"),
					Lit:  []byte("\xF5"),
					Pos:  token.Pos{Offset: 247},
				},
				{
					Type: token.TokMap.Type("INVALID"),
					Lit:  []byte("\xF6"),
					Pos:  token.Pos{Offset: 248},
				},
				{
					Type: token.TokMap.Type("INVALID"),
					Lit:  []byte("\xF7"),
					Pos:  token.Pos{Offset: 249},
				},
				{
					Type: token.TokMap.Type("INVALID"),
					Lit:  []byte("\xF8"),
					Pos:  token.Pos{Offset: 250},
				},
				{
					Type: token.TokMap.Type("INVALID"),
					Lit:  []byte("\xF9"),
					Pos:  token.Pos{Offset: 251},
				},
				{
					Type: token.TokMap.Type("INVALID"),
					Lit:  []byte("\xFA"),
					Pos:  token.Pos{Offset: 252},
				},
				{
					Type: token.TokMap.Type("INVALID"),
					Lit:  []byte("\xFB"),
					Pos:  token.Pos{Offset: 253},
				},
				{
					Type: token.TokMap.Type("INVALID"),
					Lit:  []byte("\xFC"),
					Pos:  token.Pos{Offset: 254},
				},
				{
					Type: token.TokMap.Type("INVALID"),
					Lit:  []byte("\xFD"),
					Pos:  token.Pos{Offset: 255},
				},
				{
					Type: token.TokMap.Type("INVALID"),
					Lit:  []byte("\xFE"),
					Pos:  token.Pos{Offset: 256},
				},
				{
					Type: token.TokMap.Type("INVALID"),
					Lit:  []byte("\xFF"),
					Pos:  token.Pos{Offset: 257},
				},
				{
					Type: token.TokMap.Type("+"),
					Lit:  []byte("+"),
					Pos:  token.Pos{Offset: 258},
				},
				{
					Type: token.TokMap.Type("INVALID"),
					Lit:  []byte("|"),
					Pos:  token.Pos{Offset: 259},
				},
				{
					Type: token.TokMap.Type("ident"),
					Lit:  []byte("int"),
					Pos:  token.Pos{Offset: 261},
				},
				{
					Type: token.TokMap.Type("ident"),
					Lit:  []byte("a"),
					Pos:  token.Pos{Offset: 265},
				},
				{
					Type: token.TokMap.Type(","),
					Lit:  []byte(","),
					Pos:  token.Pos{Offset: 267},
				},
				{
					Type: token.TokMap.Type("ident"),
					Lit:  []byte("b"),
					Pos:  token.Pos{Offset: 269},
				},
				{
					Type: token.TokMap.Type(","),
					Lit:  []byte(","),
					Pos:  token.Pos{Offset: 271},
				},
				{
					Type: token.TokMap.Type("ident"),
					Lit:  []byte("c"),
					Pos:  token.Pos{Offset: 273},
				},
				{
					Type: token.TokMap.Type(","),
					Lit:  []byte(","),
					Pos:  token.Pos{Offset: 275},
				},
				{
					Type: token.TokMap.Type("ident"),
					Lit:  []byte("d"),
					Pos:  token.Pos{Offset: 277},
				},
				{
					Type: token.TokMap.Type(";"),
					Lit:  []byte(";"),
					Pos:  token.Pos{Offset: 279},
				},
				{
					Type: token.TokMap.Type("ident"),
					Lit:  []byte("BEGIN"),
					Pos:  token.Pos{Offset: 281},
				},
				{
					Type: token.TokMap.Type("ident"),
					Lit:  []byte("a"),
					Pos:  token.Pos{Offset: 287},
				},
				{
					Type: token.TokMap.Type("="),
					Lit:  []byte("="),
					Pos:  token.Pos{Offset: 289},
				},
				{
					Type: token.TokMap.Type("("),
					Lit:  []byte("("),
					Pos:  token.Pos{Offset: 291},
				},
				{
					Type: token.TokMap.Type("int_lit"),
					Lit:  []byte("4711"),
					Pos:  token.Pos{Offset: 293},
				},
				{
					Type: token.TokMap.Type("+"),
					Lit:  []byte("+"),
					Pos:  token.Pos{Offset: 298},
				},
				{
					Type: token.TokMap.Type("int_lit"),
					Lit:  []byte("17"),
					Pos:  token.Pos{Offset: 300},
				},
				{
					Type: token.TokMap.Type("-"),
					Lit:  []byte("-"),
					Pos:  token.Pos{Offset: 303},
				},
				{
					Type: token.TokMap.Type("int_lit"),
					Lit:  []byte("2001"),
					Pos:  token.Pos{Offset: 305},
				},
				{
					Type: token.TokMap.Type(")"),
					Lit:  []byte(")"),
					Pos:  token.Pos{Offset: 310},
				},
				{
					Type: token.TokMap.Type("/"),
					Lit:  []byte("/"),
					Pos:  token.Pos{Offset: 312},
				},
				{
					Type: token.TokMap.Type("int_lit"),
					Lit:  []byte("100"),
					Pos:  token.Pos{Offset: 314},
				},
				{
					Type: token.TokMap.Type("-"),
					Lit:  []byte("-"),
					Pos:  token.Pos{Offset: 318},
				},
				{
					Type: token.TokMap.Type("int_lit"),
					Lit:  []byte("17"),
					Pos:  token.Pos{Offset: 320},
				},
				{
					Type: token.TokMap.Type(";"),
					Lit:  []byte(";"),
					Pos:  token.Pos{Offset: 323},
				},
				{
					Type: token.TokMap.Type("ident"),
					Lit:  []byte("b"),
					Pos:  token.Pos{Offset: 325},
				},
				{
					Type: token.TokMap.Type("="),
					Lit:  []byte("="),
					Pos:  token.Pos{Offset: 327},
				},
				{
					Type: token.TokMap.Type("ident"),
					Lit:  []byte("a"),
					Pos:  token.Pos{Offset: 329},
				},
				{
					Type: token.TokMap.Type("/"),
					Lit:  []byte("/"),
					Pos:  token.Pos{Offset: 331},
				},
				{
					Type: token.TokMap.Type("int_lit"),
					Lit:  []byte("5"),
					Pos:  token.Pos{Offset: 333},
				},
				{
					Type: token.TokMap.Type(";"),
					Lit:  []byte(";"),
					Pos:  token.Pos{Offset: 335},
				},
				{
					Type: token.TokMap.Type("ident"),
					Lit:  []byte("c"),
					Pos:  token.Pos{Offset: 337},
				},
				{
					Type: token.TokMap.Type("="),
					Lit:  []byte("="),
					Pos:  token.Pos{Offset: 339},
				},
				{
					Type: token.TokMap.Type("ident"),
					Lit:  []byte("a"),
					Pos:  token.Pos{Offset: 341},
				},
				{
					Type: token.TokMap.Type("*"),
					Lit:  []byte("*"),
					Pos:  token.Pos{Offset: 343},
				},
				{
					Type: token.TokMap.Type("ident"),
					Lit:  []byte("b"),
					Pos:  token.Pos{Offset: 345},
				},
				{
					Type: token.TokMap.Type("*"),
					Lit:  []byte("*"),
					Pos:  token.Pos{Offset: 347},
				},
				{
					Type: token.TokMap.Type("ident"),
					Lit:  []byte("b"),
					Pos:  token.Pos{Offset: 349},
				},
				{
					Type: token.TokMap.Type("*"),
					Lit:  []byte("*"),
					Pos:  token.Pos{Offset: 351},
				},
				{
					Type: token.TokMap.Type("ident"),
					Lit:  []byte("a"),
					Pos:  token.Pos{Offset: 353},
				},
				{
					Type: token.TokMap.Type("*"),
					Lit:  []byte("*"),
					Pos:  token.Pos{Offset: 355},
				},
				{
					Type: token.TokMap.Type("ident"),
					Lit:  []byte("b"),
					Pos:  token.Pos{Offset: 357},
				},
				{
					Type: token.TokMap.Type("*"),
					Lit:  []byte("*"),
					Pos:  token.Pos{Offset: 359},
				},
				{
					Type: token.TokMap.Type("ident"),
					Lit:  []byte("b"),
					Pos:  token.Pos{Offset: 361},
				},
				{
					Type: token.TokMap.Type(";"),
					Lit:  []byte(";"),
					Pos:  token.Pos{Offset: 363},
				},
				{
					Type: token.TokMap.Type("ident"),
					Lit:  []byte("d"),
					Pos:  token.Pos{Offset: 365},
				},
				{
					Type: token.TokMap.Type("="),
					Lit:  []byte("="),
					Pos:  token.Pos{Offset: 367},
				},
				{
					Type: token.TokMap.Type("ident"),
					Lit:  []byte("c"),
					Pos:  token.Pos{Offset: 369},
				},
				{
					Type: token.TokMap.Type("+"),
					Lit:  []byte("+"),
					Pos:  token.Pos{Offset: 371},
				},
				{
					Type: token.TokMap.Type("ident"),
					Lit:  []byte("c"),
					Pos:  token.Pos{Offset: 373},
				},
				{
					Type: token.TokMap.Type("/"),
					Lit:  []byte("/"),
					Pos:  token.Pos{Offset: 375},
				},
				{
					Type: token.TokMap.Type("("),
					Lit:  []byte("("),
					Pos:  token.Pos{Offset: 377},
				},
				{
					Type: token.TokMap.Type("int_lit"),
					Lit:  []byte("100"),
					Pos:  token.Pos{Offset: 379},
				},
				{
					Type: token.TokMap.Type("/"),
					Lit:  []byte("/"),
					Pos:  token.Pos{Offset: 383},
				},
				{
					Type: token.TokMap.Type("int_lit"),
					Lit:  []byte("2"),
					Pos:  token.Pos{Offset: 385},
				},
				{
					Type: token.TokMap.Type(")"),
					Lit:  []byte(")"),
					Pos:  token.Pos{Offset: 387},
				},
				{
					Type: token.TokMap.Type(";"),
					Lit:  []byte(";"),
					Pos:  token.Pos{Offset: 389},
				},
				{
					Type: token.TokMap.Type("if"),
					Lit:  []byte("if"),
					Pos:  token.Pos{Offset: 391},
				},
				{
					Type: token.TokMap.Type("("),
					Lit:  []byte("("),
					Pos:  token.Pos{Offset: 394},
				},
				{
					Type: token.TokMap.Type("ident"),
					Lit:  []byte("a"),
					Pos:  token.Pos{Offset: 396},
				},
				{
					Type: token.TokMap.Type(">"),
					Lit:  []byte(">"),
					Pos:  token.Pos{Offset: 398},
				},
				{
					Type: token.TokMap.Type("int_lit"),
					Lit:  []byte("3"),
					Pos:  token.Pos{Offset: 400},
				},
				{
					Type: token.TokMap.Type("*"),
					Lit:  []byte("*"),
					Pos:  token.Pos{Offset: 402},
				},
				{
					Type: token.TokMap.Type("ident"),
					Lit:  []byte("b"),
					Pos:  token.Pos{Offset: 404},
				},
				{
					Type: token.TokMap.Type("ident"),
					Lit:  []byte("or"),
					Pos:  token.Pos{Offset: 406},
				},
				{
					Type: token.TokMap.Type("ident"),
					Lit:  []byte("not"),
					Pos:  token.Pos{Offset: 409},
				},
				{
					Type: token.TokMap.Type("ident"),
					Lit:  []byte("a"),
					Pos:  token.Pos{Offset: 413},
				},
				{
					Type: token.TokMap.Type("=="),
					Lit:  []byte("=="),
					Pos:  token.Pos{Offset: 415},
				},
				{
					Type: token.TokMap.Type("ident"),
					Lit:  []byte("c"),
					Pos:  token.Pos{Offset: 418},
				},
				{
					Type: token.TokMap.Type("/"),
					Lit:  []byte("/"),
					Pos:  token.Pos{Offset: 420},
				},
				{
					Type: token.TokMap.Type("int_lit"),
					Lit:  []byte("2"),
					Pos:  token.Pos{Offset: 422},
				},
				{
					Type: token.TokMap.Type("ident"),
					Lit:  []byte("and"),
					Pos:  token.Pos{Offset: 424},
				},
				{
					Type: token.TokMap.Type("ident"),
					Lit:  []byte("c"),
					Pos:  token.Pos{Offset: 428},
				},
				{
					Type: token.TokMap.Type(">"),
					Lit:  []byte(">"),
					Pos:  token.Pos{Offset: 430},
				},
				{
					Type: token.TokMap.Type("("),
					Lit:  []byte("("),
					Pos:  token.Pos{Offset: 432},
				},
				{
					Type: token.TokMap.Type("ident"),
					Lit:  []byte("d"),
					Pos:  token.Pos{Offset: 434},
				},
				{
					Type: token.TokMap.Type("-"),
					Lit:  []byte("-"),
					Pos:  token.Pos{Offset: 436},
				},
				{
					Type: token.TokMap.Type("int_lit"),
					Lit:  []byte("3"),
					Pos:  token.Pos{Offset: 438},
				},
				{
					Type: token.TokMap.Type(")"),
					Lit:  []byte(")"),
					Pos:  token.Pos{Offset: 440},
				},
				{
					Type: token.TokMap.Type("*"),
					Lit:  []byte("*"),
					Pos:  token.Pos{Offset: 442},
				},
				{
					Type: token.TokMap.Type("("),
					Lit:  []byte("("),
					Pos:  token.Pos{Offset: 444},
				},
				{
					Type: token.TokMap.Type("ident"),
					Lit:  []byte("a"),
					Pos:  token.Pos{Offset: 446},
				},
				{
					Type: token.TokMap.Type("-"),
					Lit:  []byte("-"),
					Pos:  token.Pos{Offset: 448},
				},
				{
					Type: token.TokMap.Type("ident"),
					Lit:  []byte("b"),
					Pos:  token.Pos{Offset: 450},
				},
				{
					Type: token.TokMap.Type(")"),
					Lit:  []byte(")"),
					Pos:  token.Pos{Offset: 452},
				},
				{
					Type: token.TokMap.Type(")"),
					Lit:  []byte(")"),
					Pos:  token.Pos{Offset: 454},
				},
				{
					Type: token.TokMap.Type("ident"),
					Lit:  []byte("and"),
					Pos:  token.Pos{Offset: 456},
				},
				{
					Type: token.TokMap.Type("ident"),
					Lit:  []byte("a"),
					Pos:  token.Pos{Offset: 460},
				},
				{
					Type: token.TokMap.Type("<"),
					Lit:  []byte("<"),
					Pos:  token.Pos{Offset: 462},
				},
				{
					Type: token.TokMap.Type("int_lit"),
					Lit:  []byte("0"),
					Pos:  token.Pos{Offset: 464},
				},
				{
					Type: token.TokMap.Type("ident"),
					Lit:  []byte("then"),
					Pos:  token.Pos{Offset: 466},
				},
				{
					Type: token.TokMap.Type("ident"),
					Lit:  []byte("begin"),
					Pos:  token.Pos{Offset: 471},
				},
				{
					Type: token.TokMap.Type("ident"),
					Lit:  []byte("a"),
					Pos:  token.Pos{Offset: 477},
				},
				{
					Type: token.TokMap.Type("="),
					Lit:  []byte("="),
					Pos:  token.Pos{Offset: 479},
				},
				{
					Type: token.TokMap.Type("int_lit"),
					Lit:  []byte("1"),
					Pos:  token.Pos{Offset: 481},
				},
				{
					Type: token.TokMap.Type(";"),
					Lit:  []byte(";"),
					Pos:  token.Pos{Offset: 483},
				},
				{
					Type: token.TokMap.Type("ident"),
					Lit:  []byte("read"),
					Pos:  token.Pos{Offset: 485},
				},
				{
					Type: token.TokMap.Type("ident"),
					Lit:  []byte("b"),
					Pos:  token.Pos{Offset: 490},
				},
				{
					Type: token.TokMap.Type(";"),
					Lit:  []byte(";"),
					Pos:  token.Pos{Offset: 492},
				},
				{
					Type: token.TokMap.Type("ident"),
					Lit:  []byte("end"),
					Pos:  token.Pos{Offset: 494},
				},
				{
					Type: token.TokMap.Type("else"),
					Lit:  []byte("else"),
					Pos:  token.Pos{Offset: 498},
				},
				{
					Type: token.TokMap.Type("ident"),
					Lit:  []byte("begin"),
					Pos:  token.Pos{Offset: 503},
				},
				{
					Type: token.TokMap.Type("ident"),
					Lit:  []byte("a"),
					Pos:  token.Pos{Offset: 509},
				},
				{
					Type: token.TokMap.Type("="),
					Lit:  []byte("="),
					Pos:  token.Pos{Offset: 511},
				},
				{
					Type: token.TokMap.Type("int_lit"),
					Lit:  []byte("0"),
					Pos:  token.Pos{Offset: 513},
				},
				{
					Type: token.TokMap.Type(";"),
					Lit:  []byte(";"),
					Pos:  token.Pos{Offset: 515},
				},
				{
					Type: token.TokMap.Type("ident"),
					Lit:  []byte("read"),
					Pos:  token.Pos{Offset: 517},
				},
				{
					Type: token.TokMap.Type("ident"),
					Lit:  []byte("c"),
					Pos:  token.Pos{Offset: 522},
				},
				{
					Type: token.TokMap.Type(";"),
					Lit:  []byte(";"),
					Pos:  token.Pos{Offset: 524},
				},
				{
					Type: token.TokMap.Type("ident"),
					Lit:  []byte("end"),
					Pos:  token.Pos{Offset: 526},
				},
				{
					Type: token.TokMap.Type("ident"),
					Lit:  []byte("END"),
					Pos:  token.Pos{Offset: 530},
				},
				{
					Type: token.TokMap.Type("$"),
					Lit:  []byte(""),
					Pos:  token.Pos{Offset: 533},
				},
			},
		},

		{
			path: "../../testdata/quiet/lexer/l01.c",
			toks: []*token.Token{
				{
					Type: token.TokMap.Type("ident"),
					Lit:  []byte("int"),
					Pos:  token.Pos{Offset: 0},
				},
				{
					Type: token.TokMap.Type("ident"),
					Lit:  []byte("main"),
					Pos:  token.Pos{Offset: 4},
				},
				{
					Type: token.TokMap.Type("("),
					Lit:  []byte("("),
					Pos:  token.Pos{Offset: 9},
				},
				{
					Type: token.TokMap.Type("ident"),
					Lit:  []byte("void"),
					Pos:  token.Pos{Offset: 10},
				},
				{
					Type: token.TokMap.Type(")"),
					Lit:  []byte(")"),
					Pos:  token.Pos{Offset: 14},
				},
				{
					Type: token.TokMap.Type("{"),
					Lit:  []byte("{"),
					Pos:  token.Pos{Offset: 16},
				},
				{
					Type: token.TokMap.Type(";"),
					Lit:  []byte(";"),
					Pos:  token.Pos{Offset: 18},
				},
				{
					Type: token.TokMap.Type("}"),
					Lit:  []byte("}"),
					Pos:  token.Pos{Offset: 20},
				},
				{
					Type: token.TokMap.Type("$"),
					Lit:  []byte(""),
					Pos:  token.Pos{Offset: 23},
				},
			},
		},

		{
			path: "../../testdata/quiet/lexer/l02.c",
			toks: []*token.Token{
				{
					Type: token.TokMap.Type("ident"),
					Lit:  []byte("int"),
					Pos:  token.Pos{Offset: 0},
				},
				{
					Type: token.TokMap.Type("ident"),
					Lit:  []byte("foo"),
					Pos:  token.Pos{Offset: 4},
				},
				{
					Type: token.TokMap.Type(";"),
					Lit:  []byte(";"),
					Pos:  token.Pos{Offset: 7},
				},
				{
					Type: token.TokMap.Type("ident"),
					Lit:  []byte("int"),
					Pos:  token.Pos{Offset: 10},
				},
				{
					Type: token.TokMap.Type("ident"),
					Lit:  []byte("BarBara"),
					Pos:  token.Pos{Offset: 14},
				},
				{
					Type: token.TokMap.Type(";"),
					Lit:  []byte(";"),
					Pos:  token.Pos{Offset: 21},
				},
				{
					Type: token.TokMap.Type("ident"),
					Lit:  []byte("int"),
					Pos:  token.Pos{Offset: 24},
				},
				{
					Type: token.TokMap.Type("ident"),
					Lit:  []byte("bar_bara"),
					Pos:  token.Pos{Offset: 28},
				},
				{
					Type: token.TokMap.Type(";"),
					Lit:  []byte(";"),
					Pos:  token.Pos{Offset: 36},
				},
				{
					Type: token.TokMap.Type("ident"),
					Lit:  []byte("int"),
					Pos:  token.Pos{Offset: 39},
				},
				{
					Type: token.TokMap.Type("ident"),
					Lit:  []byte("bar4711"),
					Pos:  token.Pos{Offset: 43},
				},
				{
					Type: token.TokMap.Type(";"),
					Lit:  []byte(";"),
					Pos:  token.Pos{Offset: 50},
				},
				{
					Type: token.TokMap.Type("ident"),
					Lit:  []byte("int"),
					Pos:  token.Pos{Offset: 53},
				},
				{
					Type: token.TokMap.Type("ident"),
					Lit:  []byte("b4rb4r4"),
					Pos:  token.Pos{Offset: 57},
				},
				{
					Type: token.TokMap.Type(";"),
					Lit:  []byte(";"),
					Pos:  token.Pos{Offset: 64},
				},
				{
					Type: token.TokMap.Type("ident"),
					Lit:  []byte("int"),
					Pos:  token.Pos{Offset: 67},
				},
				{
					Type: token.TokMap.Type("ident"),
					Lit:  []byte("abcdefghijklmnopqrstuvwxyz_ABCDEFGHIJKLMNOPQRSTUVWXYZ_0123456789"),
					Pos:  token.Pos{Offset: 71},
				},
				{
					Type: token.TokMap.Type(";"),
					Lit:  []byte(";"),
					Pos:  token.Pos{Offset: 135},
				},
				{
					Type: token.TokMap.Type("ident"),
					Lit:  []byte("int"),
					Pos:  token.Pos{Offset: 138},
				},
				{
					Type: token.TokMap.Type("ident"),
					Lit:  []byte("main"),
					Pos:  token.Pos{Offset: 142},
				},
				{
					Type: token.TokMap.Type("("),
					Lit:  []byte("("),
					Pos:  token.Pos{Offset: 146},
				},
				{
					Type: token.TokMap.Type("ident"),
					Lit:  []byte("void"),
					Pos:  token.Pos{Offset: 147},
				},
				{
					Type: token.TokMap.Type(")"),
					Lit:  []byte(")"),
					Pos:  token.Pos{Offset: 151},
				},
				{
					Type: token.TokMap.Type("{"),
					Lit:  []byte("{"),
					Pos:  token.Pos{Offset: 153},
				},
				{
					Type: token.TokMap.Type(";"),
					Lit:  []byte(";"),
					Pos:  token.Pos{Offset: 155},
				},
				{
					Type: token.TokMap.Type("}"),
					Lit:  []byte("}"),
					Pos:  token.Pos{Offset: 157},
				},
				{
					Type: token.TokMap.Type("$"),
					Lit:  []byte(""),
					Pos:  token.Pos{Offset: 159},
				},
			},
		},

		{
			path: "../../testdata/quiet/lexer/l03.c",
			toks: []*token.Token{
				{
					Type: token.TokMap.Type("ident"),
					Lit:  []byte("int"),
					Pos:  token.Pos{Offset: 0},
				},
				{
					Type: token.TokMap.Type("ident"),
					Lit:  []byte("main"),
					Pos:  token.Pos{Offset: 4},
				},
				{
					Type: token.TokMap.Type("("),
					Lit:  []byte("("),
					Pos:  token.Pos{Offset: 8},
				},
				{
					Type: token.TokMap.Type("ident"),
					Lit:  []byte("void"),
					Pos:  token.Pos{Offset: 9},
				},
				{
					Type: token.TokMap.Type(")"),
					Lit:  []byte(")"),
					Pos:  token.Pos{Offset: 13},
				},
				{
					Type: token.TokMap.Type("{"),
					Lit:  []byte("{"),
					Pos:  token.Pos{Offset: 15},
				},
				{
					Type: token.TokMap.Type("ident"),
					Lit:  []byte("int"),
					Pos:  token.Pos{Offset: 19},
				},
				{
					Type: token.TokMap.Type("ident"),
					Lit:  []byte("i"),
					Pos:  token.Pos{Offset: 23},
				},
				{
					Type: token.TokMap.Type(";"),
					Lit:  []byte(";"),
					Pos:  token.Pos{Offset: 24},
				},
				{
					Type: token.TokMap.Type("ident"),
					Lit:  []byte("i"),
					Pos:  token.Pos{Offset: 28},
				},
				{
					Type: token.TokMap.Type("="),
					Lit:  []byte("="),
					Pos:  token.Pos{Offset: 30},
				},
				{
					Type: token.TokMap.Type("int_lit"),
					Lit:  []byte("123456789"),
					Pos:  token.Pos{Offset: 32},
				},
				{
					Type: token.TokMap.Type(";"),
					Lit:  []byte(";"),
					Pos:  token.Pos{Offset: 41},
				},
				{
					Type: token.TokMap.Type("comment"),
					Lit:  []byte("// Was:   i = 1234567890;\n"),
					Pos:  token.Pos{Offset: 43},
				},
				{
					Type: token.TokMap.Type("ident"),
					Lit:  []byte("i"),
					Pos:  token.Pos{Offset: 71},
				},
				{
					Type: token.TokMap.Type("="),
					Lit:  []byte("="),
					Pos:  token.Pos{Offset: 73},
				},
				{
					Type: token.TokMap.Type("int_lit"), // TODO: Figure out if char_lit and int_lit should be separate tokens.
					Lit:  []byte("'0'"),
					Pos:  token.Pos{Offset: 75},
				},
				{
					Type: token.TokMap.Type(";"),
					Lit:  []byte(";"),
					Pos:  token.Pos{Offset: 78},
				},
				{
					Type: token.TokMap.Type("ident"),
					Lit:  []byte("i"),
					Pos:  token.Pos{Offset: 82},
				},
				{
					Type: token.TokMap.Type("="),
					Lit:  []byte("="),
					Pos:  token.Pos{Offset: 84},
				},
				{
					Type: token.TokMap.Type("int_lit"), // TODO: Figure out if char_lit and int_lit should be separate tokens.
					Lit:  []byte("'a'"),
					Pos:  token.Pos{Offset: 86},
				},
				{
					Type: token.TokMap.Type(";"),
					Lit:  []byte(";"),
					Pos:  token.Pos{Offset: 89},
				},
				{
					Type: token.TokMap.Type("ident"),
					Lit:  []byte("i"),
					Pos:  token.Pos{Offset: 93},
				},
				{
					Type: token.TokMap.Type("="),
					Lit:  []byte("="),
					Pos:  token.Pos{Offset: 95},
				},
				{
					Type: token.TokMap.Type("int_lit"), // TODO: Figure out if char_lit and int_lit should be separate tokens.
					Lit:  []byte("' '"),
					Pos:  token.Pos{Offset: 97},
				},
				{
					Type: token.TokMap.Type(";"),
					Lit:  []byte(";"),
					Pos:  token.Pos{Offset: 100},
				},
				{
					Type: token.TokMap.Type("ident"),
					Lit:  []byte("i"),
					Pos:  token.Pos{Offset: 104},
				},
				{
					Type: token.TokMap.Type("="),
					Lit:  []byte("="),
					Pos:  token.Pos{Offset: 106},
				},
				{
					Type: token.TokMap.Type("int_lit"), // TODO: Figure out if char_lit and int_lit should be separate tokens.
					Lit:  []byte(`'\n'`),
					Pos:  token.Pos{Offset: 108},
				},
				{
					Type: token.TokMap.Type(";"),
					Lit:  []byte(";"),
					Pos:  token.Pos{Offset: 112},
				},
				{
					Type: token.TokMap.Type("}"),
					Lit:  []byte("}"),
					Pos:  token.Pos{Offset: 114},
				},
				{
					Type: token.TokMap.Type("$"),
					Lit:  []byte(""),
					Pos:  token.Pos{Offset: 117},
				},
			},
		},

		{
			path: "../../testdata/quiet/lexer/l04.c",
			toks: []*token.Token{
				{
					Type: token.TokMap.Type("ident"),
					Lit:  []byte("int"),
					Pos:  token.Pos{Offset: 0},
				},
				{
					Type: token.TokMap.Type("ident"),
					Lit:  []byte("main"),
					Pos:  token.Pos{Offset: 4},
				},
				{
					Type: token.TokMap.Type("("),
					Lit:  []byte("("),
					Pos:  token.Pos{Offset: 8},
				},
				{
					Type: token.TokMap.Type("ident"),
					Lit:  []byte("void"),
					Pos:  token.Pos{Offset: 9},
				},
				{
					Type: token.TokMap.Type(")"),
					Lit:  []byte(")"),
					Pos:  token.Pos{Offset: 13},
				},
				{
					Type: token.TokMap.Type("{"),
					Lit:  []byte("{"),
					Pos:  token.Pos{Offset: 15},
				},
				{
					Type: token.TokMap.Type("ident"),
					Lit:  []byte("int"),
					Pos:  token.Pos{Offset: 19},
				},
				{
					Type: token.TokMap.Type("ident"),
					Lit:  []byte("i"),
					Pos:  token.Pos{Offset: 23},
				},
				{
					Type: token.TokMap.Type(";"),
					Lit:  []byte(";"),
					Pos:  token.Pos{Offset: 24},
				},
				{
					Type: token.TokMap.Type("ident"),
					Lit:  []byte("char"),
					Pos:  token.Pos{Offset: 28},
				},
				{
					Type: token.TokMap.Type("ident"),
					Lit:  []byte("j"),
					Pos:  token.Pos{Offset: 33},
				},
				{
					Type: token.TokMap.Type(";"),
					Lit:  []byte(";"),
					Pos:  token.Pos{Offset: 34},
				},
				{
					Type: token.TokMap.Type("if"),
					Lit:  []byte("if"),
					Pos:  token.Pos{Offset: 38},
				},
				{
					Type: token.TokMap.Type("("),
					Lit:  []byte("("),
					Pos:  token.Pos{Offset: 41},
				},
				{
					Type: token.TokMap.Type("int_lit"),
					Lit:  []byte("1"),
					Pos:  token.Pos{Offset: 42},
				},
				{
					Type: token.TokMap.Type("=="),
					Lit:  []byte("=="),
					Pos:  token.Pos{Offset: 43},
				},
				{
					Type: token.TokMap.Type("int_lit"),
					Lit:  []byte("0"),
					Pos:  token.Pos{Offset: 45},
				},
				{
					Type: token.TokMap.Type(")"),
					Lit:  []byte(")"),
					Pos:  token.Pos{Offset: 46},
				},
				{
					Type: token.TokMap.Type("ident"),
					Lit:  []byte("i"),
					Pos:  token.Pos{Offset: 48},
				},
				{
					Type: token.TokMap.Type("="),
					Lit:  []byte("="),
					Pos:  token.Pos{Offset: 50},
				},
				{
					Type: token.TokMap.Type("int_lit"),
					Lit:  []byte("0"),
					Pos:  token.Pos{Offset: 52},
				},
				{
					Type: token.TokMap.Type(";"),
					Lit:  []byte(";"),
					Pos:  token.Pos{Offset: 53},
				},
				{
					Type: token.TokMap.Type("else"),
					Lit:  []byte("else"),
					Pos:  token.Pos{Offset: 58},
				},
				{
					Type: token.TokMap.Type("ident"),
					Lit:  []byte("i"),
					Pos:  token.Pos{Offset: 63},
				},
				{
					Type: token.TokMap.Type("="),
					Lit:  []byte("="),
					Pos:  token.Pos{Offset: 65},
				},
				{
					Type: token.TokMap.Type("int_lit"),
					Lit:  []byte("1"),
					Pos:  token.Pos{Offset: 67},
				},
				{
					Type: token.TokMap.Type(";"),
					Lit:  []byte(";"),
					Pos:  token.Pos{Offset: 68},
				},
				{
					Type: token.TokMap.Type("while"),
					Lit:  []byte("while"),
					Pos:  token.Pos{Offset: 72},
				},
				{
					Type: token.TokMap.Type("("),
					Lit:  []byte("("),
					Pos:  token.Pos{Offset: 78},
				},
				{
					Type: token.TokMap.Type("int_lit"),
					Lit:  []byte("1"),
					Pos:  token.Pos{Offset: 79},
				},
				{
					Type: token.TokMap.Type("=="),
					Lit:  []byte("=="),
					Pos:  token.Pos{Offset: 80},
				},
				{
					Type: token.TokMap.Type("int_lit"),
					Lit:  []byte("0"),
					Pos:  token.Pos{Offset: 82},
				},
				{
					Type: token.TokMap.Type(")"),
					Lit:  []byte(")"),
					Pos:  token.Pos{Offset: 83},
				},
				{
					Type: token.TokMap.Type("ident"),
					Lit:  []byte("i"),
					Pos:  token.Pos{Offset: 85},
				},
				{
					Type: token.TokMap.Type("="),
					Lit:  []byte("="),
					Pos:  token.Pos{Offset: 87},
				},
				{
					Type: token.TokMap.Type("int_lit"),
					Lit:  []byte("0"),
					Pos:  token.Pos{Offset: 89},
				},
				{
					Type: token.TokMap.Type(";"),
					Lit:  []byte(";"),
					Pos:  token.Pos{Offset: 90},
				},
				{
					Type: token.TokMap.Type("return"),
					Lit:  []byte("return"),
					Pos:  token.Pos{Offset: 94},
				},
				{
					Type: token.TokMap.Type("int_lit"),
					Lit:  []byte("42"),
					Pos:  token.Pos{Offset: 101},
				},
				{
					Type: token.TokMap.Type(";"),
					Lit:  []byte(";"),
					Pos:  token.Pos{Offset: 103},
				},
				{
					Type: token.TokMap.Type("}"),
					Lit:  []byte("}"),
					Pos:  token.Pos{Offset: 105},
				},
				{
					Type: token.TokMap.Type("$"),
					Lit:  []byte(""),
					Pos:  token.Pos{Offset: 107},
				},
			},
		},

		{
			path: "../../testdata/quiet/lexer/l05.c",
			toks: []*token.Token{
				{
					Type: token.TokMap.Type("ident"),
					Lit:  []byte("int"),
					Pos:  token.Pos{Offset: 0},
				},
				{
					Type: token.TokMap.Type("ident"),
					Lit:  []byte("main"),
					Pos:  token.Pos{Offset: 4},
				},
				{
					Type: token.TokMap.Type("("),
					Lit:  []byte("("),
					Pos:  token.Pos{Offset: 8},
				},
				{
					Type: token.TokMap.Type("ident"),
					Lit:  []byte("void"),
					Pos:  token.Pos{Offset: 9},
				},
				{
					Type: token.TokMap.Type(")"),
					Lit:  []byte(")"),
					Pos:  token.Pos{Offset: 13},
				},
				{
					Type: token.TokMap.Type("{"),
					Lit:  []byte("{"),
					Pos:  token.Pos{Offset: 15},
				},
				{
					Type: token.TokMap.Type("ident"),
					Lit:  []byte("int"),
					Pos:  token.Pos{Offset: 19},
				},
				{
					Type: token.TokMap.Type("ident"),
					Lit:  []byte("i"),
					Pos:  token.Pos{Offset: 23},
				},
				{
					Type: token.TokMap.Type(";"),
					Lit:  []byte(";"),
					Pos:  token.Pos{Offset: 24},
				},
				{
					Type: token.TokMap.Type("int_lit"),
					Lit:  []byte("1"),
					Pos:  token.Pos{Offset: 28},
				},
				{
					Type: token.TokMap.Type("!="),
					Lit:  []byte("!="),
					Pos:  token.Pos{Offset: 29},
				},
				{
					Type: token.TokMap.Type("!"),
					Lit:  []byte("!"),
					Pos:  token.Pos{Offset: 31},
				},
				{
					Type: token.TokMap.Type("int_lit"),
					Lit:  []byte("3"),
					Pos:  token.Pos{Offset: 32},
				},
				{
					Type: token.TokMap.Type(";"),
					Lit:  []byte(";"),
					Pos:  token.Pos{Offset: 33},
				},
				{
					Type: token.TokMap.Type("int_lit"),
					Lit:  []byte("4"),
					Pos:  token.Pos{Offset: 37},
				},
				{
					Type: token.TokMap.Type("&&"),
					Lit:  []byte("&&"),
					Pos:  token.Pos{Offset: 38},
				},
				{
					Type: token.TokMap.Type("("),
					Lit:  []byte("("),
					Pos:  token.Pos{Offset: 40},
				},
				{
					Type: token.TokMap.Type("int_lit"),
					Lit:  []byte("6"),
					Pos:  token.Pos{Offset: 41},
				},
				{
					Type: token.TokMap.Type(")"),
					Lit:  []byte(")"),
					Pos:  token.Pos{Offset: 42},
				},
				{
					Type: token.TokMap.Type(";"),
					Lit:  []byte(";"),
					Pos:  token.Pos{Offset: 43},
				},
				{
					Type: token.TokMap.Type("int_lit"),
					Lit:  []byte("7"),
					Pos:  token.Pos{Offset: 47},
				},
				{
					Type: token.TokMap.Type("*"),
					Lit:  []byte("*"),
					Pos:  token.Pos{Offset: 48},
				},
				{
					Type: token.TokMap.Type("int_lit"),
					Lit:  []byte("8"),
					Pos:  token.Pos{Offset: 50},
				},
				{
					Type: token.TokMap.Type("+"),
					Lit:  []byte("+"),
					Pos:  token.Pos{Offset: 51},
				},
				{
					Type: token.TokMap.Type("int_lit"),
					Lit:  []byte("10"),
					Pos:  token.Pos{Offset: 52},
				},
				{
					Type: token.TokMap.Type(";"),
					Lit:  []byte(";"),
					Pos:  token.Pos{Offset: 54},
				},
				{
					Type: token.TokMap.Type("("),
					Lit:  []byte("("),
					Pos:  token.Pos{Offset: 58},
				},
				{
					Type: token.TokMap.Type("int_lit"),
					Lit:  []byte("11"),
					Pos:  token.Pos{Offset: 59},
				},
				{
					Type: token.TokMap.Type("-"),
					Lit:  []byte("-"),
					Pos:  token.Pos{Offset: 61},
				},
				{
					Type: token.TokMap.Type("int_lit"),
					Lit:  []byte("12"),
					Pos:  token.Pos{Offset: 62},
				},
				{
					Type: token.TokMap.Type(")"),
					Lit:  []byte(")"),
					Pos:  token.Pos{Offset: 64},
				},
				{
					Type: token.TokMap.Type("+"),
					Lit:  []byte("+"),
					Pos:  token.Pos{Offset: 65},
				},
				{
					Type: token.TokMap.Type("("),
					Lit:  []byte("("),
					Pos:  token.Pos{Offset: 66},
				},
				{
					Type: token.TokMap.Type("int_lit"),
					Lit:  []byte("12"),
					Pos:  token.Pos{Offset: 67},
				},
				{
					Type: token.TokMap.Type("/"),
					Lit:  []byte("/"),
					Pos:  token.Pos{Offset: 69},
				},
				{
					Type: token.TokMap.Type("int_lit"),
					Lit:  []byte("16"),
					Pos:  token.Pos{Offset: 70},
				},
				{
					Type: token.TokMap.Type(")"),
					Lit:  []byte(")"),
					Pos:  token.Pos{Offset: 72},
				},
				{
					Type: token.TokMap.Type(";"),
					Lit:  []byte(";"),
					Pos:  token.Pos{Offset: 73},
				},
				{
					Type: token.TokMap.Type("int_lit"),
					Lit:  []byte("17"),
					Pos:  token.Pos{Offset: 77},
				},
				{
					Type: token.TokMap.Type("<="),
					Lit:  []byte("<="),
					Pos:  token.Pos{Offset: 79},
				},
				{
					Type: token.TokMap.Type("int_lit"),
					Lit:  []byte("18"),
					Pos:  token.Pos{Offset: 81},
				},
				{
					Type: token.TokMap.Type("<"),
					Lit:  []byte("<"),
					Pos:  token.Pos{Offset: 84},
				},
				{
					Type: token.TokMap.Type("-"),
					Lit:  []byte("-"),
					Pos:  token.Pos{Offset: 85},
				},
				{
					Type: token.TokMap.Type("int_lit"),
					Lit:  []byte("20"),
					Pos:  token.Pos{Offset: 86},
				},
				{
					Type: token.TokMap.Type(";"),
					Lit:  []byte(";"),
					Pos:  token.Pos{Offset: 88},
				},
				{
					Type: token.TokMap.Type("ident"),
					Lit:  []byte("i"),
					Pos:  token.Pos{Offset: 92},
				},
				{
					Type: token.TokMap.Type("="),
					Lit:  []byte("="),
					Pos:  token.Pos{Offset: 93},
				},
				{
					Type: token.TokMap.Type("int_lit"),
					Lit:  []byte("21"),
					Pos:  token.Pos{Offset: 94},
				},
				{
					Type: token.TokMap.Type("=="),
					Lit:  []byte("=="),
					Pos:  token.Pos{Offset: 96},
				},
				{
					Type: token.TokMap.Type("int_lit"),
					Lit:  []byte("22"),
					Pos:  token.Pos{Offset: 98},
				},
				{
					Type: token.TokMap.Type(";"),
					Lit:  []byte(";"),
					Pos:  token.Pos{Offset: 100},
				},
				{
					Type: token.TokMap.Type("int_lit"),
					Lit:  []byte("25"),
					Pos:  token.Pos{Offset: 104},
				},
				{
					Type: token.TokMap.Type(">="),
					Lit:  []byte(">="),
					Pos:  token.Pos{Offset: 107},
				},
				{
					Type: token.TokMap.Type("int_lit"),
					Lit:  []byte("27"),
					Pos:  token.Pos{Offset: 109},
				},
				{
					Type: token.TokMap.Type(">"),
					Lit:  []byte(">"),
					Pos:  token.Pos{Offset: 111},
				},
				{
					Type: token.TokMap.Type("int_lit"),
					Lit:  []byte("28"),
					Pos:  token.Pos{Offset: 112},
				},
				{
					Type: token.TokMap.Type(";"),
					Lit:  []byte(";"),
					Pos:  token.Pos{Offset: 114},
				},
				{
					Type: token.TokMap.Type("}"),
					Lit:  []byte("}"),
					Pos:  token.Pos{Offset: 116},
				},
				{
					Type: token.TokMap.Type("$"),
					Lit:  []byte(""),
					Pos:  token.Pos{Offset: 118},
				},
			},
		},

		{
			path: "../../testdata/quiet/lexer/l06.c",
			toks: []*token.Token{
				{
					Type: token.TokMap.Type("comment"),
					Lit:  []byte("/* This file contains examples of various types of white space and comments. */"),
					Pos:  token.Pos{Offset: 0},
				},
				{
					Type: token.TokMap.Type("comment"),
					Lit:  []byte("// A blank (32)\n"),
					Pos:  token.Pos{Offset: 80},
				},
				{
					Type: token.TokMap.Type("comment"),
					Lit:  []byte("// A new-line (10)\n"),
					Pos:  token.Pos{Offset: 98},
				},
				{
					Type: token.TokMap.Type("comment"),
					Lit:  []byte("// a carriage-return (13)\n"),
					Pos:  token.Pos{Offset: 118},
				},
				{
					Type: token.TokMap.Type("comment"),
					Lit:  []byte("// form-feed (12)\n"),
					Pos:  token.Pos{Offset: 146},
				},
				{
					Type: token.TokMap.Type("comment"),
					Lit:  []byte("// tab(9)\n"),
					Pos:  token.Pos{Offset: 166},
				},
				{
					Type: token.TokMap.Type("comment"),
					Lit:  []byte("/* In uC, each file must contain at least one declaration */"),
					Pos:  token.Pos{Offset: 178},
				},
				{
					Type: token.TokMap.Type("comment"),
					Lit:  []byte("/* a comment... */"),
					Pos:  token.Pos{Offset: 240},
				},
				{
					Type: token.TokMap.Type("ident"),
					Lit:  []byte("int"),
					Pos:  token.Pos{Offset: 259},
				},
				{
					Type: token.TokMap.Type("ident"),
					Lit:  []byte("main"),
					Pos:  token.Pos{Offset: 263},
				},
				{
					Type: token.TokMap.Type("("),
					Lit:  []byte("("),
					Pos:  token.Pos{Offset: 267},
				},
				{
					Type: token.TokMap.Type("comment"),
					Lit:  []byte("/* ...and another */"),
					Pos:  token.Pos{Offset: 269},
				},
				{
					Type: token.TokMap.Type("ident"),
					Lit:  []byte("void"),
					Pos:  token.Pos{Offset: 297},
				},
				{
					Type: token.TokMap.Type(")"),
					Lit:  []byte(")"),
					Pos:  token.Pos{Offset: 302},
				},
				{
					Type: token.TokMap.Type("{"),
					Lit:  []byte("{"),
					Pos:  token.Pos{Offset: 304},
				},
				{
					Type: token.TokMap.Type(";"),
					Lit:  []byte(";"),
					Pos:  token.Pos{Offset: 306},
				},
				{
					Type: token.TokMap.Type("}"),
					Lit:  []byte("}"),
					Pos:  token.Pos{Offset: 308},
				},
				{
					Type: token.TokMap.Type("$"),
					Lit:  []byte(""),
					Pos:  token.Pos{Offset: 313},
				},
			},
		},
	}

	for _, g := range golden {
		log.Println("path:", g.path)
		s, err := lexer.NewLexerFile(g.path)
		if err != nil {
			t.Error(err)
			continue
		}
		for j := 0; ; j++ {
			got := s.Scan()
			if j >= len(g.toks) {
				t.Errorf("%s: invalid number of tokens; expected %d tokens, got > %d", g.path, len(g.toks), j)
				break
			}
			if want := g.toks[j]; !tokenEqual(got, want) {
				t.Errorf("%s: token %d mismatch; expected %#v, got %#v", g.path, j, want, got)
				fmt.Printf("gocc token at %d: %q\n", got.Pos.Offset, string(got.Lit)) // TODO: Remove.
			}
			if got.Type == token.EOF {
				if j != len(g.toks)-1 {
					t.Errorf("%s: invalid number of tokens; expected %d tokens, got %d", g.path, len(g.toks), j+1)
				}
				break
			}
		}
	}
}

// tokenEqual reports whether the given tokens are equal.
func tokenEqual(t1, t2 *token.Token) bool {
	return bytes.Equal(t1.Lit, t2.Lit) && t1.Type == t2.Type && t1.Pos.Offset == t2.Pos.Offset
}

func BenchmarkLexer(b *testing.B) {
	src, err := ioutil.ReadFile("../../testdata/noisy/advanced/eval.c")
	if err != nil {
		b.Fatal(err)
	}

	for i := 0; i < b.N; i++ {
		s := lexer.NewLexer(src)
		for {
			tok := s.Scan()
			if tok.Type == token.EOF {
				break
			}
		}
	}
}
