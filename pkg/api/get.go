package api

import (
	"errors"

	"github.com/needl3/redis-cli-lite/pkg/constants/commands"
	"github.com/needl3/redis-cli-lite/pkg/constants/identifier"
	"github.com/needl3/redis-cli-lite/pkg/serializer"
)

func (api Api) get(data string) (string, error) {
	encoded := api.Encoder.Encode(commands.GET + " " + data)
	response, err := api.Lib.SendRaw(encoded)
	if err != nil {
		return "", err
	}

	token, _, err := api.Parser(response)
	if err != nil {
		return "", err
	}

	if token.TokenType == identifier.SIMPLE_STRING || token.TokenType == identifier.BULK_STRING {
		if val, ok := token.Value.(serializer.TokenizerOutput); ok {
			return string(val), nil
		}
		return "", errors.New("Invalid get response")
	}
	if token.TokenType == identifier.SIMPLE_ERROR {
		if val, ok := token.Value.(serializer.TokenizerOutput); ok {
			return string(val), errors.New(string(val))
		}
	}
	return "", errors.New("Invalid get response")
}
