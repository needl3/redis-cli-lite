// Parser is used to parse RESP specific response given back by the redis server
// Into client specific representation

package serializer

import (
	"errors"
)

type Parser struct {
	SimpleStringIdentifier byte
	SimpleErrorIdentifier  byte
	IntegersIdentifier     byte
	BulkStringIdentifier   byte
	ArrayIdentifier        byte
}

var ParserClient Parser

func init() {
	ParserClient = Parser{
		SimpleStringIdentifier: '+',
		SimpleErrorIdentifier:  '-',
		IntegersIdentifier:     ':',
		BulkStringIdentifier:   '$',
		ArrayIdentifier:        '*',
	}
}

func (prsr Parser) Parse(expr []byte) (Token[any], []byte, error) {
	if len(expr) < 3 {
		return Token[any]{}, nil, errors.New("Invalid expression to parse")
	}

	identifier := expr[0]
	if identifier == prsr.SimpleStringIdentifier || identifier == prsr.SimpleErrorIdentifier {
		tokenizer := NewSimpleStringTokenizer()
		token, expr := tokenizer.Parse(expr[1:])
		return Token[any]{Value: token.Value, TokenType: token.TokenType}, expr, nil
	} else if identifier == prsr.BulkStringIdentifier {
		tokenizer := NewBulkStringTokenizer()
		token, expr := tokenizer.Parse(expr[1:])
		return Token[any]{Value: token.Value, TokenType: token.TokenType}, expr, nil
	} else if identifier == prsr.IntegersIdentifier {
		tokenizer := NewIntegersTokenizer()
		token, expr := tokenizer.Parse(expr[1:])
		return Token[any]{Value: token.Value, TokenType: token.TokenType}, expr, nil
	} else if identifier == prsr.ArrayIdentifier {
		tokenizer := NewArrayTokenizer()
		token, expr := tokenizer.Parse(expr[1:])
		return Token[any]{Value: token.Value, TokenType: token.TokenType}, expr, nil
	}

	return Token[any]{}, nil, errors.ErrUnsupported
}

func (prsr Parser) Pretty(token Token[any]) string {
	switch token.TokenType {
	case identifier.SIMPLE_STRING:
		return token.Value.(string)
	case identifier.SIMPLE_ERROR:
		return token.Value.(string)
	case identifier.INTEGER:
		return string(token.Value.(int))
	case identifier.ARRAY:
		finalString := "["
		arr, ok := token.Value.([]Token[any])
		if !ok {
			return "[]"
		}
		for idx, val := range arr {
			finalString += prsr.Pretty(val)
			if idx < len(arr)-1 {
				finalString += ", "
			}
		}
		return finalString + "]"
	default:
		return ""
	}
}
