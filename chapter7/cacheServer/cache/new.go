package cache

import "log"

func New(typ string) Cache {
	defer log.Println(typ, "ready to serve")

	if typ == "inmemory" {
		return NewInMemoryCache()
	} else if typ == "rocksdb" {
		return NewRocksdbCache()
	}

	panic("unknown cache type " + typ)
}
