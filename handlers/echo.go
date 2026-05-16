package handlers

import (
	"fmt"
)

func echoHandler(args []string) (string) {
	if len(args) < 1 {
		return "-ERR wrong number of arguments for ECHO"
	}
	message := args[0]
	return fmt.Sprintf("$%d\r\n%s\r\n", len(message), message)
}