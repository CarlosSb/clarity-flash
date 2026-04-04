package cache

import "sync"

// Cache e um cache in-memory com TTL basico
type Cache struct {
	mu    sync.RWMutex
	items map[string]*item
}

type item struct {
	value    interface{}
	expiring bool
}

func NewCache() *Cache {
	return &Cache{
		items: make(map[string]*item),
	}
}

func (c *Cache) Set(key string, value interface{}) {
	c.mu.Lock()
	c.items[key] = &item{value: value}
	c.mu.Unlock()
}

func (c *Cache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	v, ok := c.items[key]
	c.mu.RUnlock()
	if !ok {
		return nil, false
	}
	return v.value, true
}

func (c *Cache) Delete(key string) {
	c.mu.Lock()
	delete(c.items, key)
	c.mu.Unlock()
}
