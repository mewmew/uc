// Package token defines constants representing the lexical tokens of the µC
// programming language.
package token

import "fmt"

// A Token represents a lexical token of the µC programming language.
type Token struct {
	// The token type.
	Kind
	// The string value of the token.
	Val string
	// Start position in the input string.
	Pos int
}

func (tok Token) String() string {
	return fmt.Sprintf(`token.Token{Kind: token.%v, Val: %q, Pos: %v}`, tok.Kind, tok.Val, tok.Pos)
}

//go:generate stringer -type Kind

// Kind is the set of lexical token types of the µC programming language.
type Kind uint16

// Token types.
const (
	// Special tokens.
	EOF     Kind = iota // End of file
	Error               // Token value holds an error message (e.g. unterminated string)
	Comment             // /* block comment */ or // line comment

	literalStart

	// Identifiers and basic literals.
	Ident   // main (also includes type names)
	IntLit  // 123
	CharLit // 'a', '\n'

	literalEnd

	operatorStart

	// Operators and delimiters.
	Add       // +
	Sub       // -
	Mul       // *
	Div       // /
	Assign    // =
	Eq        // ==
	Ne        // !=
	Lt        // <
	Le        // <=
	Gt        // >
	Ge        // >=
	Land      // &&
	Not       // !
	Lparen    // (
	Rparen    // )
	Lbracket  // [
	Rbracket  // ]
	Lbrace    // {
	Rbrace    // }
	Comma     // ,
	Semicolon // ;

	operatorEnd

	keywordStart

	// Keywords.
	KwElse   // else
	KwIf     // if
	KwReturn // return
	KwWhile  // while

	keywordEnd
)

// IsKeyword reports whether kind is a keyword.
func (kind Kind) IsKeyword() bool {
	return keywordStart < kind && kind < keywordEnd
}

// IsLiteral reports whether kind is an identifier or a basic literal.
func (kind Kind) IsLiteral() bool {
	return literalStart < kind && kind < literalEnd
}

// IsOperator reports whether kind is an operator or a delimiter.
func (kind Kind) IsOperator() bool {
	return operatorStart < kind && kind < operatorEnd
}

// Keywords is the set of valid keywords in the µC programming language.
var Keywords = map[string]Kind{
	"else":   KwElse,
	"if":     KwIf,
	"return": KwReturn,
	"while":  KwWhile,
}
