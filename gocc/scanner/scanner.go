// Package scanner implements a lexer which satifies the Gocc Scanner interface.
package scanner

import (
	"bytes"
	"io"
	"io/ioutil"
	"strings"

	"github.com/mewkiz/pkg/errutil"
	"github.com/mewmew/uc/gocc/lexer"
	"github.com/mewmew/uc/gocc/token"
)

// Scanner represents the lexer interface used by the Gocc parser.
type Scanner interface {
	// Scan lexes and returns the next token of the source input.
	Scan() *token.Token
}

const (
	// whitespace specifies the white-space characters: space (0x20), horizontal
	// tab (0x09), new line (line-feed (0x0A) or carriage-return (0x0D)),
	// vertical tab (0x0B), and form-feed (0x0C) (§6.4).
	//
	// Note: Even though not explicitly mentioned in the C11 specification,
	// carriage return (0x0D) is treated as a white-space character by the lexer,
	// as it is part of the CRLF newline character sequence.
	whitespace = " \t\n\v\f\r"
)

// New returns a new scanner lexing from r.
func New(r io.Reader) (Scanner, error) {
	input, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, errutil.Err(err)
	}
	return NewFromBytes(input), nil
}

// Open returns a new scanner lexing from path.
func Open(path string) (Scanner, error) {
	input, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errutil.Err(err)
	}
	return NewFromBytes(input), nil
}

// NewFromString returns a new scanner lexing from input.
func NewFromString(input string) Scanner {
	return NewFromBytes([]byte(input))
}

// NewFromBytes returns a new scanner lexing from input.
func NewFromBytes(input []byte) Scanner {
	// Append new line to files not ending with new line.
	appendNewLine := false
	lastNewLine := bytes.LastIndexByte(input, '\n')
	for _, r := range string(input[lastNewLine+1:]) {
		if !strings.ContainsRune(whitespace, r) {
			// Non-whitespace character located after the last new line character.
			appendNewLine = true
			break
		}
	}
	if appendNewLine {
		input = append(input, '\n')
	}

	return lexer.NewLexer(input)
}
