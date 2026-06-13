package handlers

import (
	"container/list"
	"sync"
	"time"
)

type entry struct {
	key       string
	value     any
	expiresAt *time.Time
	waiters  []chan string
}

const maxCapacity = 10000

var (
	store        = make(map[string]*list.Element)
	evictionList = list.New()
	mu           sync.Mutex
)
