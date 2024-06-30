package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:6379")
	if err != nil {
		fmt.Println("Couldn't connect to redis server")
		fmt.Println(err.Error())
	}

	handleConnection(conn)
}

func handleConnection(con net.Conn) {
	defer con.Close()

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

		_, err = con.Write(searilizeIn(cmd))
		if err != nil {
			fmt.Println("Error sending data to server: ", err.Error())
			continue
		}

		rawResponse := make([]byte, 1024)
		_, err = con.Read(rawResponse)
		if err != nil {
			fmt.Println("Error reading server response")
			fmt.Println(err.Error())
		}
		parsedResponse := parseResponse(rawResponse)

		fmt.Printf(parsedResponse)
	}
}

func searilizeIn(cmd string) []byte {
	searilized := fmt.Sprintf("*1\r\n$4\r\n%s\r\n", cmd)
	return []byte(searilized)
}

func parseResponse(res []byte) string {
	return strings.Trim(string(res), "+")
}
