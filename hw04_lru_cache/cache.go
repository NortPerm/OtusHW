package hw04lrucache

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
}

type cacheItem struct {
	key   string
	value interface{}
}

func (c *lruCache) Set(key Key, value interface{}) bool {
	val, ok := c.items[key]
	if ok {
		val.Value = value
		return true
	}
	c.queue.PushFront(value)
	c.items[key] = c.queue.Front()
	return false
}

func (c *lruCache) Get(key Key) (interface{}, bool) {
	val, ok := c.items[key]
	if ok {
		c.queue.MoveToFront(val)
		return val.Value, true
	}
	return nil, false
}

func (c *lruCache) Clear() {
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
