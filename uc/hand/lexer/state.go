// Sections from the public committee draft of the C11 ISO standard [1] have
// been referenced throughout the source code (e.g. §6.4).
//
// References:
//    [1]: http://www.open-std.org/jtc1/sc22/wg14/www/docs/n1570.pdf

package lexer

import (
	"strings"
	"unicode/utf8"

	"github.com/mewmew/uc/uc/token"
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
		l.emitErrorf("illegal UTF-8 encoding") // TODO: Check if we ever hit this case.
		return lexToken
	case eof:
		l.emitEOF()
		// Terminate the lexer with a nil state function.
		return nil

	// Operators and delimiters.
	case '(':
		l.emit(token.Lparen)
		return lexToken
	case ')':
		l.emit(token.Rparen)
		return lexToken
	case '[':
		l.emit(token.Lbracket)
		return lexToken
	case ']':
		l.emit(token.Rbracket)
		return lexToken
	case '{':
		l.emit(token.Lbrace)
		return lexToken
	case '}':
		l.emit(token.Rbrace)
		return lexToken
	case ';':
		l.emit(token.Semicolon)
		return lexToken
	case ',':
		l.emit(token.Comma)
		return lexToken
	case '+':
		l.emit(token.Add)
		return lexToken
	case '-':
		l.emit(token.Sub)
		return lexToken
	case '*':
		l.emit(token.Mul)
		return lexToken
	case '/':
		return lexSlash
	case '<':
		return lexLess
	case '>':
		return lexGreater
	case '=':
		return lexEqual
	case '&':
		return lexAmpersand
	case '\'':
		return lexCharLit
	case '!':
		return lexExclaim
	}

	// Lex integer literal.
	if isDigit(r) {
		return lexIntLit // 123
	}

	// Lex identifier or keyword.
	if r == '_' || isAlpha(r) {
		return lexLetter // foo, return
	}

	// Emit error token but continue lexing next token.
	l.emitErrorf("unexpected %#U", r)
	return lexToken
}

// lexSlash lexes a division operator (/), a line comment (//), or a block
// comment (/*). A slash character (/) has already been consumed.
func lexSlash(l *lexer) stateFn {
	switch l.next() {
	case '/':
		// Line comment (//).
		return lexLineComment
	case '*':
		// Block comment (/*).
		return lexBlockComment
	default:
		// Division operator (/).
		l.backup()
		l.emit(token.Div)
		return lexToken
	}
}

// lexLineComment lexes a line comment which acts like a newline. Two slash
// characters (//) have already been consumed.
//
//    Line comment: //[^\n]*
func lexLineComment(l *lexer) stateFn {
	for {
		switch l.next() {
		// TODO: Remove the utf8.RuneError case as comments should be able to
		// include other encodings than UTF-8, such as ISO-8859-1
		case utf8.RuneError:
			// Append error but continue lexing line comment.
			l.errorfCur("illegal UTF-8 encoding")
		case eof:
			s := l.input[l.start:l.cur]
			s = strings.TrimRight(s, "\r") // skip trailing carriage returns.
			l.emitCustom(token.Comment, s)
			l.emitEOF()
			// Terminate the lexer with a nil state function.
			return nil
		case '\n':
			s := l.input[l.start:l.cur]
			s = strings.TrimRight(s, "\r\n") // skip trailing carriage returns and newlines.
			l.emitCustom(token.Comment, s)
			return lexToken
		}
	}
}

// lexBlockComment lexes a block comment which acts like a white-space
// character. A slash and an asterisk character (/*) have already been consumed.
//
//    Block comment: "/*" .* "*/"
func lexBlockComment(l *lexer) stateFn {
	for !strings.HasSuffix(l.input[l.start+2:l.cur], "*/") {
		switch l.next() {
		// TODO: Remove the utf8.RuneError case as comments should be able to
		// include other encodings than UTF-8, such as ISO-8859-1
		case utf8.RuneError:
			// Append error but continue lexing line comment.
			l.errorfCur("illegal UTF-8 encoding")
		case eof:
			l.emitErrorf("unexpected eof in block comment")
			l.emitEOF()
			// Terminate the lexer with a nil state function.
			return nil
		}
	}

	// Strip carriage returns.
	s := strings.Replace(l.input[l.start:l.cur], "\r", "", -1)
	l.emitCustom(token.Comment, s)

	return lexToken
}

