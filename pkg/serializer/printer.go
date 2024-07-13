package serializer

import (
	"github.com/needl3/redis-cli-lite/pkg/constants/identifier"
)

func Pretty(token Token[any]) string {
	switch token.TokenType {
	case identifier.SIMPLE_STRING:
		return token.Value.(string)
	case identifier.SIMPLE_ERROR:
		return "ERROR: " + token.Value.(string)
	case identifier.INTEGER:
		return token.Value.(string)
	case identifier.BULK_STRING:
		return token.Value.(string)
	case identifier.ARRAY:
		finalString := "["
		arr, ok := token.Value.([]Token[any])
		if !ok {
			return "[]"
		}
		for idx, val := range arr {
			finalString += Pretty(val)
			if idx < len(arr)-1 {
				finalString += ", "
			}
		}
		return finalString + "]"
	default:
		return ""
	}
}
