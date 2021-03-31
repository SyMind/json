package parser

import (
	"github.com/SyMind/json/ast"
	"github.com/SyMind/json/lexer"
)

type parser struct {
	source string
	lexer  lexer.Lexer
}

func (p *parser) parseJSONValue() ast.JSONValue {
	switch p.lexer.Token {
	case lexer.TLeftCurlyBracket:
		// TODO

	case lexer.TLeftSquareBracket:
		// TODO

	case lexer.TNumber:
		p.lexer.Next()
		return ast.JSONValue{Value: ast.Number{
			Value: p.lexer.Number,
		}}

	case lexer.TString:
		p.lexer.Next()
		return ast.JSONValue{Value: ast.String{
			Value: p.lexer.String,
		}}

	case lexer.TTrue:
		p.lexer.Next()
		return ast.JSONValue{Value: ast.True{}}

	case lexer.TFalse:
		p.lexer.Next()
		return ast.JSONValue{Value: ast.False{}}

	case lexer.TNull:
		p.lexer.Next()
		return ast.JSONValue{Value: ast.Null{}}
	}

	panic("Syntax error")
}

func Parse(source string) (result ast.JSONValue, ok bool) {
	ok = true

	defer func() {
		r := recover()
		if _, isLexerPanic := r.(lexer.LexerPanic); isLexerPanic {
			ok = false
		} else if r != nil {
			panic(r)
		}
	}()

	p := &parser{
		source: source,
		lexer:  lexer.NewLexer(source),
	}

	result = p.parseJSONValue()
	return
}
