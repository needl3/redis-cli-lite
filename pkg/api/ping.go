package api

import (
	"errors"

	"github.com/needl3/redis-cli-lite/pkg/constants/commands"
	"github.com/needl3/redis-cli-lite/pkg/constants/identifier"
	"github.com/needl3/redis-cli-lite/pkg/serializer"
)

func (api Api) Ping() (string, error) {
	encoded := api.Encoder.Encode(commands.PING)
	response, err := api.Lib.SendRaw(encoded)
	if err != nil {
		return "", err
	}

	token, _, err := serializer.Parse(response)
	if err != nil {
		return "", err
	}

	if token.TokenType == identifier.SIMPLE_STRING {

		if val, ok := token.Value.(serializer.TokenizerOutput); ok {
			return string(val), nil
		}
		return "", errors.New("Invalid ping response")
	}
	return "", errors.New("Invalid ping response")
}
