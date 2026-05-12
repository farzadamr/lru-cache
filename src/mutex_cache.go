package lru

import (
	"container/list"
	"sync"
)

type mutexCache struct {
	capacity int
	items    map[string]*list.Element
	ll       *list.List
	mu       sync.RWMutex
}

func NewMutexCache(capacity int) Cache {
	return &mutexCache{
		capacity: capacity,
		items:    make(map[string]*list.Element),
		ll:       list.New(),
	}
}

func (c *mutexCache) Get(key string) (interface{}, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	elem, ok := c.items[key]
	if !ok {
		return nil, false
	}
	c.ll.MoveToFront(elem)
	return elem.Value.(*entry).value, true
}

func (c *mutexCache) Set(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if elem, ok := c.items[key]; ok {
		elem.Value.(*entry).value = value
		c.ll.MoveToFront(elem)
		return
	}
	if c.ll.Len() >= c.capacity {
		c.evictOldest()
	}
	newEntry := &entry{key: key, value: value}
	elem := c.ll.PushFront(newEntry)
	c.items[key] = elem
}

func (c *mutexCache) Remove(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if elem, ok := c.items[key]; ok {
		delete(c.items, key)
		c.ll.Remove(elem)
	}
}

func (c *mutexCache) Len() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.ll.Len()
}

func (c *mutexCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items = make(map[string]*list.Element)
	c.ll.Init()
}

func (c *mutexCache) evictOldest() {
	oldest := c.ll.Back()
	if oldest != nil {
		delete(c.items, oldest.Value.(*entry).key)
		c.ll.Remove(oldest)
	}
}
