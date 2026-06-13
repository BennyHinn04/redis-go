package handlers

import (
	"fmt"
)

func rpushCommand(args []string) string {
	if len(args) < 2 {
		return "-ERR wrong number of arguments for 'rpush' command\r\n"
	}

	key := args[0]
	valuesToPush := args[1:] 
	
	var listData []string
	var finalLength int

	mu.Lock()
	defer mu.Unlock()

	if element, exists := store[key]; exists {
		
		entryElement := element.Value.(entry)
		existingList, ok := entryElement.value.([]string)
		
		if !ok {
			return "-WRONGTYPE Operation against a key holding the wrong kind of value\r\n"
		}

		listData = append(existingList, valuesToPush...)
		entryElement.value = listData
		finalLength = len(listData)
		for len(entryElement.waiters) > 0 && len(listData) > 0{
			oldestClientChan := entryElement.waiters[0]
			entryElement.waiters = entryElement.waiters[1:]

			element.Value = entryElement
			dataValue := listData[0]
			listData = listData[1:]
			oldestClientChan <- dataValue
		}

		element.Value = entryElement 
		evictionList.MoveToFront(element)
	} else {
		listData = valuesToPush

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

	return fmt.Sprintf(":%d\r\n", finalLength)
}