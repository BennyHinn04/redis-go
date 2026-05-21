package handlers

import "fmt"

func getCommand(args []string) string {
	if len(args) < 1 {
		return "-ERR: "
	}
	
	key := args[0]
	mu.RLock()

	value,exists := store[key]

	mu.RUnlock()

	if !exists {
		return "$-1/r/n"
	}
	return fmt.Sprintf("$%d\r\n%s\r\n", len(value), value)

}