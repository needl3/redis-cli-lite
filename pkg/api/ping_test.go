package api

import (
	"log"
	"testing"
)

func TestPing(t *testing.T) {
	api, err := Initialize("localhost", 6379, 1)
	if err != nil {
		log.Fatal(err)
	}

	resp, err := api.ping()
	if err != nil {
		t.Fatal(err)
	}

	if resp != "PONG" {
		t.Fatalf("Expected PONG, got %s", resp)
	}
}
