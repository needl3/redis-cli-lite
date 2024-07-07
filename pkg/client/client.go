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
	conn       net.Conn
	serializer *serializer.Searializer
	host       string
	port       string
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

		encoded := cli.serializer.Encoder.Encode(strings.Trim(cmd, "\n"))
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
		parsedResponse, _, err := cli.serializer.Parser.Parse(rawResponse)
		if err != nil {
			fmt.Println("[X] Invalid response. Failed to parse")
			fmt.Println(err.Error())
			return
		}

		fmt.Println(parsedResponse.Value)
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
		conn:       conn,
		serializer: serializer.New(),
		host:       host,
		port:       port,
	}
}
