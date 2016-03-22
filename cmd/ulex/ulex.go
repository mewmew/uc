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
	"github.com/mewmew/uc/uc/hand/token"
)

func usage() {
	const use = `
Usage: ulex FILE...

If FILE is -, read standard input.`
	fmt.Fprintln(os.Stderr, use[1:])
}

func main() {
	flag.Usage = usage
	flag.Parse()
	paths := flag.Args()
	if flag.NArg() == 0 {
		flag.Usage()
		os.Exit(1)
	}
	for _, path := range paths {
		err := lexFile(path)
		if err != nil {
			log.Print(err)
		}
	}
}

// lexFile lexes the given file and pretty-prints its tokens to standard output.
func lexFile(path string) (err error) {
	var toks []token.Token
	if path == "-" {
		fmt.Println("Lexing from standard input")
		toks, err = lexer.Parse(os.Stdin)
		if err != nil {
			return err
		}
	} else {
		fmt.Printf("Lexing %q\n", path)
		toks, err = lexer.ParseFile(path)
		if err != nil {
			return err
		}
	}
	fmt.Println()
	pad := int(math.Ceil(math.Log10(float64(len(toks)))))
	for i, tok := range toks {
		if tok.Kind == token.Error {
			elog.Printf("ERROR %*d:   %v\n", pad, i, tok)
		} else {
			fmt.Printf("token %*d:   %v\n", pad, i, tok)
		}
	}
	fmt.Println()
	return nil
}

// elog represents a logger with no prefix or flags, which logs errors to
// standard error.
var elog = log.New(os.Stderr, "", 0)
