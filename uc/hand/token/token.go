package token

// Example taken from https://blog.gopheracademy.com/advent-2014/parsers-lexers/

type Token int

const (
	// Special tokens
	ILLEGAL Token = iota
	EOF
	WS // blank (32), newline (10), carriage return (13), form feed (12), and tab (9)

	IDENT // identifiers (function names, variables, ...) [a-zA-Z_][a-zA-Z0-9_]*

	// Literals
	STRING
	NUM

	// Misc characters
	NEQ       // !=
	NOT       // !
	AND       // &&
	LPARAN    // (
	RPARAN    // )
	ASTERIX   // *
	PLUS      // +
	COMMA     // , (comma)
	MINUS     // -
	DIVIDE    // /
	SEMICOLON // ;
	LE        // <=
	LT        // <
	EQ        // ==
	ASSIGN    // =
	GE        // >=
	GT        // >
	LBRACKET  // [
	RBRACKET  // ]
	LBRACE    // {
	RBRACE    // }

	// Keywords
	CHAR
	ELSE
	IF
	INT
	RETURN
	VOID // function signatures return type only
	WHILE
)

func (t Token) String() string {
	switch t {
	case ILLEGAL:
		return "ILLEGAL"
	}
	return "< token name not yet implemented >"
}
