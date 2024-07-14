package api

import (
	"log"
	"testing"
)

func TestLPush(t *testing.T) {
	api, err := Initialize("localhost", 6379, 1)
	if err != nil {
		log.Fatal(err)
	}

	key := "arr"
	val := "1"
	_, err = api.lpush(key, val)
	if err != nil {
		t.Fatal(err)
	}
}

func TestLPop(t *testing.T) {
	api, err := Initialize("localhost", 6379, 1)
	if err != nil {
		log.Fatal(err)
	}

	key := "arr"
	val := "1011"
	_, err = api.lpush(key, val)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := api.lpop(key)
	if err != nil {
		t.Fatal(err)
	}

	if resp != val {
		t.Fatalf("Expected %v, got %v", val, resp)
	}
}

func TestLRange(t *testing.T) {
	api, err := Initialize("localhost", 6379, 1)
	if err != nil {
		log.Fatal(err)
	}

	key := "arr"
	start := 0
	end := 2
	vals := []string{"sdfdf", "afewl", "fj890"}
	for _, val := range vals {
		_, err = api.lpush(key, val)
		if err != nil {
			t.Fatal(err)
		}
	}

	resp, err := api.lrange(key, start, end)
	if err != nil {
		t.Fatal(err)
	}

	for idx, val := range resp {
		if val != vals[len(vals)-1-idx] {
			t.Fatalf("Expected %v, got %v", val, resp[0])
		}
	}
}
