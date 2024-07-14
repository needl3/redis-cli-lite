package api

import (
	"log"
	"testing"
)

var api *Api

func init() {
	_api, err := Initialize("localhost", 6379, 1)
	if err != nil {
		log.Fatal(err)
	}
	api = _api
}

func TestPing(t *testing.T) {
	resp, err := api.ping()
	if err != nil {
		t.Fatal(err)
	}

	if resp != "PONG" {
		t.Fatalf("Expected PONG, got %s", resp)
	}
}
