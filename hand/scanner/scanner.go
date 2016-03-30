// Package scanner implements a lexer which satifies the Gocc Scanner interface.
package scanner

import (
	"io"

	"github.com/mewmew/uc/gocc/token"
	"github.com/mewmew/uc/hand/lexer"
	uctoken "github.com/mewmew/uc/token"
)

// Scanner represents the lexer interface used by the Gocc parser.
type Scanner interface {
	// Scan lexes and returns the next token of the source input.
	Scan() *token.Token
}

// a scanner is a lexer which implements the Gocc Scanner interface.
type scanner struct {
	// Lexed tokens.
	toks []uctoken.Token
	// Current token.
	cur int
}

// Ensure that scanner implements the Gocc Scanner interface.
var _ Scanner = &scanner{}

// eof represents an end-of-file token with unknown source position.
var eof = &token.Token{Type: token.TokMap.Type("$")}

// Scan lexes and returns the next token of the source input.
func (s *scanner) Scan() *token.Token {
	if s.cur >= len(s.toks) {
		return eof
	}
	tok := s.toks[s.cur]
	s.cur++
	var typ token.Type
	switch tok.Kind {
	case uctoken.EOF:
		typ = token.TokMap.Type("$")
	case uctoken.Error:
		typ = token.TokMap.Type("INVALID")
	case uctoken.Comment:
		typ = token.TokMap.Type("comment")
	case uctoken.Ident:
		typ = token.TokMap.Type("ident")
	case uctoken.IntLit:
		typ = token.TokMap.Type("int_lit")
	case uctoken.CharLit:
		typ = token.TokMap.Type("char_lit")
	default:
		typ = token.TokMap.Type(tok.Val)
	}
	lit := []byte(tok.Val)
	pos := token.Pos{Offset: tok.Pos}
	return &token.Token{
		Type: typ,
		Lit:  lit,
		Pos:  pos,
	}
}

// New returns a new scanner lexing from r.
func New(r io.Reader) (Scanner, error) {
	toks, err := lexer.Parse(r)
	if err != nil {
		return nil, err
	}
	return &scanner{toks: toks}, nil
}

// Open returns a new scanner lexing from path.
func Open(path string) (Scanner, error) {
	toks, err := lexer.ParseFile(path)
	if err != nil {
		return nil, err
	}
	return &scanner{toks: toks}, nil
}

// NewFromString returns a new scanner lexing from input.
func NewFromString(input string) Scanner {
	toks := lexer.ParseString(input)
	return &scanner{toks: toks}
}
