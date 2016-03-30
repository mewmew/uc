// uparse is a parser for the µC language which pretty-prints abstract syntax
// trees to standard output.
//
// Usage: uparse FILE...
//
// If FILE is -, read standard input.
package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/mewkiz/pkg/errutil"
	"github.com/mewmew/uc/gocc/parser"
	"github.com/mewmew/uc/gocc/token"
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
	paths := flag.Args()
	for _, path := range paths {
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
	var s parser.Scanner
	if path == "-" {
		fmt.Fprintln(os.Stderr, "Parsing from standard input")
		s, err = scanner.New(os.Stdin)
	} else {
		fmt.Fprintf(os.Stderr, "Parsing %q\n", path)
		s, err = scanner.Open(path)
	}
	if err != nil {
		return errutil.Err(err)
	}

	// Create debug scanner.
	ds := newDebugScanner(s)

	// Parse input.
	p := parser.NewParser()
	foo, err := p.Parse(ds)
	if err != nil {
		return errutil.Err(err)
	}
	fmt.Println("foo:", foo)

	return nil
}

// A debugScanner wraps a parser.Scanner to produce debug output while scanning.
type debugScanner struct {
	s parser.Scanner
}

// newDebugScanner returns a debug scanner which produces debug output while
// scanner from s.
func newDebugScanner(s parser.Scanner) parser.Scanner {
	return &debugScanner{s: s}
}

// Scan scans from the underlying scanner and prints scanned tokens to standard
// output.
func (ds *debugScanner) Scan() *token.Token {
	tok := ds.s.Scan()
	fmt.Println("pos:", tok.Pos.Offset)
	fmt.Println("typ:", tok.Type)
	fmt.Println("lit:", string(tok.Lit))
	fmt.Println()
	return tok
}
