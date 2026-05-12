package lru

import "container/list"

type command int

const (
	cmdGet command = iota
	cmdSet
	cmdRemove
	cmdLen
	cmdClear
)

type request struct {
	cmd        command
	key        string
	value      interface{}
	respGetVal chan interface{}
	respGetOk  chan bool
	respLen    chan int
}

type actorCache struct {
	capacity int
	items    map[string]*list.Element
	ll       *list.List
	reqChan  chan request
	stopChan chan struct{}
}

func NewActorCache(capacity int) Cache {
	c := &actorCache{
		capacity: capacity,
		items:    make(map[string]*list.Element),
		ll:       list.New(),
		reqChan:  make(chan request, 100),
		stopChan: make(chan struct{}),
	}
	go c.loop()
	return c
}

func (c *actorCache) loop() {
	for {
		select {
		case <-c.stopChan:
			return
		case req := <-c.reqChan:
			switch req.cmd {
			case cmdGet:
				if elem, ok := c.items[req.key]; !ok {
					req.respGetVal <- nil
					req.respGetOk <- false
				} else {
					c.ll.MoveToFront(elem)
					req.respGetVal <- elem.Value.(*entry).value
					req.respGetOk <- true
				}
			case cmdSet:
				if elem, ok := c.items[req.key]; ok {
					elem.Value.(*entry).value = req.value
					c.ll.MoveToFront(elem)
					continue
				}
				if c.ll.Len() >= c.capacity {
					c.evictOldest()
				}
				newEntry := &entry{key: req.key, value: req.value}
				elem := c.ll.PushFront(newEntry)
				c.items[req.key] = elem
			case cmdRemove:
				if elem, ok := c.items[req.key]; ok {
					delete(c.items, req.key)
					c.ll.Remove(elem)
				}
			case cmdLen:
				req.respLen <- c.ll.Len()
			case cmdClear:
				c.items = make(map[string]*list.Element)
				c.ll.Init()
			}
		}
	}
}

func (c *actorCache) Get(key string) (interface{}, bool) {
	respVal := make(chan interface{}, 1)
	respOk := make(chan bool, 1)
	c.reqChan <- request{
		cmd:        cmdGet,
		key:        key,
		respGetVal: respVal,
		respGetOk:  respOk,
	}
	return <-respVal, <-respOk
}

func (c *actorCache) Set(key string, value interface{}) {
	c.reqChan <- request{
		cmd:   cmdSet,
		key:   key,
		value: value,
	}
}

func (c *actorCache) Remove(key string) {
	c.reqChan <- request{
		cmd: cmdRemove,
		key: key,
	}
}

func (c *actorCache) Len() int {
	resp := make(chan int, 1)
	c.reqChan <- request{cmd: cmdLen, respLen: resp}
	return <-resp
}

func (c *actorCache) Clear() {
	c.reqChan <- request{cmd: cmdClear}
}

func (c *actorCache) evictOldest() {
	oldest := c.ll.Back()
	if oldest != nil {
		delete(c.items, oldest.Value.(*entry).key)
	}
}

func (c *actorCache) Close() {
	close(c.stopChan)
}
