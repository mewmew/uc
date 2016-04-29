// usem is a static semantic checker for the µC language which validates the
// input and reports errors to standard error.
//
// Usage: usem [OPTION]... FILE...
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

	"github.com/mewkiz/pkg/errutil"
	"github.com/mewkiz/pkg/ioutilx"
	"github.com/mewmew/uc/ast"
	goccerrors "github.com/mewmew/uc/gocc/errors"
	"github.com/mewmew/uc/gocc/parser"
	goccscanner "github.com/mewmew/uc/gocc/scanner"
	handscanner "github.com/mewmew/uc/hand/scanner"
	"github.com/mewmew/uc/sem"
	semerrors "github.com/mewmew/uc/sem/errors"
)

func usage() {
	const use = `
Usage: usem [OPTION]... FILE...

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
		err := checkFile(path, hand)
		if err != nil {
			log.Print(err)
		}
	}
}

// checkFile performs a static semantic analysis check on the given file.
func checkFile(path string, hand bool) error {
	// Lexical analysis
	// Syntactic analysis (skip function bodies)
	// Top-level declarations; used for forward-declarations.
	// Syntactic analysis (including function bodies)

	// NOTE: "For each method body, we rewind the lexer to the point where the
	// method body began and parse the method body."
	//
	// ref: https://blogs.msdn.microsoft.com/ericlippert/2010/02/04/how-many-passes/

	// Semantic analysis

	// Create lexer for the input.
	buf, err := ioutilx.ReadFile(path)
	if err != nil {
		return errutil.Err(err)
	}
	if path == "-" {
		path = "<stdin>"
	}
	fmt.Fprintf(os.Stderr, "Checking %q\n", path)
	var s parser.Scanner
	if hand {
		s = handscanner.NewFromBytes(buf)
	} else {
		s = goccscanner.NewFromBytes(buf)
	}

	// Parse input.
	p := parser.NewParser()
	f, err := p.Parse(s)
	if err != nil {
		if err, ok := err.(*goccerrors.Error); ok {
			// Unwrap Gocc error.
			return parser.NewError(err)
		}
		return errutil.Err(err)
	}
	file := f.(*ast.File)
	input := string(buf)
	src := semerrors.NewSource(path, input)
	if err := sem.Check(file); err != nil {
		if err, ok := err.(*errutil.ErrInfo); ok {
			// Unwrap errutil error.
			if err, ok := err.Err.(*semerrors.Error); ok {
				// Unwrap semantic analysis error, and add input source information.
				err.Src = src
				return err
			}
		}
		return errutil.Err(err)
	}

	//pretty.Print(file)

	return nil
}
