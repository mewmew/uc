// usem is a static semantic checker for the µC language which validates the
// input and reports errors to standard error.
//
// Usage: usem [OPTION]... FILE...
//
// If FILE is -, read standard input.
//
//   -gocc-lexer
//        use Gocc generated lexer
//   -no-colors
//        disable colors in output
//   -no-nested-functions
//        disable support for nested functions
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
	"github.com/mewmew/uc/sem/semcheck"
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
		// goccLexer specifies whether to use the Gocc generated lexer, instead of
		// the hand-written lexer.
		goccLexer bool
		// noColors specifies whether to disable colors in output.
		noColors bool
	)
	flag.BoolVar(&goccLexer, "gocc-lexer", false, "use Gocc generated lexer")
	flag.BoolVar(&noColors, "no-colors", false, "disable colors in output")
	flag.BoolVar(&semcheck.NoNestedFunctions, "no-nested-functions", false, "disable support for nested functions")
	flag.Usage = usage
	flag.Parse()
	semerrors.UseColor = !noColors
	if flag.NArg() == 0 {
		flag.Usage()
		os.Exit(1)
	}

	// Parse input.
	for _, path := range flag.Args() {
		err := checkFile(path, goccLexer)
		if err != nil {
			if _, ok := err.(*semerrors.Error); ok {
				elog.Print(err)
			} else {
				log.Print(err)
			}
		}
	}
}

// checkFile performs a static semantic analysis check on the given file.
func checkFile(path string, goccLexer bool) error {
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
	if goccLexer {
		s = goccscanner.NewFromBytes(buf)
	} else {
		s = handscanner.NewFromBytes(buf)
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
	if _, err := sem.Check(file); err != nil {
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

	return nil
}

// elog represents a logger with no prefix or flags, which logs errors to
// standard error.
var elog = log.New(os.Stderr, "", 0)
