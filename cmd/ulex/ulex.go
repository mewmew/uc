package main

import (
	"flag"
	"fmt"
	"github.com/mewmew/uc/uc/hand/lexer"
	"log"
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
		fmt.Println(token)
	}
	return nil
}
