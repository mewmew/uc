// 3rdpartycompile is a compiler for the µC language which through the uclang
// tool chain validates the input, compiles to LLVM and through clang links
// with the supplied lib uc.c and outputs the corresponding binary.
//
// Usage: 3rdpartycompile [OPTION]... FILE...
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
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/mewkiz/pkg/errutil"
	"github.com/mewkiz/pkg/goutil"
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
Usage: 3rdpartycompile [OPTION]... FILE...

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
	flag.StringVar(&outputPath, "o", "a.out", "output path")
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
	for _, path := range flag.Args() {
		err := compileFile(path, outputPath, goccLexer)
		if err != nil {
			log.Fatal(err)
		}
	}
}

// checkFile performs a static semantic analysis check on the given file.
func compileFile(path string, outputPath string, goccLexer bool) error {
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

	// Add path to uc lib.
	lib, err := goutil.SrcDir("github.com/mewmew/uc/testdata")
	if err != nil {
		return errutil.Err(err)
	}
	lib = filepath.Join(lib, "uc.ll")

	// Link and create binary through clang
	clang := exec.Command("clang", "-o", outputPath, "-x", "ir", lib, "-")
	clang.Stdin = strings.NewReader(module.String())
	clang.Stderr = os.Stderr
	clang.Stdout = os.Stdout
	if err := clang.Run(); err != nil {
		return errutil.Err(err)
	}
	return nil
}
