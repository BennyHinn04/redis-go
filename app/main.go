package main

import (
	"fmt"
	"net"
	"os"

	"github.com/codecrafters-io/redis-starter-go/resp"
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

	buffer := make([]byte, 1024)

	bytesRead, err := conn.Read(buffer)
	if err == nil {
		fmt.Printf("Successfull read %d bytes, buffer content %q", bytesRead,buffer)
	} else {
		fmt.Printf("Error occured : %q", err.Error())
	}

	command, bytesConsumed, err := resp.ParseSimpleString(buffer)

	if err == nil {
		fmt.Printf("Successfully parsed the command %q of %d bytes", command, bytesConsumed)
		response := "+PONG\r\n"
		conn.Write([]byte(response))
	}

	if err != nil {
		fmt.Println(err.Error())
	}

}
