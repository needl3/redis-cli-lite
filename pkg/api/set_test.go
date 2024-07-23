package api

import (
	"log"
	"testing"

	"github.com/needl3/redis-cli-lite/pkg/utils"
)

func TestSet(t *testing.T) {
	tlsConfig, err := utils.PrepareTLSConfig("../../redis.crt", "../../redis.key")
	if err != nil {
		t.Fatal(err)
	}

	api, err := Initialize("localhost", 6379, 1, tlsConfig)
	if err != nil {
		log.Fatal(err)
	}
	err = api.Set("name", "uwu")
	if err != nil {
		t.Fatal("Setting value failed: ", err)
	}
}

func TestGet(t *testing.T) {
	tlsConfig, err := utils.PrepareTLSConfig("../../redis.crt", "../../redis.key")
	if err != nil {
		t.Fatal(err)
	}

	api, err := Initialize("localhost", 6379, 1, tlsConfig)
	if err != nil {
		log.Fatal(err)
	}

	key := "name"
	val := "uwu"
	err = api.Set(key, val)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := api.Get(key)
	if err != nil {
		t.Fatal(err)
	}

	if val != resp {
		t.Fatalf("Expected %v, got %v", val, resp)
	}
}
