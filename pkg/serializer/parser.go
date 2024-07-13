// Parser is used to parse RESP specific response given back by the redis server
// Into client specific representation

package serializer

import (
	"errors"

	"github.com/needl3/redis-cli-lite/pkg/constants/identifier"
)

func Parse(expr []byte) (Token[any], []byte, error) {
	if len(expr) < 3 {
		return Token[any]{}, nil, errors.New("Invalid expression to parse")
	}

	_identifier := expr[0]
	if _identifier == identifier.SIMPLE_STRING {
		tokenizer := NewSimpleStringTokenizer()
		token, expr := tokenizer.Parse(expr[1:])
		return Token[any]{Value: token.Value, TokenType: token.TokenType}, expr, nil
	} else if _identifier == identifier.SIMPLE_ERROR {
		tokenizer := NewSimpleErrorTokenizer()
		token, expr := tokenizer.Parse(expr[1:])
		return Token[any]{Value: token.Value, TokenType: token.TokenType}, expr, nil
	} else if _identifier == identifier.BULK_STRING {
		tokenizer := NewBulkStringTokenizer()
		token, expr := tokenizer.Parse(expr[1:])
		return Token[any]{Value: token.Value, TokenType: token.TokenType}, expr, nil
	} else if _identifier == identifier.INTEGER {
		tokenizer := NewIntegersTokenizer()
		token, expr := tokenizer.Parse(expr[1:])
		return Token[any]{Value: token.Value, TokenType: token.TokenType}, expr, nil
	} else if _identifier == identifier.ARRAY {
		tokenizer := NewArrayTokenizer()
		token, expr := tokenizer.Parse(expr[1:])
		return Token[any]{Value: token.Value, TokenType: token.TokenType}, expr, nil
	}

	return Token[any]{}, nil, errors.ErrUnsupported
}
