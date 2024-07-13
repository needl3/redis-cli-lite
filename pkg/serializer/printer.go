package serializer

import (
	"fmt"

	"github.com/needl3/redis-cli-lite/pkg/constants/identifier"
)

func Pretty(token Token[any]) string {
	switch token.TokenType {
	case identifier.SIMPLE_STRING:
		return string(token.Value.(TokenizerOutput))
	case identifier.SIMPLE_ERROR:
		return "ERROR: " + string(token.Value.(TokenizerOutput))
	case identifier.INTEGER:
		return fmt.Sprint(token.Value.(int))
	case identifier.BULK_STRING:
		return string(token.Value.(TokenizerOutput))
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
