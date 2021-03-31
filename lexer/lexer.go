package lexer

import (
	"unicode/utf8"
)

type Token uint

const (
	TEndOfFile Token = iota

	// Structural tokens
	TLeftSquareBracket  // [
	TLeftCurlyBracket   // {
	TRightSquareBracket // ]
	TRightCurlyBracket  // }
	TColon              // ;
	TComma              // ,

	// Literal tokens
	TTrue  // true
	TFalse // false
	TNull  // null
	TNumber
	TString
)

type Lexer struct {
	source    string
	current   int
	start     int
	end       int
	Token     Token
	codePoint rune
	Number    float64
	String    string
}

type LexerPanic struct{}

func (l *Lexer) Next() {
	switch l.codePoint {
	case -1:
		l.Token = TEndOfFile

	case '\r', '\n', '\u2028', '\u2029':
		l.step()

	case '[':
		l.step()
		l.Token = TLeftSquareBracket

	case ']':
		l.step()
		l.Token = TRightSquareBracket

	case '{':
		l.step()
		l.Token = TLeftCurlyBracket

	case '}':
		l.step()
		l.Token = TRightCurlyBracket

	case '-', '.', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		// TODO
	case '"':
		// TODO
	case 't', 'f', 'n':
		// TODO
	}
}

func (l *Lexer) step() {
	codePoint, width := utf8.DecodeRuneInString(l.source[l.current:])

	// 使用 -1 表示文件的结尾
	if width == 0 {
		codePoint = -1
	}

	l.codePoint = codePoint
	l.end = l.current
	l.current += width
}

func NewLexer(s string) Lexer {
	lexer := Lexer{
		source: s,
	}
	lexer.step()
	lexer.Next()
	return lexer
}
