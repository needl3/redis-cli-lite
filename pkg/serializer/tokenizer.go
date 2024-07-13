package serializer

import (
	"fmt"
	"strconv"

	"github.com/needl3/redis-cli-lite/pkg/constants/identifier"
)

type Token[T any] struct {
	TokenType byte
	Value     T
}

type Tokenizer[T any] struct {
	Identifier byte
	Parse      func(expr []byte) (Token[T], []byte)
}

func scanToken(expr []byte, function func(byte)) []byte {
	var delimeter []byte
	for idx, val := range expr {
		if val == '\r' {
			if len(delimeter) == 0 {
				delimeter = append(delimeter, val)
			}
		} else if val == '\n' {
			if len(delimeter) != 0 {
				// We have reached token termination point
				if len(expr) == idx+1 {
					expr = nil
				} else {
					expr = expr[idx+1:]
				}
				break
			}
		} else {
			delimeter = nil
			function(val)
		}
	}
	return expr
}

func ExtractLength(expr []byte) (int, []byte) {
	var parsedInt []byte
	var delimeter []byte
	for idx, val := range expr {
		if val == '\r' {
			if len(delimeter) == 0 {
				delimeter = append(delimeter, val)
			}
		} else if val == '\n' {
			if len(delimeter) != 0 {
				// We have reached token termination point
				length, error := strconv.Atoi(string(parsedInt))
				if error != nil {
					fmt.Println(error)
					return 0, expr
				}
				return length, expr[idx+1:]
			}
		} else {
			delimeter = nil
			parsedInt = append(parsedInt, val)
		}
	}
	return 0, expr
}

func NewSimpleStringTokenizer() Tokenizer[string] {
	return Tokenizer[string]{
		Identifier: identifier.SIMPLE_STRING,
		Parse: func(expr []byte) (Token[string], []byte) {
			token := Token[string]{
				TokenType: identifier.SIMPLE_STRING,
				Value:     "",
			}
			expr = scanToken(expr, func(b byte) {
				token.Value += string(b)
			})
			return token, expr
		},
	}
}

func NewSimpleErrorTokenizer() Tokenizer[string] {
	return Tokenizer[string]{
		Identifier: identifier.SIMPLE_ERROR,
		Parse: func(expr []byte) (Token[string], []byte) {
			token := Token[string]{
				TokenType: identifier.SIMPLE_STRING,
				Value:     "",
			}
			expr = scanToken(expr, func(b byte) {
				token.Value += string(b)
			})
			return token, expr
		},
	}
}

func NewBulkStringTokenizer() Tokenizer[string] {
	return Tokenizer[string]{
		Identifier: identifier.BULK_STRING,
		Parse: func(expr []byte) (Token[string], []byte) {
			token := Token[string]{
				TokenType: identifier.BULK_STRING,
				Value:     "",
			}

			_, expr = ExtractLength(expr)
			expr = scanToken(expr, func(b byte) {
				token.Value += string(b)
			})
			return token, expr
		},
	}
}

func NewIntegersTokenizer() Tokenizer[string] {
	return Tokenizer[string]{
		Identifier: identifier.INTEGER,
		Parse: func(expr []byte) (Token[string], []byte) {
			sign := ""
			if expr[0] == '-' {
				expr = expr[1:]
				sign = "-"
			}
			token := Token[string]{
				TokenType: identifier.INTEGER,
				Value:     sign,
			}

			expr = scanToken(expr, func(b byte) {
				token.Value += string(b)
			})

			return token, expr
		},
	}
}

func NewArrayTokenizer() Tokenizer[[]Token[any]] {
	return Tokenizer[[]Token[any]]{
		Identifier: identifier.ARRAY,
		Parse: func(expr []byte) (Token[[]Token[any]], []byte) {
			token := Token[[]Token[any]]{TokenType: identifier.ARRAY, Value: nil}
			arrLength, expr := ExtractLength(expr)
			token.Value = make([]Token[any], arrLength)

			var err error
			var tk Token[any]
			for i := 0; i < arrLength; i++ {
				tk, expr, err = Parse(expr)
				if err != nil {
					fmt.Println(err)
					continue
				}
				token.Value[i] = Token[any]{TokenType: tk.TokenType, Value: tk.Value}
			}
			return token, expr
		},
	}
}
