package parser

import (
	"github.com/SyMind/json/ast"
	"github.com/SyMind/json/lexer"
)

type parser struct {
	source string
	lexer  lexer.Lexer
}

func (p *parser) parseMaybeTrailingComma(closeToken lexer.Token) {
	p.lexer.Next()
	if p.lexer.Token == lexer.TComma {
		p.lexer.Next()
	}

	if p.lexer.Token == closeToken {
		panic("JSON does not support trailing commas")
	}
}

func (p *parser) parseJSONValue() ast.JSONValue {
	switch p.lexer.Token {
	case lexer.TLeftCurlyBracket:
		p.lexer.Next()
		properties := []ast.Property{}
		for p.lexer.Token != lexer.TRightCurlyBracket {
			if len(properties) > 0 {
				p.lexer.Next()
				if p.lexer.Token == lexer.TComma {
					p.lexer.Next()
				}
			}
			nameString := p.lexer.String
			name := ast.String{Value: nameString}

			p.lexer.Next()
			if p.lexer.Token == lexer.TColon {
				p.lexer.Next()
			}
			jsonValue := p.parseJSONValue()

			property := ast.Property{
				Name:  name,
				Value: jsonValue,
			}
			properties = append(properties, property)
		}

		return ast.JSONValue{Value: ast.Object{
			Properties: properties,
		}}

	case lexer.TLeftSquareBracket:
		p.lexer.Next()
		items := []ast.JSONValue{}
		for p.lexer.Token != lexer.TRightSquareBracket {
			if len(items) > 0 {
				p.parseMaybeTrailingComma(lexer.TRightSquareBracket)
			}
			item := p.parseJSONValue()
			items = append(items, item)
		}
		p.lexer.Next()
		return ast.JSONValue{Value: ast.Array{
			Items: items,
		}}

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
