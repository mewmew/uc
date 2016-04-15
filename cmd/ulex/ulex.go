// ulex is a lexer for the µC language which pretty-prints tokens to standard
// output.
//
// Usage: ulex FILE...
//
// If FILE is -, read standard input.
//
//   -hand
//        use hand-written lexer (default true)
//   -n int
//        number of tokens to lex
package main

import (
	"flag"
	"fmt"
	"log"
	"math"
	"os"

	"github.com/mewkiz/pkg/errutil"
	"github.com/mewmew/uc/gocc/scanner"
	gocctoken "github.com/mewmew/uc/gocc/token"
	"github.com/mewmew/uc/hand/lexer"
	"github.com/mewmew/uc/token"
)

func usage() {
	const use = `
Usage: ulex FILE...

If FILE is -, read standard input.
`
	fmt.Fprintln(os.Stderr, use[1:])
	flag.PrintDefaults()
}

func main() {
	// Command line flags.
	var (
		// hand specifies whether to use the hand-written lexer, instead of the
		// Gocc generated.
		hand bool
		// n specifies the number of tokens to lex.
		n int
	)
	flag.BoolVar(&hand, "hand", true, "use hand-written lexer")
	flag.IntVar(&n, "n", 0, "number of tokens to lex")
	flag.Usage = usage
	flag.Parse()
	if flag.NArg() == 0 {
		flag.Usage()
		os.Exit(1)
	}

	// Lex input.
	lexFile := lexFileGocc
	if hand {
		lexFile = lexFileHand
	}
	for _, path := range flag.Args() {
		if err := lexFile(path, n); err != nil {
			log.Print(err)
		}
	}
}

// lexFileHand lexes the given file and pretty-prints the n first tokens to
// standard output, using the hand-written lexer.
func lexFileHand(path string, n int) (err error) {
	var toks []token.Token
	if path == "-" {
		fmt.Fprintln(os.Stderr, "Lexing from standard input")
		toks, err = lexer.Parse(os.Stdin)
	} else {
		fmt.Fprintf(os.Stderr, "Lexing %q\n", path)
		toks, err = lexer.ParseFile(path)
	}
	if err != nil {
		return errutil.Err(err)
	}

	ntoks := len(toks)
	if n > ntoks {
		ntoks = n
	}
	pad := int(math.Ceil(math.Log10(float64(ntoks))))
	for i, tok := range toks {
		if n != 0 && i == n {
			break
		}
		if tok.Kind == token.Error {
			elog.Printf("ERROR %*d:   %v\n", pad, i, tok)
		} else {
			fmt.Printf("token %*d:   %v\n", pad, i, tok)
		}
	}
	fmt.Fprintln(os.Stderr)
	return nil
}

// lexFileGocc lexes the given file and pretty-prints the n first tokens to
// standard output, using the Gocc generated lexer.
func lexFileGocc(path string, n int) (err error) {
	var s scanner.Scanner
	if path == "-" {
		fmt.Fprintln(os.Stderr, "Lexing from standard input")
		s, err = scanner.New(os.Stdin)
	} else {
		fmt.Fprintf(os.Stderr, "Lexing %q\n", path)
		s, err = scanner.Open(path)
	}
	if err != nil {
		return errutil.Err(err)
	}

	for i := 0; ; i++ {
		if n != 0 && i == n {
			break
		}
		tok := s.Scan()
		if tok.Type == gocctoken.INVALID {
			elog.Printf("ERROR %d:   %#v\n", i, tok)
			fmt.Printf("   lit: %q\n", string(tok.Lit))
		} else {
			fmt.Printf("token %d:    %#v\n", i, tok)
			fmt.Printf("   lit: %q\n", string(tok.Lit))
		}
		if tok.Type == gocctoken.EOF {
			break
		}
	}
	fmt.Fprintln(os.Stderr)
	return nil
}

// elog represents a logger with no prefix or flags, which logs errors to
// standard error.
var elog = log.New(os.Stderr, "", 0)
