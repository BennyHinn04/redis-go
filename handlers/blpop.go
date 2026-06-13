// Here the blpop command waits for value to be inserted if the list is empty
// Since we lock the store before reading, other client can't inserted in the list since lock is held by the 
// client running blpop
// So the goroutine, running blpop must communicate with goroutine adding the value to the list
// Also we need to ensure the value inserted must be popped by client that is waiting for long(comparitively)
// So we will maintian a wait queue (for pushing values to the earliest client) with channels to specific entry in store 
// So other clients can send value to the client running blpop
package handlers

import (
	"fmt"
	"strconv"
	"time"
)

func blpopCommand(args []string) string {
	if len(args) < 2 {
		return "-ERR wrong number of arguments for 'blpop' command\r\n"
	}

	key := args[0]
	timeout := args[1]

	mu.Lock()

	element,exists := store[key]

	if exists {
		curEntry := element.Value.(entry)
		listData,ok := curEntry.value.([]string)
		if ok && len(listData) > 0 {
			poppedValue := listData[0]
			curEntry.value = listData[1:]
			element.Value = curEntry

			evictionList.MoveToFront(element)
			
			mu.Unlock()

			return fmt.Sprintf("$%d\r\n%s\r\n", len(poppedValue), poppedValue)
		} 

	} else {
		newEntry := entry{
			key	: key,
			value : []string{},
			waiters: []chan string{}
		}
		store[key] = evictionList.PushFront(newEntry)
	}

	newChannel := make(chan string)
	element = store[key]
	curEntry := element.Value.(entry)
	curEntry.waiters = append(curEntry.waiters,newChannel)

	element.Value = curEntry

	mu.Unlock()
	// Goroutine sleeps until the newChannel has a value
	timeoutSec,err := strconv.Atoi(timeout)
	if timeoutSec == 0 {
		poppedValue := <-newChannel
		return fmt.Sprintf("$%d\r\n%s\r\n", len(poppedValue), poppedValue)
	} else {
		select {
		case poppedValue := <-newChannel:
			return fmt.Sprintf("$%d\r\n%s\r\n", len(poppedValue), poppedValue)
		case <-time.After(time.Duration(timeoutSec) * time.Second):
			return "$-1\r\n"
		}
	}


}