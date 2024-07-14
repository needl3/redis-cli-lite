package client

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strings"

	"github.com/needl3/redis-cli-lite/pkg/serializer"
)

type Client struct {
	conn    net.Conn
	encoder serializer.Encoder
	parser  func(expr []byte) (serializer.Token[any], []byte, error)
	printer func(token serializer.Token[any]) string
	host    string
	port    string
}

func (cli Client) HandleConnection() {
	defer cli.conn.Close()

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("%s:%s>", cli.host, cli.port)
		cmd, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				fmt.Println("\n[-] Bye Bye")
				return
			}

			fmt.Println("\n[X] Error reading input")
			continue
		}

		encoded := cli.encoder.Encode(strings.Trim(cmd, "\n"))
		_, err = cli.conn.Write(encoded)
		if err != nil {
			fmt.Println("[X] Error sending data to server: ", err.Error())
			continue
		}

		rawResponse := make([]byte, 1024)
		_, err = cli.conn.Read(rawResponse)
		if err != nil {
			fmt.Println("[X] Error reading server response")
			fmt.Println(err.Error())
		}
		parsedResponse, _, err := cli.parser(rawResponse)
		if err != nil {
			fmt.Println("[X] Invalid response. Failed to parse")
			fmt.Println(err.Error())
			return
		}

		fmt.Println(cli.printer(parsedResponse))
	}
}

func New(host string, port string) Client {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", host, port))
	if err != nil {
		fmt.Println("[X] Couldn't connect to redis server")
		fmt.Println(err.Error())
		os.Exit(1)
	}

	return Client{
		conn:    conn,
		host:    host,
		port:    port,
		encoder: serializer.EncoderClient,
		parser:  serializer.Parse,
		printer: serializer.Pretty,
	}
}
