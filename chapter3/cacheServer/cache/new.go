package cache

import "log"

func New(typ string) Cache {
	var c Cache
	if typ == "inmemory" {
		c = NewInMemoryCache()
	}
	if typ == "rocksdb" {
		c = NewRocksdbCache()
	}
	if c == nil {
		panic("unknown cache type " + typ)
	}
	log.Println(typ, "ready to serve")
	return c
}
