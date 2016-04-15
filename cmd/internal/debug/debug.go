// Package debug implements lexer and parser debugging utility functions.
package debug

import (
	"fmt"

	"github.com/mewmew/uc/gocc/parser"
	"github.com/mewmew/uc/gocc/token"
)

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
	fmt.Println("pos:", tok.Pos.Offset)
	fmt.Println("typ:", tok.Type)
	fmt.Println("lit:", string(tok.Lit))
	fmt.Println()
	return tok
}
