package handlers

import (
	"fmt"
	"time"
)

func getCommand(args []string) string {
	if len(args) < 1 {
		return "-ERR: "
	}
	
	key := args[0]

	mu.Lock()
	defer mu.Unlock()

	storeElement,exists := store[key]
	if !exists {
		return "$-1/r/n"
	}

	curEntry := storeElement.Value.(entry)
	if curEntry.expiresAt != nil && time.Now().After(*curEntry.expiresAt) {
		delete(store,curEntry.key)
		evictionList.Remove(storeElement)
		return "$-1\r\n"
	}

	evictionList.MoveToFront(storeElement)

	return fmt.Sprintf("$%d\r\n%s\r\n", len(curEntry.value), curEntry.value)

}