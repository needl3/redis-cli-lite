package client

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
)

type Client struct {
	conn net.Conn
}

func (cli Client) searilizeIn(cmd string) []byte {
	searilized := fmt.Sprintf("*1\r\n$4\r\n%s\r\n", cmd)
	return []byte(searilized)
}

func (cli Client) parseResponse(res []byte) string {
	return strings.Trim(string(res), "+")
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

		_, err = cli.conn.Write(cli.searilizeIn(cmd))
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
		parsedResponse := cli.parseResponse(rawResponse)

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
		conn: conn,
	}
}
