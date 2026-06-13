package handlers

import (
	"fmt"
)

func lpushCommand(args []string) string {
	if len(args) < 2 {
		return "-ERR wrong number of arguments for 'lpush' command\r\n"
	}

	key := args[0]
	valuesToPush := args[1:]

	reversed := make([]string, len(valuesToPush))
	for i, j := 0, len(valuesToPush)-1; i < len(valuesToPush); i, j = i+1, j-1 {
		reversed[i] = valuesToPush[j]
	}

	var listData []string
	mu.Lock()
	defer mu.Unlock()

	if element, exists := store[key]; exists {
		entryElement := element.Value.(entry)
		existingList, ok := entryElement.value.([]string)

		if !ok {
			return "-WRONGTYPE Operation against a key holding the wrong kind of value\r\n"
		}

		listData = append(reversed, existingList...)
		
		for len(entryElement.waiters) > 0 && len(listData) > 0{
			oldestClientChan := entryElement.waiters[0]
			entryElement.waiters = entryElement.waiters[1:]

			element.Value = entryElement
			dataValue := listData[0]
			listData = listData[1:]
			oldestClientChan <- dataValue
		}
		entryElement.value = listData
		element.Value = entryElement

		evictionList.MoveToFront(element)
	} else {
		listData = reversed

		newEntry := entry{
			key:       key,
			value:     listData,
			expiresAt: nil,
		}

		if evictionList.Len() >= maxCapacity {
			oldestElement := evictionList.Back()
			if oldestElement != nil {
				evictionList.Remove(oldestElement)
				oldestEntry := oldestElement.Value.(entry)
				delete(store, oldestEntry.key)
			}
		}

		newElement := evictionList.PushFront(newEntry)
		store[key] = newElement
	}



	return fmt.Sprintf(":%d\r\n", len(listData))
}