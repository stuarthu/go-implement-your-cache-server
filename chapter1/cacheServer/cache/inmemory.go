package cache

import "sync"

type InMemoryCache struct {
	c     map[string][]byte
	mutex sync.RWMutex
	Stat
}

func (c *InMemoryCache) Set(k string, v []byte) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.c[k] = v
	c.add(k, v)
}

func (c *InMemoryCache) Get(k string) []byte {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.c[k]
}

func (c *InMemoryCache) Del(k string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	v, exist := c.c[k]
	if !exist {
		return
	}
	delete(c.c, k)
	c.del(k, v)
}

func (c *InMemoryCache) GetStat() Stat {
	return c.Stat
}

func NewInMemoryCache() *InMemoryCache {
	return &InMemoryCache{make(map[string][]byte), sync.RWMutex{}, Stat{}}
}
