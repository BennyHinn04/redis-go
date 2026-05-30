package handlers

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func lrangeCommand(args []string) string {
	if len(args) < 3 {
		return "-ERR wrong number of arguments for 'lrange' command\r\n"
	}

	key := args[0]

	startIndex, err := strconv.Atoi(args[1])
	if err != nil {
		return "-ERR value is not an integer or out of range\r\n"
	}
	endIndex, err := strconv.Atoi(args[2])
	if err != nil {
		return "-ERR value is not an integer or out of range\r\n"
	}

	mu.Lock()
	defer mu.Unlock()

	element, exists := store[key]
	if !exists {
		return "*0\r\n"
	}

	curEntry := element.Value.(entry)

	if curEntry.expiresAt != nil && time.Now().After(*curEntry.expiresAt) {
		delete(store, key)
		evictionList.Remove(element)
		return "*0\r\n"
	}

	list, ok := curEntry.value.([]string)
	if !ok {
		return "-WRONGTYPE Operation against a key holding the wrong kind of value\r\n"
	}

	listLen := len(list)

	if startIndex < 0 {
		startIndex = listLen + startIndex
	}
	if endIndex < 0 {
		endIndex = listLen + endIndex
	}

	if startIndex < 0 {
		startIndex = 0
	}
	if endIndex >= listLen {
		endIndex = listLen - 1
	}

	if startIndex > endIndex || startIndex >= listLen {
		return "*0\r\n"
	}

	slicedData := list[startIndex : endIndex+1]

	evictionList.MoveToFront(element)

	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("*%d\r\n", len(slicedData)))

	for _, item := range slicedData {
		sb.WriteString(fmt.Sprintf("$%d\r\n%s\r\n", len(item), item))
	}

	return sb.String()
}
