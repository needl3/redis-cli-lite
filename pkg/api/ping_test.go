package api

import (
	"log"
	"testing"

	"github.com/needl3/redis-cli-lite/pkg/utils"
)

func TestPing(t *testing.T) {
	tlsConfig, err := utils.PrepareTLSConfig("../../redis.crt", "../../redis.key")
	if err != nil {
		t.Fatal(err)
	}

	api, err := Initialize("localhost", 6379, 1, tlsConfig)
	if err != nil {
		log.Fatal(err)
	}

	resp, err := api.Ping()
	if err != nil {
		t.Fatal(err)
	}

	if resp != "PONG" {
		t.Fatalf("Expected PONG, got %s", resp)
	}
}
