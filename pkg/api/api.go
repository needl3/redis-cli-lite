package api

//
// Api is used to interact with redis server
// TODO apis
// rpush, rpop, ltrim, lindex, lset, lrem
//

import (
	"crypto/tls"
	"net"

	"github.com/needl3/redis-cli-lite/pkg/client"
	"github.com/needl3/redis-cli-lite/pkg/serializer"
)

type Api struct {
	Lib     *client.Library
	Encoder serializer.Encoder
	Parser  func(expr []byte) (serializer.Token[any], []byte, error)
}

func Initialize(host string, port int, pool int, tlsConfig *tls.Config) (*Api, error) {
	lib := &client.Library{
		Host:     host,
		Port:     port,
		ConnPool: make(chan net.Conn, pool),
	}
	var err error
	err = lib.ConnectPool(pool, tlsConfig)

	if err != nil {
		return nil, err
	}
	return &Api{
		Lib:     lib,
		Encoder: serializer.EncoderClient,
		Parser:  serializer.Parse,
	}, nil
}
