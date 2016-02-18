package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/kr/pretty"
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
	tokens, err := lexer.ParseFile(path)
	if err != nil {
		return err
	}
	for _, token := range tokens {
		fmt.Printf("=== [ %s ] ===\n", token.Kind)
		pretty.Println(token)
		fmt.Println()
	}
	return nil
}
