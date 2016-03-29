// ulex is a lexer for the µC language which pretty-prints tokens to standard
// output.
//
// Usage: ulex FILE...
//
// If FILE is -, read standard input.
package main

import (
	"flag"
	"fmt"
	"log"
	"math"
	"os"

	"github.com/mewmew/uc/uc/hand/lexer"
	"github.com/mewmew/uc/uc/token"
)

func usage() {
	const use = `
Usage: ulex FILE...

If FILE is -, read standard input.`
	fmt.Fprintln(os.Stderr, use[1:])
}

func main() {
	// Command line flags.
	var (
		// n specifies the number of tokens to lex.
		n int
	)
	flag.IntVar(&n, "n", 0, "number of tokens to lex")
	flag.Usage = usage
	flag.Parse()
	paths := flag.Args()
	if flag.NArg() == 0 {
		flag.Usage()
		os.Exit(1)
	}
	for _, path := range paths {
		err := lexFile(path, n)
		if err != nil {
			log.Print(err)
		}
	}
}

// lexFile lexes the given file and pretty-prints the n first tokens to standard
// output.
func lexFile(path string, n int) (err error) {
	var l lexer.Lexer
	if path == "-" {
		fmt.Fprintln(os.Stderr, "Lexing from standard input")
		l, err = lexer.UlexParse(os.Stdin)
	} else {
		fmt.Fprintf(os.Stderr, "Lexing %q\n", path)
		l, err = lexer.UlexParseFile(path)
	}
	if err != nil {
		return err
	}
	toks := l.Tokens()
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
			line, col := l.Position(tok.Pos)
			elog.Printf("%*s   line:%v, col:%v\n", col, "^", line, col)
		} else {
			fmt.Printf("token %*d:   %v\n", pad, i, tok)
			line, col := l.Position(tok.Pos)
			fmt.Printf("%*s   line:%v, col:%v\n", col, "^", line, col)
		}
	}
	fmt.Fprintln(os.Stderr)
	return nil
}

// elog represents a logger with no prefix or flags, which logs errors to
// standard error.
var elog = log.New(os.Stderr, "", 0)
