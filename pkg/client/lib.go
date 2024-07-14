package client

import (
	"errors"
	"fmt"
	"net"
)

type Library struct {
	Host     string
	Port     int
	ConnPool []net.Conn
}

func (lib *Library) addr() string {
	return fmt.Sprintf("%s:%d", lib.Host, lib.Port)
}

func (lib *Library) Connect() error {
	return lib.ConnectPool(1)
}

func (lib *Library) ConnectPool(pool int) error {
	for i := 0; i < pool; i++ {
		maxRetries := 3

		var err error
		for ; maxRetries > 0; maxRetries-- {
			conn, _err := net.Dial("tcp", lib.addr())
			if _err == nil {
				lib.ConnPool = append(lib.ConnPool, conn)
				break
			}
			fmt.Println("[X] Couldn't connect to redis server")
			err = _err

		}
		if maxRetries == 0 {
			fmt.Println("Couldn't connect to server. Closing all connections")
			for _, c := range lib.ConnPool {
				c.Close()
			}
			return err
		}
	}
	return nil
}

func (lib *Library) getConnection() (net.Conn, error) {
	if len(lib.ConnPool) == 0 {
		return nil, errors.New("No connection available")
	}

	conn := lib.ConnPool[0]
	lib.ConnPool = lib.ConnPool[1:]
	return conn, nil
}

func (lib *Library) releaseConnection(conn net.Conn) {
	lib.ConnPool = append(lib.ConnPool, conn)
}

func (lib *Library) SendRaw(message []byte) ([]byte, error) {
	conn, err := lib.getConnection()
	if err != nil {
		return nil, err
	}

	defer lib.releaseConnection(conn)

	_, err = conn.Write(message)
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
