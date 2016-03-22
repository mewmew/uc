package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/mewmew/uc/uc/hand/lexer"
)

func main() {
	flag.Parse()
	for _, path := range flag.Args() {
		err := lexFile(path)
		if err != nil {
			log.Print(err)
		}
	}
}

func lexFile(path string) error {
	toks, err := lexer.ParseFile(path)
	if err != nil {
		return err
	}
	for _, tok := range toks {
		fmt.Println(tok)
	}
	return nil
}
