package lexer

import "testing"

func assertEqual(t *testing.T, a interface{}, b interface{}) {
	t.Helper()
	if a != b {
		t.Fatalf("%s != %s", a, b)
	}
}

func TestTokens(t *testing.T) {
	expected := []struct {
		contents string
		token    Token
	}{
		{"", TEndOfFile},

		// Structural tokens
		{"[", TLeftSquareBracket},
		{"]", TRightSquareBracket},
		{"{", TLeftCurlyBracket},
		{"}", TRightCurlyBracket},

		// Literal tokens
		// {"true", TTrue},
		// {"false", TFalse},
		// {"null", TNull},
	}

	for _, it := range expected {
		contents := it.contents
		token := it.token
		t.Run(contents, func(t *testing.T) {
			l := NewLexer(contents)
			assertEqual(t, l.Token, token)
		})
	}
}
