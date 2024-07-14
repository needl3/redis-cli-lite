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

type TokenizerOutput []byte

func scanToken(expr []byte, function func(byte, int)) []byte {
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
			function(val, idx)
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

func NewSimpleStringTokenizer() Tokenizer[TokenizerOutput] {
	return Tokenizer[TokenizerOutput]{
		Identifier: identifier.SIMPLE_STRING,
		Parse: func(expr []byte) (Token[TokenizerOutput], []byte) {
			token := Token[TokenizerOutput]{
				TokenType: identifier.SIMPLE_STRING,
				Value:     []byte{},
			}
			expr = scanToken(expr, func(b byte, _ int) {
				token.Value = append(token.Value, b)
			})
			return token, expr
		},
	}
}

func NewSimpleErrorTokenizer() Tokenizer[TokenizerOutput] {
	return Tokenizer[TokenizerOutput]{
		Identifier: identifier.SIMPLE_ERROR,
		Parse: func(expr []byte) (Token[TokenizerOutput], []byte) {
			token := Token[TokenizerOutput]{
				TokenType: identifier.SIMPLE_ERROR,
				Value:     []byte{},
			}
			expr = scanToken(expr, func(b byte, _idx int) {
				token.Value = append(token.Value, b)
			})
			return token, expr
		},
	}
}

func NewBulkStringTokenizer() Tokenizer[TokenizerOutput] {
	return Tokenizer[TokenizerOutput]{
		Identifier: identifier.BULK_STRING,
		Parse: func(expr []byte) (Token[TokenizerOutput], []byte) {
			token := Token[TokenizerOutput]{
				TokenType: identifier.BULK_STRING,
				Value:     nil,
			}

			size, expr := ExtractLength(expr)
			if size < 0 {
				return token, expr
			}
			token.Value = make([]byte, size)

			expr = scanToken(expr, func(b byte, idx int) {
				token.Value[idx] = b
			})
			return token, expr
		},
	}
}

func NewIntegersTokenizer() Tokenizer[int] {
	return Tokenizer[int]{
		Identifier: identifier.INTEGER,
		Parse: func(expr []byte) (Token[int], []byte) {
			sign := 1
			if expr[0] == '-' {
				expr = expr[1:]
				sign = -1
			}
			token := Token[int]{
				TokenType: identifier.INTEGER,
				Value:     0,
			}

			expr = scanToken(expr, func(b byte, _ int) {
				token.Value = token.Value*10 + int(b-'0')
			})
			token.Value *= sign

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
