// Sections from the public committee draft of the C11 ISO standard [1] have
// referenced throughout the source code (e.g. §6.4).
//
//    [1]: http://www.open-std.org/jtc1/sc22/wg14/www/docs/n1570.pdf

package lexer

import (
	"strings"
	"unicode/utf8"
)

const (
	// whitespace specifies the white-space characters: space (0x20), horizontal
	// tab (0x09), line-feed (0x0A), vertical tab (0x0B), and form-feed (0x0C)
	// (§6.4).
	//
	// Note: Even though not explicitly mentioned in the C11 specification,
	// carriage return (0x0D) is treated as a white-space character by the lexer,
	// as it is part of the CRLF newline character sequence.
	whitespace = " \t\n\v\f\r"
	// decimal specifies the decimal digit characters.
	decimal = "0123456789"
	// upper specifies the uppercase letters.
	upper = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	// lower specifies the lowercase letters.
	lower = "abcdefghijklmnopqrstuvwxyz"
	// alpha specifies the alphabetic characters.
	alpha = upper + lower
	// head is the set of valid characters for the first character of an
	// identifier.
	head = alpha + "_"
	// tail is the set of valid characters for the remaining characters of an
	// identifier (i.e. all characters in the identifier except the first).
	tail = head + decimal
)

// A stateFn represents the state of the lexer as a function that returns a
// state function.
type stateFn func(l *lexer) stateFn

// lexToken lexes a token of the µC programming language. It is the initial
// state function of the lexer.
func lexToken(l *lexer) stateFn {
	l.ignoreRun(whitespace)

	r := l.next()
	switch r {
	// Special tokens.
	case utf8.RuneError:
		// Emit error token but continue lexing next token.
		l.emitErrorf("illegal UTF-8 encoding")
		return lexToken
	case eof:
		l.emitEOF()
		// Terminate the lexer with a nil state function.
		return nil

		// Operators and delimiters.

		// Constants.
	}

	// Lex integer literals.
	if isDigit(r) {
		return lexDigit // 123
	}

	// Lex identifier or keyword.
	if r == '_' || isAlpha(r) {
		return lexLetter // foo, return
	}

	// Emit error token but continue lexing next token.
	l.emitErrorf("unexpected %q", r)
	return lexToken
}

func lexDigit(l *lexer) stateFn {
	panic("not yet implemented")
}

func lexLetter(l *lexer) stateFn {
	panic("not yet implemented")
}

// isDigit returns true if r is a digit (0-9), and false otherwise.
func isDigit(r rune) bool {
	return strings.ContainsRune(decimal, r)
}

// isAlpha returns true if r is an alphabetic character, and false otherwise.
func isAlpha(r rune) bool {
	return strings.ContainsRune(alpha, r)
}
