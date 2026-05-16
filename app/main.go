package main

import (
	"fmt"
	"net"
	"strings"
	"os"
	"redis-go/resp"
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	// Uncomment the code below to pass the first stage
	//
	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}
	conn, err := l.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}

	for {

		buffer := make([]byte, 1024)

		bytesRead, err := conn.Read(buffer)
		if err == nil {
			fmt.Printf("Successfull read %d bytes, buffer content %q", bytesRead,buffer)
		} else {
			fmt.Printf("Error occured : %q", err.Error())
			break;
		}

		commands, bytesConsumed, err := resp.ParseArray(buffer)

		if err != nil {
			fmt.Println(err.Error())
		}

		fmt.Println("Successfully consumed %d", bytesConsumed)

		command := strings.ToUpper(commands[0])

		switch command {
		case "PING":
			response := "+PONG\r\n"
			conn.Write([]byte(response))
		}
	}
	

}
