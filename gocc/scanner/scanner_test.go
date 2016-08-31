package scanner_test

import (
	"bytes"
	"io/ioutil"
	"log"
	"testing"

	"github.com/mewmew/uc/gocc/scanner"
	"github.com/mewmew/uc/gocc/token"
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
					Type: token.TokMap.Type("$"),
					Lit:  []byte(""),
					// Note, a new line character has been inserted to ensure that
					// the file ends with a new line; thus the actual offset of the
					// EOF is 1099.
					Pos: token.Pos{Offset: 1100},
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
					Type: token.TokMap.Type("char_lit"),
					Lit:  []byte("'c'"),
					Pos:  token.Pos{Offset: 33},
				},
				{
					Type: token.TokMap.Type(";"),
					Lit:  []byte(";"),
					Pos:  token.Pos{Offset: 36},
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
					Type: token.TokMap.Type("$"),
					Lit:  []byte(""),
					// Note, a new line character has been inserted to ensure that
					// the file ends with a new line; thus the actual offset of the
					// EOF is 533.
					Pos: token.Pos{Offset: 534},
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
					Type: token.TokMap.Type("char_lit"),
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
					Type: token.TokMap.Type("char_lit"),
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
					Type: token.TokMap.Type("char_lit"),
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
					Type: token.TokMap.Type("char_lit"),
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
		s, err := scanner.Open(g.path)
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
		s := scanner.NewFromBytes(src)
		for {
			tok := s.Scan()
			if tok.Type == token.EOF {
				break
			}
		}
	}
}
