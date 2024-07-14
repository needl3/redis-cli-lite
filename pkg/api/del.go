package api

import (
	"github.com/needl3/redis-cli-lite/pkg/constants/commands"
)

func (api Api) del(key string) error {
	encoded := api.Encoder.Encode(commands.DEL + " " + key)
	response, err := api.Lib.SendRaw(encoded)
	if err != nil {
		return err
	}

	_, _, err = api.Parser(response)
	if err != nil {
		return err
	}
	return nil
}
