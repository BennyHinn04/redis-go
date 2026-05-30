package handlers

import (
	"fmt"
	"time"
)

func llenCommand(args []string) string {

	if len(args) < 1 {
		return "-ERR wrong number of arguments for 'llen' command\r\n"
	}

	key := args[0]

	mu.Lock()
	defer mu.Unlock()

	element, exists := store[key]

	if !exists {
		return ":0\r\n"
	}

	curEntry := element.Value.(entry)

	if curEntry.expiresAt != nil && time.Now().After(*curEntry.expiresAt) {
		delete(store, key)
		evictionList.Remove(element)
		return ":0\r\n"
	}

	listData, ok := curEntry.value.([]string)

	if !ok {
		return "-WRONGTYPE Operation against a key holding the wrong kind of value\r\n"
	}

	evictionList.MoveToFront(element)
	
	return fmt.Sprintf(":%d\r\n", len(listData))
}