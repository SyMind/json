package ast

type T interface{}

type JSONValue struct {
	Value T
}

type Object struct {
	Properties []Property
}

type Property struct {
	Name  String
	Value JSONValue
}

type Array struct {
	Items []JSONValue
}

type Number struct {
	Value float64
}

type String struct {
	Value string
}

type True struct{}

type False struct{}

type Null struct{}
