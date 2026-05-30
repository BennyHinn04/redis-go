package handlers

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func lpopCommand(args []string) string {
	
	if len(args) < 1 || len(args) > 2 {
		return "-ERR wrong number of arguments for 'lpop' command\r\n"
	}

	key := args[0]
	hasCountArg := false
	count := 1 

	if len(args) == 2 {
		hasCountArg = true
		cleanStr := strings.TrimSpace(args[1])
		parsedCount, err := strconv.Atoi(cleanStr)
		
		if err != nil || parsedCount < 0 {
			return "-ERR value is out of range, must be positive\r\n"
		}
		count = parsedCount
	}

	mu.Lock()
	defer mu.Unlock()

	element, exists := store[key]
	if !exists {
		return "$-1\r\n"
	}

	curEntry := element.Value.(entry)

	if curEntry.expiresAt != nil && time.Now().After(*curEntry.expiresAt) {
		delete(store, key)
		evictionList.Remove(element)
		return "$-1\r\n"
	}

	list, ok := curEntry.value.([]string)
	if !ok {
		return "-WRONGTYPE Operation against a key holding the wrong kind of value\r\n"
	}

	if len(list) == 0 {
		delete(store, key)
		evictionList.Remove(element)
		return "$-1\r\n"
	}

	popCount := count
	if popCount > len(list) {
		popCount = len(list)
	}

	poppedValues := list[:popCount]
	remainingList := list[popCount:]

	if len(remainingList) == 0 {
		delete(store, key)
		evictionList.Remove(element)
	} else {
		curEntry.value = remainingList
		element.Value = curEntry
		evictionList.MoveToFront(element)
	}

	if !hasCountArg {
		singleValue := poppedValues[0]
		return fmt.Sprintf("$%d\r\n%s\r\n", len(singleValue), singleValue)
	} else {
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("*%d\r\n", len(poppedValues)))
		
		for _, val := range poppedValues {
			sb.WriteString(fmt.Sprintf("$%d\r\n%s\r\n", len(val), val))
		}
		return sb.String()
	}
}