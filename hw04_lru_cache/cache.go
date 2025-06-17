package hw04lrucache

import "sync"

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type CacheItem struct {
	key   Key
	value interface{}
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*ListItem
	mu       sync.Mutex
}

func (c *lruCache) Set(key Key, value interface{}) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	item, exists := c.items[key]
	if exists {
		item.Value.(*CacheItem).value = value
		c.queue.MoveToFront(item)
		return true
	}
	cacheItem := &CacheItem{key: key, value: value}
	listItem := c.queue.PushFront(cacheItem)
	c.items[key] = listItem

	if c.queue.Len() > c.capacity {
		itemBack := c.queue.Back()
		if itemBack != nil {
			delete(c.items, itemBack.Value.(*CacheItem).key)
			c.queue.Remove(itemBack)
		}
	}

	return false
}

func (c *lruCache) Get(key Key) (interface{}, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	item, exists := c.items[key]
	if exists {
		c.queue.MoveToFront(item)
		return item.Value.(*CacheItem).value, true
	}
	return nil, false
}

func (c *lruCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.queue = NewList()
	c.items = make(map[Key]*ListItem, c.capacity)
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}
