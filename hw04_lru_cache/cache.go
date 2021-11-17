package hw04lrucache

import "sync"

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
	mu       *sync.Mutex
}

func (c *lruCache) Clear() {
	c.items = make(map[Key]*ListItem)
	c.queue = &list{}
	c.capacity = 0
}

func (c *lruCache) Get(key Key) (interface{}, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	item, found := c.items[key]

	if found {
		c.queue.MoveToFront(item)
		return item.Value.(cacheItem).value, found
	}
	return nil, found
}

func (c *lruCache) Set(key Key, value interface{}) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	item, found := c.items[key]

	cacheItemStruct := cacheItem{
		key:   string(key),
		value: value,
	}

	if found {
		item.Value = cacheItemStruct
		c.queue.MoveToFront(item)
	} else {
		newItem := c.queue.PushFront(cacheItemStruct)
		c.items[key] = newItem
		if len(c.items) > c.capacity {
			el := c.queue.Back()
			key := el.Value.(cacheItem).key
			delete(c.items, Key(key))
			c.queue.Remove(c.queue.Back())
		}
	}
	return found
}

type cacheItem struct {
	key   string
	value interface{}
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
		mu:       &sync.Mutex{},
	}
}
