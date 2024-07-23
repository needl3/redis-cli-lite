package api

import (
	"log"
	"testing"

	"github.com/needl3/redis-cli-lite/pkg/utils"
)

func TestLPush(t *testing.T) {
	tlsConfig, err := utils.PrepareTLSConfig("../../redis.crt", "../../redis.key")
	if err != nil {
		t.Fatal(err)
	}

	api, err := Initialize("localhost", 6379, 1, tlsConfig)
	if err != nil {
		log.Fatal(err)
	}

	key := "arr"
	val := "1"
	_, err = api.Lpush(key, val)
	if err != nil {
		t.Fatal(err)
	}
}

func TestLPop(t *testing.T) {
	tlsConfig, err := utils.PrepareTLSConfig("../../redis.crt", "../../redis.key")
	if err != nil {
		t.Fatal(err)
	}

	api, err := Initialize("localhost", 6379, 1, tlsConfig)
	if err != nil {
		log.Fatal(err)
	}

	key := "arr"
	val := "1011"
	_, err = api.Lpush(key, val)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := api.Lpop(key)
	if err != nil {
		t.Fatal(err)
	}

	if resp != val {
		t.Fatalf("Expected %v, got %v", val, resp)
	}
}

func TestLRange(t *testing.T) {
	tlsConfig, err := utils.PrepareTLSConfig("../../redis.crt", "../../redis.key")
	if err != nil {
		t.Fatal(err)
	}

	api, err := Initialize("localhost", 6379, 1, tlsConfig)
	if err != nil {
		log.Fatal(err)
	}

	key := "arr"
	start := 0
	end := 2
	vals := []string{"sdfdf", "afewl", "fj890"}
	for _, val := range vals {
		_, err = api.Lpush(key, val)
		if err != nil {
			t.Fatal(err)
		}
	}

	resp, err := api.Lrange(key, start, end)
	if err != nil {
		t.Fatal(err)
	}

	for idx, val := range resp {
		if val != vals[len(vals)-1-idx] {
			t.Fatalf("Expected %v, got %v", val, resp[0])
		}
	}
}
