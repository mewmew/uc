package main

import (
	"flag"
	"fmt"
	"log"
	"math"

	"github.com/mewmew/uc/uc/hand/lexer"
)

func main() {
	flag.Parse()
	for _, path := range flag.Args() {
		fmt.Printf("Lexing %q\n", path)
		fmt.Println()
		err := lexFile(path)
		if err != nil {
			log.Print(err)
		}
		fmt.Println()
	}
}

func lexFile(path string) error {
	toks, err := lexer.ParseFile(path)
	if err != nil {
		return err
	}
	pad := int(math.Ceil(math.Log10(float64(len(toks)))))
	for i, tok := range toks {
		fmt.Printf("token %*d:   %v\n", pad, i, tok)
	}
	return nil
}
