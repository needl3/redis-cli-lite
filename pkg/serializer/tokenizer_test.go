package serializer

import (
	"fmt"
	"reflect"
	"testing"
)

func TestSimpleStringTokenizer(t *testing.T) {
	expression := "abcd\r\n"
	tokenizer := NewSimpleStringTokenizer()
	token, expr := tokenizer.Parse([]byte(expression))
	if token.Value != "abcd" || expr != nil {
		t.Fatalf("Simple string tokenizer failed to tokenize string %s", expression)
	}
}

func TestBulkStringTokenizer(t *testing.T) {
	expression := "4\r\nabcd\r\n"
	tokenizer := NewBulkStringTokenizer()
	token, expr := tokenizer.Parse([]byte(expression))
	if token.Value != "abcd" || expr != nil {
		t.Fatalf("Bulk string tokenizer failed to tokenize string %s", expression)
	}
}

func TestIntegerTokenizer(t *testing.T) {
	expression := "5\r\n"
	tokenizer := NewIntegersTokenizer()
	token, expr := tokenizer.Parse([]byte(expression))
	if token.Value != "5" || expr != nil {
		t.Fatalf("Integer tokenizer failed to tokenize string negative number %s", expression)
	}
}

func TestNegativeIntegerTokenizer(t *testing.T) {
	expression := "-5\r\n"
	tokenizer := NewIntegersTokenizer()
	token, expr := tokenizer.Parse([]byte(expression))
	if token.Value != "-5" || expr != nil {
		t.Fatalf("Integer tokenizer failed to tokenize string negative number %s", expression)
	}
}

func TestArrayTokenizer(t *testing.T) {
	expression := "2\r\n:5\r\n$2\r\nab\r\n"
	expectedResult := Token[[]Token[any]]{
		TokenType: "array",
		Value: []Token[any]{{TokenType: "integer", Value: "5"},
			{TokenType: "bulkstring", Value: "ab"}},
	}
	tokenizer := NewArrayTokenizer()
	token, _ := tokenizer.Parse([]byte(expression))
	for idx, tkn := range token.Value {
		if tkn.Value != expectedResult.Value[idx].Value || tkn.TokenType != expectedResult.Value[idx].TokenType {
			t.Fatalf("Array tokenizer failed to tokenize array. Expected: %v\nReceived: %v", expectedResult, token.Value)
		}
	}
}

func Test2DArrayTokenizer(t *testing.T) {
	expression := "2\r\n:5\r\n*2\r\n$2\r\nab\r\n:-1\r\n"
	expectedResult := Token[[]Token[any]]{
		TokenType: "array",
		Value: []Token[any]{
			{TokenType: "integer", Value: "5"},
			{TokenType: "array", Value: []Token[any]{
				{TokenType: "bulkstring", Value: "ab"},
				{TokenType: "integer", Value: "-1"},
			},
			},
		},
	}
	tokenizer := NewArrayTokenizer()
	token, _ := tokenizer.Parse([]byte(expression))
	if !reflect.DeepEqual(token.Value, expectedResult.Value) {
		t.Fatalf("Array tokenizer failed to tokenize array.\nExpected: %v\nReceived: %v", expectedResult, token)
	}
}
