package api

import (
	"errors"

	"github.com/needl3/redis-cli-lite/pkg/constants/commands"
	"github.com/needl3/redis-cli-lite/pkg/constants/identifier"
	"github.com/needl3/redis-cli-lite/pkg/serializer"
)

func (api Api) set(key string, val string) error {
	encoded := api.Encoder.Encode(commands.SET + " " + key + " " + val)
	response, err := api.Lib.SendRaw(encoded)
	if err != nil {
		return err
	}

	token, _, err := api.Parser(response)
	if err != nil {
		return err
	}

	if token.TokenType == identifier.SIMPLE_STRING {
		if _, ok := token.Value.(serializer.TokenizerOutput); ok {
			return nil
		}
		return errors.New("Invalid set response")
	}
	return errors.New("Invalid set response")
}
