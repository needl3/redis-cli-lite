package api

import (
	"log"
	"testing"

	"github.com/needl3/redis-cli-lite/pkg/utils"
)

func TestDel(t *testing.T) {
	tlsConfig, err := utils.PrepareTLSConfig("../../redis.crt", "../../redis.key")
	if err != nil {
		t.Fatal(err)
	}

	api, err := Initialize("localhost", 6379, 1, tlsConfig)
	if err != nil {
		log.Fatal(err)
	}

	key := "todelete"
	value := "uwu"
	err = api.set(key, value)
	if err != nil {
		t.Fatal("Cannot set value: ", err)
	}

	err = api.del(key)
	if err != nil {
		t.Fatal("Cannot delete value: ", err)
	}

	val, err := api.get(key)
	if err != nil {
		t.Fatal("Cannot get value: ", err)
	}

	if val != "" {
		t.Fatalf("Expected %v, got %v", "", val)
	}

}
