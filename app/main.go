package main

import (
	"fmt"
	"net"
	"os"
	"redis-go/resp"
	"redis-go/handlers"
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
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}
		go handleConnection(conn)
	}

}
func handleConnection(conn net.Conn) {
	defer conn.Close()
	for {
		buffer := make([]byte, 1024)
		bytesRead, err := conn.Read(buffer)
		if err == nil {
			fmt.Printf("Successfull read %d bytes", bytesRead)
		} else {
			fmt.Printf("Error occured : %q", err.Error())
			break
		}

		commands, bytesConsumed, err := resp.ParseArray(buffer)

		if err != nil {
			fmt.Println(err.Error())
		}

		fmt.Printf("Successfully consumed %d", bytesConsumed)

		if len(commands) > 0 {
			args := commands[1:]
			response := handlers.Execute(commands[0],args)
			conn.Write([]byte(response))
		}
		
	}
}
