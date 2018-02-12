package cache

import "sync"

type InMemoryCache struct {
	c     map[string][]byte
	mutex sync.RWMutex
}

func (c *InMemoryCache) Set(k string, v []byte) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.c[k] = v
}

func (c *InMemoryCache) Get(k string) []byte {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.c[k]
}

func NewInMemoryCache() *InMemoryCache {
	return &InMemoryCache{make(map[string][]byte), sync.RWMutex{}}
}

func (c *InMemoryCache) Del(k string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	delete(c.c, k)
}
