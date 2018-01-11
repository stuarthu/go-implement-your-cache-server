package cache

import "sync"

type InMemoryCache struct {
	c     map[string][]byte
	mutex sync.RWMutex
}

func (c *InMemoryCache) set(s string, b []byte) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.c[s] = b
}

func (c *InMemoryCache) get(s string) []byte {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.c[s]
}

func NewInMemoryCache() *InMemoryCache {
	return &InMemoryCache{make(map[string][]byte), sync.RWMutex{}}
}
