// uclang is compiler for the µC language which validates the input, generates
// corresponding LLVM IR assembly and reports errors to standard error.
//
// Usage: uclang [OPTION]... FILE...
//
// If FILE is -, read standard input.
//
//   -gocc-lexer
//        use Gocc generated lexer
//   -no-colors
//        disable colors in output
//   -no-nested-functions
//        disable support for nested functions
//   -o string
//        output path
package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"

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
		// goccLexer specifies whether to use the Gocc generated lexer, instead of
		// the hand-written lexer.
		goccLexer bool
		// noColors specifies whether to disable colors in output.
		noColors bool
		// outputPath specifies the output path for the generated LLVM IR.
		outputPath string
	)
	flag.BoolVar(&goccLexer, "gocc-lexer", false, "use Gocc generated lexer")
	flag.BoolVar(&noColors, "no-colors", false, "disable colors in output")
	flag.BoolVar(&semcheck.NoNestedFunctions, "no-nested-functions", false, "disable support for nested functions")
	flag.StringVar(&outputPath, "o", "", "output path")
	flag.Usage = usage
	flag.Parse()
	semerrors.UseColor = !noColors
	if flag.NArg() == 0 {
		flag.Usage()
		os.Exit(1)
	}
	// TODO: Remove once nested functions are supported. For now, disallow during
	// semantic analysis.
	semcheck.NoNestedFunctions = true

	// Parse input.
	output := os.Stdout
	if len(outputPath) > 0 {
		var err error
		output, err = os.Create(outputPath)
		if err != nil {
			log.Fatal(errutil.Err(err))
		}
		defer output.Close()
	}
	for _, path := range flag.Args() {
		err := compileFile(path, output, goccLexer)
		if err != nil {
			log.Fatal(err)
		}
	}
}

// checkFile performs a static semantic analysis check on the given file.
func compileFile(path string, output io.Writer, goccLexer bool) error {
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

	// Generate LLVM IR module based on the syntax tree of the given file.
	module := irgen.Gen(file, info)
	if _, err := fmt.Fprint(output, module); err != nil {
		return errutil.Err(err)
	}

	return nil
}
