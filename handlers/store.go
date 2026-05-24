package handlers

import (
	"container/list"
	"sync"
	"time"
)

type entry struct {
	key       string
	value     string
	expiresAt *time.Time
}

const maxCapacity = 10000

var (
	store        = make(map[string]*list.Element)
	evictionList = list.New()
	mu           sync.Mutex
)
