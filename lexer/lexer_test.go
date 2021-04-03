package lexer

import "testing"

func assertEqual(t *testing.T, a interface{}, b interface{}) {
	t.Helper()
	if a != b {
		t.Fatalf("%s != %s", a, b)
	}
}

func expectNumber(t *testing.T, contents string, expected float64) {
	t.Helper()
	l := NewLexer(contents)
	assertEqual(t, l.Token, TNumber)
	assertEqual(t, l.Number, expected)
}

func TestNumber(t *testing.T) {
	expectNumber(t, "0", 0.0)
	expectNumber(t, "123", 123.0)
	expectNumber(t, "1289.345", 1289.345)
}

func expectString(t *testing.T, contents string, expected string) {
	t.Helper()
	l := NewLexer(contents)
	assertEqual(t, l.Token, TString)
	assertEqual(t, l.String, expected)
}

func TestString(t *testing.T) {
	expectString(t, "\"\"", "")
	expectString(t, "\"123\"", "123")
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
		{",", TComma},
		{":", TColon},

		// Literal tokens
		{"true", TTrue},
		{"false", TFalse},
		{"null", TNull},
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

func expectLexerError(t *testing.T, contents string, expected string) {
	t.Helper()
	t.Run(contents, func(t *testing.T) {
		t.Helper()
		text := ""
		func() {
			defer func() {
				r := recover()
				lexerPanic, isLexerPanic := r.(LexerPanic)
				if r != nil && !isLexerPanic {
					panic(r)
				} else {
					text = lexerPanic.Msg
				}
			}()
			NewLexer(contents)
		}()
		assertEqual(t, text, expected)
	})
}

func TestSyntaxError(t *testing.T) {
	expectLexerError(t, "(", "Unexpected token (")
	expectLexerError(t, "trae", "Unexpected token a")
	expectLexerError(t, "faise", "Unexpected token i")
	expectLexerError(t, "00", "Unexpected number")
}
