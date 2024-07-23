package client

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"time"
)

type Library struct {
	Host     string
	Port     int
	ConnPool chan net.Conn
}

func (lib *Library) addr() string {
	return fmt.Sprintf("%s:%d", lib.Host, lib.Port)
}

func (lib *Library) connect(secure *tls.Config) (net.Conn, error) {
	if secure == nil {
		return net.Dial("tcp", lib.addr())
	}

	return tls.Dial("tcp", lib.addr(), secure)

}

func (lib *Library) ConnectPool(pool int, secure *tls.Config) error {
	for i := 0; i < pool; i++ {
		maxRetries := 3

		var err error
		for ; maxRetries > 0; maxRetries-- {
			var conn net.Conn
			conn, err = lib.connect(secure)
			if err == nil {
				lib.ConnPool <- conn
				break
			}
		}
		if maxRetries == 0 {
			err = errors.New("Couldn't connect to server. Closing all connections")
			for {
				select {
				case conn := <-lib.ConnPool:
					conn.Close()
				default:
					return err
				}
			}
		}
	}
	return nil
}

func (lib *Library) releaseConnection(conn net.Conn) {
	lib.ConnPool <- conn
}

func (lib *Library) SendRaw(message []byte) ([]byte, error) {
	conn, ok := <-lib.ConnPool
	if !ok {
		return nil, errors.New("No connections available. Pool closed.")
	}

	defer lib.releaseConnection(conn)

	_, err := conn.Write(message)
	if err != nil {
		return nil, err
	}

	rawResponse := make([]byte, 1024)
	_, err = conn.Read(rawResponse)
	if err != nil {
		return nil, err
	}
	return rawResponse, nil
}

// Currently being used on testing
func (lib *Library) Monitor(ctx context.Context) <-chan string {
	logChan := make(chan string)
	go func() {
		for {
			logChan <- fmt.Sprintf("Connections: %d", len(lib.ConnPool))

			select {
			case <-ctx.Done():
				logChan <- fmt.Sprintln("[!] Exiting monitor")
				return
			default:
				time.Sleep(50 * time.Millisecond)
				continue
			}
		}
	}()
	return logChan
}
