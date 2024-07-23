package utils

import (
	"crypto/tls"
	"crypto/x509"
	"os"
)

func PrepareTLSConfig(certfile string, keyfile string) (*tls.Config, error) {
	cert, err := tls.LoadX509KeyPair(certfile, keyfile)
	if err != nil {
		return nil, err
	}
	caCert, err := os.ReadFile(certfile)
	if err != nil {
		return nil, err
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)
	return &tls.Config{
		RootCAs:      caCertPool,
		Certificates: []tls.Certificate{cert},
	}, nil
}
