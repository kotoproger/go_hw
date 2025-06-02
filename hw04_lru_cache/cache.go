package hw04lrucache

import (
	"sync"
)

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*ListItem
	keys     map[*ListItem]Key
	mutex    sync.Mutex
}

func (c *lruCache) Set(key Key, value interface{}) bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.queue.Len() == c.capacity {
		toRemove := c.queue.Back()
		c.queue.Remove(toRemove)
		keyToRemove := c.keys[toRemove]

		delete(c.keys, toRemove)
		delete(c.items, keyToRemove)
	}

	item, ok := c.items[key]
	if !ok {
		item = c.queue.PushFront(value)

		c.items[key] = item
		c.keys[item] = key

		return false
	}

	item.Value = value
	c.queue.MoveToFront(item)

	return true
}

func (c *lruCache) Get(key Key) (interface{}, bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	item, ok := c.items[key]
	if !ok {
		return nil, false
	}

	// fmt.Println(c.queue.Front(), c.queue.Back())
	// fmt.Println(key, c.queue.Len())
	c.queue.MoveToFront(item)
	// fmt.Println(c.queue.Front(), c.queue.Back())

	return item.Value, true
}

func (c *lruCache) Clear() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.queue = NewList()
	c.items = make(map[Key]*ListItem, len(c.items))
	c.keys = make(map[*ListItem]Key, len(c.items))
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
		keys:     make(map[*ListItem]Key, capacity),
	}
}
