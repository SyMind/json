package parser

import (
	"reflect"
	"testing"

	"github.com/SyMind/json/ast"
)

func AssertEqual(t *testing.T, a interface{}, b interface{}) {
	t.Helper()
	if a != b {
		t.Fatalf("%s != %s", a, b)
	}
}

func AssertEqualJSONValue(t *testing.T, jsonValue ast.JSONValue, expected ast.JSONValue) {
	t.Helper()
	AssertEqual(t, reflect.TypeOf(jsonValue.Value), reflect.TypeOf(expected.Value))

	switch jsonValue.Value.(type) {
	case ast.Object:
		properties := jsonValue.Value.(ast.Object).Properties
		expectedProperties := expected.Value.(ast.Object).Properties
		AssertEqual(t, len(properties), len(expectedProperties))

		for i := 0; i < len(properties); i++ {
			AssertEqual(t, properties[i].Name.Value, expectedProperties[i].Name.Value)
			AssertEqualJSONValue(t, properties[i].Value, expectedProperties[i].Value)
		}
	case ast.Array:
		items := jsonValue.Value.(ast.Array).Items
		expectedItems := expected.Value.(ast.Array).Items
		AssertEqual(t, len(items), len(expectedItems))

		for i := 0; i < len(items); i++ {
			AssertEqualJSONValue(t, items[i], expectedItems[i])
		}
	case ast.Number:
		AssertEqual(t, jsonValue.Value.(ast.Number).Value, expected.Value.(ast.Number).Value)
	case ast.String:
		AssertEqual(t, jsonValue.Value.(ast.String).Value, expected.Value.(ast.String).Value)
	case ast.True:
	case ast.False:
	case ast.Null:
	}
}

func expectJSONValue(t *testing.T, contents string, expected ast.JSONValue) {
	t.Helper()
	jsonValue, _ := Parse(contents)
	AssertEqualJSONValue(t, jsonValue, expected)
}

func TestJSONParser(t *testing.T) {
	expectJSONValue(t, "true", ast.JSONValue{Value: ast.True{}})
	expectJSONValue(t, "false", ast.JSONValue{Value: ast.False{}})
	expectJSONValue(t, "null", ast.JSONValue{Value: ast.Null{}})
	expectJSONValue(t, "\"x\"", ast.JSONValue{Value: ast.String{
		Value: "x",
	}})
	expectJSONValue(t, "123", ast.JSONValue{Value: ast.Number{
		Value: 123,
	}})
	expectJSONValue(t, "[\"a\", 123, true, false, null]", ast.JSONValue{Value: ast.Array{
		Items: []ast.JSONValue{
			{Value: ast.String{
				Value: "a",
			}},
			{Value: ast.Number{
				Value: 123.0,
			}},
			{Value: ast.True{}},
			{Value: ast.False{}},
			{Value: ast.Null{}},
		},
	}})
	expectJSONValue(t, "{\"a\": 123, \"b\": \"v\", \"c\": true, \"d\": false, \"e\": null}", ast.JSONValue{Value: ast.Object{
		Properties: []ast.Property{
			{
				Name: ast.String{Value: "a"},
				Value: ast.JSONValue{Value: ast.Number{
					Value: 123.0,
				}},
			},
			{
				Name: ast.String{Value: "b"},
				Value: ast.JSONValue{Value: ast.String{
					Value: "v",
				}},
			},
			{
				Name:  ast.String{Value: "c"},
				Value: ast.JSONValue{Value: ast.True{}},
			},
			{
				Name:  ast.String{Value: "d"},
				Value: ast.JSONValue{Value: ast.False{}},
			},
			{
				Name:  ast.String{Value: "e"},
				Value: ast.JSONValue{Value: ast.Null{}},
			},
		},
	}})
}
