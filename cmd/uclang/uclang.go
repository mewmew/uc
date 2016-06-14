// uclang is compiler for the µC language which validates the input, creates
// instructions from the input and reports errors to standard error.
//
// Usage: uclang [OPTION]... FILE...
//
// If FILE is -, read standard input.
//
//   -gocc-lexer
//        use the gocc uC lexer (default false)
//   -no-colors
//        do not use ANSI escape codes in output (default false)
//   -no-nested-functions
//        do not allow nested function declarations (default false)
package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/kr/pretty"
	"github.com/mewkiz/pkg/errutil"
	"github.com/mewkiz/pkg/ioutilx"
	"github.com/mewmew/uc/ast"
	goccerrors "github.com/mewmew/uc/gocc/errors"
	"github.com/mewmew/uc/gocc/parser"
	goccscanner "github.com/mewmew/uc/gocc/scanner"
	handscanner "github.com/mewmew/uc/hand/scanner"
	"github.com/mewmew/uc/irgen"
	"github.com/mewmew/uc/sem"
	semerrors "github.com/mewmew/uc/sem/errors"
	"github.com/mewmew/uc/sem/semcheck"
)

func usage() {
	const use = `
Usage: uclang [OPTION]... FILE...

If FILE is -, read standard input.
`
	fmt.Fprintln(os.Stderr, use[1:])
	flag.PrintDefaults()
}

func main() {
	var (
		// hand specifies whether to use the hand-written lexer, instead of the
		// Gocc generated.
		useGoccLexer bool
		noColors     bool
	)
	flag.BoolVar(&useGoccLexer, "gocc-lexer", false, "use gocc lexer")
	flag.BoolVar(&noColors, "no-colors", false, "use colors in output")
	flag.BoolVar(&semcheck.NoNestedFunctions, "no-nested-functions", false, "use colors in output")
	flag.Usage = usage
	flag.Parse()
	semerrors.UseColor = !noColors
	if flag.NArg() == 0 {
		flag.Usage()
		os.Exit(1)
	}

	// Parse input.
	for _, path := range flag.Args() {
		err := compileFile(path, useGoccLexer)
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
func compileFile(path string, useGoccLexer bool) error {
	// Lexical analysis
	// Syntactic analysis
	// Semantic analysis
	// Intermediate representation generation

	// Create lexer for the input.
	buf, err := ioutilx.ReadFile(path)
	if err != nil {
		return errutil.Err(err)
	}
	if path == "-" {
		path = "<stdin>"
	}

	fmt.Fprintf(os.Stderr, "Compiling %q\n", path)

	var s parser.Scanner
	if useGoccLexer {
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
	info, err := sem.Check(file)
	if err != nil {
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

	module := irgen.Gen(file, info)
	log.Println("=== [ Pretty module ] ===\n")
	pretty.Println(module)
	log.Println("=== [ LLVM IR module ] ===\n")
	fmt.Println(module)

	return nil
}

// elog represents a logger with no prefix or flags, which logs errors to
// standard error.
var elog = log.New(os.Stderr, "", 0)
