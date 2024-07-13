package serializer

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/needl3/redis-cli-lite/pkg/constants/identifier"
)

func TestSimpleStringTokenizer(t *testing.T) {
	expression := "abcd\r\n"
	tokenizer := NewSimpleStringTokenizer()
	token, expr := tokenizer.Parse([]byte(expression))
	if string(token.Value) != "abcd" || expr != nil {
		t.Fatalf("Simple string tokenizer failed to tokenize string %s", expression)
	}
}

func TestBulkStringTokenizer(t *testing.T) {
	expression := "4\r\nabcd\r\n"
	tokenizer := NewBulkStringTokenizer()
	token, expr := tokenizer.Parse([]byte(expression))
	if string(token.Value) != "abcd" || expr != nil {
		t.Fatalf("Bulk string tokenizer failed to tokenize string %s", expression)
	}
}

func TestIntegerTokenizer(t *testing.T) {
	expression := "5\r\n"
	tokenizer := NewIntegersTokenizer()
	token, expr := tokenizer.Parse([]byte(expression))
	if token.Value != 5 || expr != nil {
		t.Fatalf("Integer tokenizer failed to tokenize string negative number %s", expression)
	}
}

func TestNegativeIntegerTokenizer(t *testing.T) {
	expression := "-5\r\n"
	tokenizer := NewIntegersTokenizer()
	token, expr := tokenizer.Parse([]byte(expression))
	if token.Value != -5 || expr != nil {
		t.Fatalf("Integer tokenizer failed to tokenize string negative number %s", expression)
	}
}

func TestArrayTokenizer(t *testing.T) {
	expression := "2\r\n+asdf\r\n$2\r\nab\r\n"
	expectedResult := Token[[]Token[any]]{
		TokenType: identifier.ARRAY,
		Value: []Token[any]{{TokenType: identifier.SIMPLE_STRING, Value: TokenizerOutput("asdf")},
			{TokenType: identifier.BULK_STRING, Value: TokenizerOutput("ab")}},
	}
	tokenizer := NewArrayTokenizer()
	token, _ := tokenizer.Parse([]byte(expression))
	for idx, tkn := range token.Value {
		val := expectedResult.Value[idx].Value
		if !bytes.Equal(tkn.Value.(TokenizerOutput), val.(TokenizerOutput)) || tkn.TokenType != expectedResult.Value[idx].TokenType {
			t.Fatalf("Array tokenizer failed to tokenize array. Expected: %v\nReceived: %v", expectedResult.Value[idx].Value, val)
		}
	}
}

func Test2DArrayTokenizer(t *testing.T) {
	expression := "2\r\n+asdf\r\n*2\r\n$2\r\nab\r\n-asdf\r\n"
	expectedResult := Token[[]Token[any]]{
		TokenType: identifier.ARRAY,
		Value: []Token[any]{
			{TokenType: identifier.SIMPLE_STRING, Value: TokenizerOutput("asdf")},
			{TokenType: identifier.ARRAY, Value: []Token[any]{
				{TokenType: identifier.BULK_STRING, Value: TokenizerOutput("ab")},
				{TokenType: identifier.SIMPLE_ERROR, Value: TokenizerOutput("asdf")},
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
