package token

import "testing"

type test struct {
	tok  Token
	want string
}

func TestTokenString(t *testing.T) {
	golden := []test{
		{
			tok:  ILLEGAL,
			want: "ILLEGAL",
		},
	}

	for i, g := range golden {
		got := g.tok.String()
		if got != g.want {
			t.Errorf("i=%d: expected %q, got %q", i, g.want, got)
		}
	}
}
