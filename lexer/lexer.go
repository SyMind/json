package lexer

import (
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

func IsWhitespace(codePoint rune) bool {
	switch codePoint {
	case
		'\u0009', // character tabulation
		'\u000B', // line tabulation
		'\u000C', // form feed
		'\u0020', // space
		'\u00A0', // no-break space

		// Unicode "Space_Separator" code points
		'\u1680', // ogham space mark
		'\u2000', // en quad
		'\u2001', // em quad
		'\u2002', // en space
		'\u2003', // em space
		'\u2004', // three-per-em space
		'\u2005', // four-per-em space
		'\u2006', // six-per-em space
		'\u2007', // figure space
		'\u2008', // punctuation space
		'\u2009', // thin space
		'\u200A', // hair space
		'\u202F', // narrow no-break space
		'\u205F', // medium mathematical space
		'\u3000', // ideographic space

		'\uFEFF': // zero width non-breaking space
		return true

	default:
		return false
	}
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
		panic(LexerPanic{})
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
			panic(LexerPanic{})
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
			panic("Unclosed string")
		}
	}

	l.Token = TString
	l.String = l.source[l.start+1 : l.end-1]
}

func (l *Lexer) parseKeyword(t Token) {
	s := KeywordToString[t]
	width := len(s)
	text := l.source[l.current-1 : l.current+width-1]

	if s != text {
		panic("Unexpect token")
	}

	l.end = l.current + width - 1
	l.current = l.end
	l.Token = t
	l.step()
}

func NewLexer(s string) Lexer {
	lexer := Lexer{
		source: s,
	}
	lexer.step()
	lexer.Next()
	return lexer
}
