package client

import (
	"bufio"
	"fmt"
	"github.com/needl3/redis-cli-lite/pkg/searilizer"
	"io"
	"net"
	"os"
)

type Client struct {
	conn       net.Conn
	serializer *serializer.Searializer
}

func (cli Client) HandleConnection() {
	defer cli.conn.Close()

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("127.0.0.1:6379>")
		cmd, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				fmt.Println("\nBye Bye")
				return
			}

			fmt.Println("\nError reading input")
			continue
		}

		_, err = cli.conn.Write(cli.serializer.Encoder.Encode(cmd))
		if err != nil {
			fmt.Println("Error sending data to server: ", err.Error())
			continue
		}

		rawResponse := make([]byte, 1024)
		_, err = cli.conn.Read(rawResponse)
		if err != nil {
			fmt.Println("Error reading server response")
			fmt.Println(err.Error())
		}
		parsedResponse := cli.serializer.Parser.Parse(rawResponse)

		fmt.Printf(parsedResponse)
	}
}

func New(host string, port string) Client {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", host, port))
	if err != nil {
		fmt.Println("Couldn't connect to redis server")
		fmt.Println(err.Error())
		os.Exit(1)
	}

	return Client{
		conn:       conn,
		serializer: serializer.New(),
	}
}
