package api

import (
	"errors"
	"fmt"

	"github.com/needl3/redis-cli-lite/pkg/constants/commands"
	"github.com/needl3/redis-cli-lite/pkg/constants/identifier"
	"github.com/needl3/redis-cli-lite/pkg/serializer"
)

// LPUSH - Prepend one or multiple values to a list
// Datatype - Integer as operation success status
func (api Api) lpush(key string, value string) (int, error) {

	encoded := api.Encoder.Encode(commands.LPUSH + " " + key + " " + value)
	response, err := api.Lib.SendRaw(encoded)
	if err != nil {
		return -1, err
	}
	token, _, err := api.Parser(response)
	if err != nil {
		return -1, err
	}

	if token.TokenType == identifier.INTEGER {
		return token.Value.(int), nil
	}

	if token.TokenType == identifier.SIMPLE_ERROR {
		return -1, errors.New(string(token.Value.(serializer.TokenizerOutput)))
	}

	return -1, errors.New("Invalid lpush response")
}

// LPOP - Removes and returns the first element of the list stored at key.
// Datatype - Bulk string
func (api Api) lpop(key string) (string, error) {
	encoded := api.Encoder.Encode(commands.LPOP + " " + key)
	response, err := api.Lib.SendRaw(encoded)
	if err != nil {
		return "", err
	}
	token, _, err := api.Parser(response)
	if err != nil {
		return "", err
	}

	if token.TokenType == identifier.SIMPLE_ERROR {
		return "", errors.New(string(token.Value.(serializer.TokenizerOutput)))
	}

	if token.TokenType == identifier.BULK_STRING {
		return string(token.Value.(serializer.TokenizerOutput)), nil
	}

	return "", errors.New("Invalid response")
}

func (api Api) lrange(key string, start int, stop int) ([]any, error) {
	encoded := api.Encoder.Encode(commands.LRANGE + " " + key + " " + fmt.Sprint(start) + " " + fmt.Sprint(stop))
	response, err := api.Lib.SendRaw(encoded)
	if err != nil {
		return nil, err
	}

	token, _, err := api.Parser(response)
	if err != nil {
		return nil, err
	}

	arr, err := parseArray(token)
	if err != nil {
		return nil, err
	}

	if val, ok := arr.([]any); ok {
		return val, nil
	}

	return nil, errors.New("Invalid lrange response")
}

// No including integer or types other than bulk string and error because lists can't have other types other than string
// https://redis.io/docs/latest/develop/data-types/#lists
// Returntype doesn't have string because it's used recursively to parse array.
// So, future token type could be array since nested array is theoretically possible
func parseArray(token serializer.Token[any]) (any, error) {
	switch token.TokenType {
	case identifier.SIMPLE_STRING:
		return string(token.Value.(serializer.TokenizerOutput)), nil
	case identifier.SIMPLE_ERROR:
		return nil, errors.New(string(token.Value.(serializer.TokenizerOutput)))
	case identifier.BULK_STRING:
		return string(token.Value.(serializer.TokenizerOutput)), nil
	case identifier.ARRAY:
		finalArray := []any{}
		arr, ok := token.Value.([]serializer.Token[any])
		if !ok {
			return nil, errors.New("Invalid array value")
		}
		for _, val := range arr {
			val, err := parseArray(val)
			if err != nil {
				return nil, err
			}
			finalArray = append(finalArray, val)
		}
		return finalArray, nil
	default:
		return nil, errors.New("Invalid token type")
	}
}
