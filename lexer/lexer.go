package lexer

import (
	"fmt"
	"strconv"
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

var KeywordToString = map[Token]string{
	TTrue:  "true",
	TFalse: "false",
	TNull:  "null",
}

type Loc struct {
	// This is the 0-based index of this location from the start of the file, in bytes
	Start int32
}

type LexerPanic struct {
	Msg string
	Loc Loc
}

func IsWhitespace(codePoint rune) bool {
	switch codePoint {
	case
		'\u0009', // character tabulation
		'\u000B', // line tabulation
		'\u000C', // form feed
		'\u0020', // space
		'\u00A0': // no-break space

		return true

	default:
		return false
	}
}

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

func (l *Lexer) Next() {
	for {
		l.start = l.end
		l.Token = 0

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

		case ',':
			l.step()
			l.Token = TComma

		case ':':
			l.step()
			l.Token = TColon

		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			l.parseNumber()

		case '"':
			l.parseString()

		case 't':
			l.parseKeyword(TTrue)

		case 'f':
			l.parseKeyword(TFalse)

		case 'n':
			l.parseKeyword(TNull)

		default:
			// Check for insignificant whitespace characters
			if IsWhitespace(l.codePoint) {
				l.step()
				continue
			}

			l.SyntaxError(fmt.Sprintf("Unexpected token %c", l.source[l.end]))
		}

		return
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

func (l *Lexer) parseNumber() {
	first := l.codePoint
	l.step()
	if first == '0' && l.codePoint == '0' {
		l.SyntaxError("Unexpected number")
	}

	// Initial digits
	for {
		if l.codePoint < '0' || l.codePoint > '9' {
			break
		}
		l.step()
	}

	// Fractional digits
	if l.codePoint == '.' {
		l.step()
		for {
			if l.codePoint < '0' || l.codePoint > '9' {
				break
			}
			l.step()
		}
	}

	// Exponent
	if l.codePoint == 'e' || l.codePoint == 'E' {
		l.step()
		if l.codePoint == '+' || l.codePoint == '-' {
			l.step()
		}
		if l.codePoint < '0' || l.codePoint > '9' {
			l.SyntaxError(fmt.Sprintf("Unexpected token %c", l.source[l.current]))
		}
		for {
			if l.codePoint < '0' || l.codePoint > '9' {
				break
			}
			l.step()
		}
	}

	l.Token = TNumber
	l.Number, _ = strconv.ParseFloat(l.source[l.start:l.end], 64)
}

func (l *Lexer) parseString() {
	for {
		l.step()
		if l.codePoint == '"' {
			l.step()
			break
		}
		if l.codePoint == -1 {
			l.SyntaxError("Unexpected end of JSON input")
		}
	}

	l.Token = TString
	l.String = l.source[l.start+1 : l.end-1]
}

func (l *Lexer) parseKeyword(t Token) {
	s := KeywordToString[t]
	width := len(s)

	for i := 1; i < width; i++ {
		c := l.source[l.current+i-1]
		if c != s[i] {
			l.SyntaxError(fmt.Sprintf("Unexpected token %c", c))
		}
	}

	l.end = l.current + width - 1
	l.current = l.end
	l.Token = t
	l.step()
}

func (l *Lexer) SyntaxError(msg string) {
	loc := Loc{Start: int32(l.end)}
	panic(LexerPanic{
		Msg: msg,
		Loc: loc,
	})
}

func NewLexer(s string) Lexer {
	lexer := Lexer{
		source: s,
	}
	lexer.step()
	lexer.Next()
	return lexer
}
