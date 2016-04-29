package lexer_test

import (
	"io/ioutil"
	"testing"

	"github.com/mewmew/uc/gocc/lexer"
	"github.com/mewmew/uc/gocc/token"
)

func BenchmarkLexer(b *testing.B) {
	buf, err := ioutil.ReadFile("../../testdata/noisy/advanced/eval.c")
	if err != nil {
		b.Fatal(err)
	}
	b.SetBytes(int64(len(buf)))
	for i := 0; i < b.N; i++ {
		l := lexer.NewLexer(buf)
		for {
			tok := l.Scan()
			if tok.Type == token.EOF {
				break
			}
		}
	}
}
