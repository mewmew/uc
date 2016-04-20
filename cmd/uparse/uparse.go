// uparse is a parser for the µC language which pretty-prints abstract syntax
// trees to standard output.
//
// Usage: uparse [OPTION]... FILE...
//
// If FILE is -, read standard input.
//
//   -hand
//        use hand-written lexer (default true)
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
		// hand specifies whether to use the hand-written lexer, instead of the
		// Gocc generated.
		hand bool
	)
	flag.BoolVar(&hand, "hand", true, "use hand-written lexer")
	flag.Usage = usage
	flag.Parse()
	if flag.NArg() == 0 {
		flag.Usage()
		os.Exit(1)
	}

	// Parse input.
	for _, path := range flag.Args() {
		err := parseFile(path, hand)
		if err != nil {
			log.Print(err)
		}
	}
}

// parseFile parses the given file and pretty-prints its abstract syntax tree to
// standard output, optionally using the hand-written lexer.
func parseFile(path string, hand bool) error {
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
	if hand {
		s = handscanner.NewFromBytes(buf)
	} else {
		s = goccscanner.NewFromBytes(buf)
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
