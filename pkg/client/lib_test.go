package client_test

import (
	"context"
	"math/rand"
	"net"
	"testing"
	"time"

	"github.com/needl3/redis-cli-lite/pkg/client"
)

func TestConnectionPoolPopulation(t *testing.T) {
	const connectionPool = 10

	lib := client.Library{
		Host:     "127.0.0.1",
		Port:     6379,
		ConnPool: make(chan net.Conn, connectionPool),
	}
	lib.ConnectPool(connectionPool)
	if len(lib.ConnPool) != connectionPool {
		t.Fatalf("Connection pool not populated")
	}
}

func TestStressfulConnectionPool(t *testing.T) {
	const connectionPool = 10
	const testTime = 10 // In seconds

	lib := client.Library{
		Host:     "127.0.0.1",
		Port:     6379,
		ConnPool: make(chan net.Conn, connectionPool),
	}
	lib.ConnectPool(connectionPool)

	ctx, cancel := context.WithCancel(context.Background())

	// Print the monitor logs
	go func() {
		logChan := lib.Monitor(ctx)
		for {
			select {
			case log, ok := <-logChan:
				if ok {
					t.Log(log)
				}
			}
		}
	}()

	// Connection consumer
	go func() {
		for {
			conn, ok := <-lib.ConnPool
			if !ok {
				cancel()
				t.Log("Connection pool not populated")
			}
			go func() {
				time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
				lib.ConnPool <- conn
			}()
			time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
		}
	}()

	time.Sleep(testTime * time.Second)
	cancel()
}
