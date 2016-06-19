// uparse is a parser for the µC language which pretty-prints abstract syntax
// trees to standard output.
//
// Usage: uparse [OPTION]... FILE...
//
// If FILE is -, read standard input.
//
//   -gocc-lexer
//        use Gocc generated lexer
package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/davecgh/go-spew/spew"
	"github.com/kr/pretty"
	"github.com/mewkiz/pkg/errutil"
	"github.com/mewkiz/pkg/ioutilx"
	"github.com/mewmew/uc/ast"
	"github.com/mewmew/uc/gocc/errors"
	"github.com/mewmew/uc/gocc/parser"
	goccscanner "github.com/mewmew/uc/gocc/scanner"
	handscanner "github.com/mewmew/uc/hand/scanner"
)

func usage() {
	const use = `
Usage: uparse [OPTION]... FILE...

If FILE is -, read standard input.
`
	fmt.Fprintln(os.Stderr, use[1:])
	flag.PrintDefaults()
}

func main() {
	var (
		// goccLexer specifies whether to use the Gocc generated lexer, instead of
		// the hand-written lexer.
		goccLexer bool
	)
	flag.BoolVar(&goccLexer, "gocc-lexer", false, "use Gocc generated lexer")
	flag.Usage = usage
	flag.Parse()
	if flag.NArg() == 0 {
		flag.Usage()
		os.Exit(1)
	}

	// Parse input.
	for _, path := range flag.Args() {
		err := parseFile(path, goccLexer)
		if err != nil {
			log.Print(err)
		}
	}
}

// parseFile parses the given file and pretty-prints its abstract syntax tree to
// standard output, optionally using the Gocc generated lexer.
func parseFile(path string, goccLexer bool) error {
	// Create lexer for the input.
	buf, err := ioutilx.ReadFile(path)
	if err != nil {
		return errutil.Err(err)
	}
	if path == "-" {
		fmt.Fprintln(os.Stderr, "Parsing from standard input")
	} else {
		fmt.Fprintf(os.Stderr, "Parsing %q\n", path)
	}
	var s parser.Scanner
	if goccLexer {
		s = goccscanner.NewFromBytes(buf)
	} else {
		s = handscanner.NewFromBytes(buf)
	}

	// Parse input.
	p := parser.NewParser()
	file, err := p.Parse(s)
	if err != nil {
		if err, ok := err.(*errors.Error); ok {
			// Unwrap Gocc error.
			return parser.NewError(err)
		}
		return errutil.Err(err)
	}
	f := file.(*ast.File)
	for _, decl := range f.Decls {
		fmt.Println("=== [ Top-level declaration ] ===")
		fmt.Println()
		fmt.Printf("decl type: %T\n", decl)
		fmt.Println()
		fmt.Println("decl:", decl)
		fmt.Println()
		pretty.Print(decl)
		fmt.Println()
		spew.Print(decl)
		fmt.Println()
		fmt.Println()
	}

	return nil
}
