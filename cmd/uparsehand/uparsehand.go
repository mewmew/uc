// uparse is a parser for the µC language which pretty-prints abstract syntax
// trees to standard output.
//
// Usage: uparse FILE...
//
// If FILE is -, read standard input.
package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/davecgh/go-spew/spew"
	"github.com/kr/pretty"
	"github.com/mewkiz/pkg/errutil"
	"github.com/mewmew/uc/ast"
	"github.com/mewmew/uc/cmd/internal/ioutilx"
	"github.com/mewmew/uc/gocc/parser"
	"github.com/mewmew/uc/gocc/token"
	// TODO: Use hand-written scanner instead of Gocc-generated lexer once the
	// grammar has matured.
	"github.com/mewmew/uc/hand/scanner"
)

func usage() {
	const use = `
Usage: uparse FILE...

If FILE is -, read standard input.`
	fmt.Fprintln(os.Stderr, use[1:])
}

func main() {
	flag.Usage = usage
	flag.Parse()
	if flag.NArg() == 0 {
		flag.Usage()
		os.Exit(1)
	}
	for _, path := range flag.Args() {
		err := parseFile(path)
		if err != nil {
			log.Print(err)
		}
	}
}

// parseFile parses the given file and pretty-prints its abstract syntax tree to
// standard output.
func parseFile(path string) (err error) {
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
	s, err := scanner.New(bytes.NewReader(buf))
	if err != nil {
		return errutil.Err(err)
	}

	// Parse input.
	p := parser.NewParser()
	file, err := p.Parse(NewScanner(s))
	if err != nil {
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

// A Scanner wraps a parser.Scanner to produce debug output while scanning.
type Scanner struct {
	s parser.Scanner
}

// NewScanner returns a debug scanner which produces debug output while scanner
// from s.
func NewScanner(s parser.Scanner) parser.Scanner {
	return &Scanner{s: s}
}

// Scan scans from the underlying scanner and prints scanned tokens to standard
// output.
func (ds *Scanner) Scan() *token.Token {
	tok := ds.s.Scan()
	if tok.Type == token.TokMap.Type("comment") {
		// Ignoring comments
		log.Println("Ignoring comment:", string(tok.Lit))
		return ds.Scan()
	}
	return tok
}
