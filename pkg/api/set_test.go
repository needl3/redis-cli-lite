package api

import (
	"log"
	"testing"
)

func TestSet(t *testing.T) {
	api, err := Initialize("localhost", 6379, 1)
	if err != nil {
		log.Fatal(err)
	}
	err = api.set("name", "uwu")
	if err != nil {
		t.Fatal("Setting value failed: ", err)
	}
}

func TestGet(t *testing.T) {
	api, err := Initialize("localhost", 6379, 1)
	if err != nil {
		log.Fatal(err)
	}

	key := "name"
	val := "uwu"
	err = api.set(key, val)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := api.get(key)
	if err != nil {
		t.Fatal(err)
	}

	if val != resp {
		t.Fatalf("Expected %v, got %v", val, resp)
	}
}
