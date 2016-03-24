package lexer_test

import (
	"fmt"
	"log"

	"github.com/mewmew/uc/uc/hand/lexer"
)

func ExampleParseFile() {
	tokens, err := lexer.ParseFile("../../testdata/incorrect/lexer/long-char.c")
	if err != nil {
		log.Fatal(err)
	}
	for _, tok := range tokens {
		fmt.Printf("%-9v   %q\n", tok.Kind, tok.Val)
	}
	// Output:
	// Ident       "int"
	// Ident       "main"
	// Lparen      "("
	// Ident       "void"
	// Rparen      ")"
	// Lbrace      "{"
	// Ident       "char"
	// Ident       "c"
	// Semicolon   ";"
	// Ident       "c"
	// Assign      "="
	// CharLit     "'c'"
	// Semicolon   ";"
	// Comment     "// OK"
	// Ident       "c"
	// Assign      "="
	// Error       "unterminated character literal"
	// Ident       "cc"
	// Error       "unterminated character literal"
	// Semicolon   ";"
	// Comment     "// Not OK"
	// Rbrace      "}"
	// EOF         ""
}

func ExampleParseString() {
	tokens := lexer.ParseString("int main (void) { ; }")
	for i, tok := range tokens {
		fmt.Printf("token %d: %v\n", i, tok)
	}
	// Output:
	// token 0: token.Token{Kind: token.Ident, Val: "int", Pos: 0}
	// token 1: token.Token{Kind: token.Ident, Val: "main", Pos: 4}
	// token 2: token.Token{Kind: token.Lparen, Val: "(", Pos: 9}
	// token 3: token.Token{Kind: token.Ident, Val: "void", Pos: 10}
	// token 4: token.Token{Kind: token.Rparen, Val: ")", Pos: 14}
	// token 5: token.Token{Kind: token.Lbrace, Val: "{", Pos: 16}
	// token 6: token.Token{Kind: token.Semicolon, Val: ";", Pos: 18}
	// token 7: token.Token{Kind: token.Rbrace, Val: "}", Pos: 20}
	// token 8: token.Token{Kind: token.EOF, Val: "", Pos: 21}
}
