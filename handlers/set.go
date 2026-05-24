package handlers

import (
	"strconv"
	"strings"
	"time"
)

func setCommand(args []string) string {
	if len(args) < 2 {
		return "-ERR wrong number of arguments for 'set' command\r\n"
	}

	key := args[0]
	value := args[1]

	newEntry := entry {
		key: key,
		value: value,
		expiresAt: nil,
	}

	if len(args) == 4 {
		expiryArgs := strings.ToUpper(args[2])
		timeVal,err := strconv.Atoi(args[3])
		if err == nil {
			if expiryArgs == "PX" {
				expiry := time.Now().Add(time.Duration(timeVal) * time.Millisecond)
				newEntry.expiresAt = &expiry
			} else if expiryArgs == "EX" {
				expiry := time.Now().Add(time.Duration(timeVal) * time.Second)
				newEntry.expiresAt = &expiry
			}
		}

	}
	mu.Lock()
	defer mu.Unlock()
	if listElement,exists := store[key]; exists {
		evictionList.MoveToFront(listElement)
		listElement.Value = newEntry
	}

	if evictionList.Len() >= maxCapacity {
		tail := evictionList.Back()
		if tail != nil {
			evictionList.Remove(tail)
			oldestEntry := tail.Value.(entry)
			delete(store,oldestEntry.key)
		}
	}

	newElement := evictionList.PushFront(newEntry)
	store[key] = newElement
	return "+OK\r\n"
}