// lexAmpersand lexes a logical AND operator (&&). An ampersand (&) has already
// been consumed.
func lexAmpersand(l *lexer) stateFn {
	if r := l.next(); r == '&' {
		l.emit(token.Land)
	} else {
		// Emit error token but continue lexing next token.
		l.backup()
		l.emitErrorf("expected '&' after '&', got %#U", r)
	}
	return lexToken
}

// lexLess lexes a less-than operator (<) or a less-than-or-equal-to
// operator (<=). A less than character (<) has already been consumed.
func lexLess(l *lexer) stateFn {
	if l.accept("=") {
		l.emit(token.Le)
	} else {
		l.emit(token.Lt)
	}
	return lexToken
}

// lexGreater lexes a greater-than operator (>) or a greater-than-or-equal-to
// operator (>=). A greater than character (>) has already been consumed.
func lexGreater(l *lexer) stateFn {
	if l.accept("=") {
		l.emit(token.Ge)
	} else {
		l.emit(token.Gt)
	}
	return lexToken
}

// lexEqual lexes an assignment operator (=) or an equality operator (==). An
// equal character (=) has already been consumed.
func lexEqual(l *lexer) stateFn {
	if l.accept("=") {
		l.emit(token.Eq)
	} else {
		l.emit(token.Assign)
	}
	return lexToken
}

// lexExclaim lexes a negation operator (!) or an inequality operator (!=). An
// exclamation mark (!) has already been consumed.
//
//    Negation operator:   !
//    Inequality operator: !=
func lexExclaim(l *lexer) stateFn {
	if l.accept("=") {
		l.emit(token.Ne) // !=
	} else {
		l.emit(token.Not) // !
	}
	return lexToken
}

// TODO: if r > utf8.SelfRune { error. }

// lexCharLit lexes a character literal (e.g. 'a', '\n'). An apostrophe (') has
// already been consumed.
//
//    CharLit = "'" ([^'] | "\n" ) "'"
func lexCharLit(l *lexer) stateFn {
	// Store position directly after the token prefix, i.e. after the apostrophe.
	cur := l.cur
	// Consume character or escape sequence.
	if r := l.next(); r == '\\' {
		if !l.accept("n") {
			// TODO: Evaluate if errorCur is always used before restoring the
			// posision to the one directly after the token prefix. If that should
			// be the case, rewrite errorCur to take another arugment cur, and let
			// it restore the position.
			l.errorfCur(`unknown escape sequence '\%c'`, r)
			// Continue lexing directly after the token prefix.
			l.cur = cur
			l.ignore()
			return lexToken
		}
	}

	// Consume apostrophe.
	if !l.accept("'") {
		l.errorf("unterminated character literal")
		// Continue lexing directly after the token prefix.
		l.cur = cur
		l.ignore()
		return lexToken
	}

	l.emit(token.CharLit)
	return lexToken
}

// lexIntLit lexes an integer literal (e.g. 123). A decimal digit (0-9) has
// already been consumed.
//
//    IntLit = [0-9]+
func lexIntLit(l *lexer) stateFn {
	l.acceptRun(decimal)
	l.emit(token.IntLit)
	return lexToken
}

// lexLetter lexes an identifier (e.g. main) or a keyword (e.g. return). An
// alphabetic character (a-z or A-Z) or an underscore (_) has already been
// consumed.
//
//    Identifier = [a-zA-Z_][a-zA-Z0-9_]*
//    Keyword    = if, return, …
func lexLetter(l *lexer) stateFn {
	l.acceptRun(alpha + decimal + "_")
	// TODO: Optimize using binary search on keyword.
	if kind, ok := token.Keywords[l.input[l.start:l.cur]]; ok {
		l.emit(kind)
	} else {
		l.emit(token.Ident)

	}
	return lexToken
}

// isDigit returns true if r is a digit (0-9), and false otherwise.
func isDigit(r rune) bool {
	return strings.ContainsRune(decimal, r)
}

// isAlpha returns true if r is an alphabetic character, and false otherwise.
func isAlpha(r rune) bool {
	return strings.ContainsRune(alpha, r)
}
