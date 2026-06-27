package handlers

import (
	"time"
)

func typeCommand(args []string) string {
	if len(args) < 1 {
		return "-ERR wrong number of arguments for 'type' command\r\n"
	}

	key := args[0]
	mu.Lock()
	defer mu.Unlock()
	element,exists := store[key]
	if !exists {
		return "+none\r\n"
	}

	curEntry := element.Value.(entry)

	if curEntry.expiresAt != nil && time.Now().After(*curEntry.expiresAt) {
		delete(store,key)
		evictionList.Remove(element)
		return "+none\r\n"
	}

	switch curEntry.value.(type) {
	case string:
		return "+string\r\n"
	case []string:
		return "+list\r\n"
	default:
		return "+unknown\r\n"
	}
}