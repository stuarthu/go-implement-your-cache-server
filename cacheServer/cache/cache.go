package cache

import (
	"log"
	"os"
)

type Cache interface {
	set(string, []byte)
	get(string) []byte
}

func newCache(typ string) Cache {
	if typ == "inmemory" {
		return NewInMemoryCache()
	} else if typ == "rocksdb" {
		return NewRocksdbCache()
	} else if typ == "rocksdb_simple" {
		return NewSimpleRocksdbCache()
	}
	panic("unknown cache type " + typ)
}

var c Cache

func init() {
	if len(os.Args) != 2 {
		panic("please specify cache type")
	}
	c = newCache(os.Args[1])
	log.Println(os.Args[1], "ready to serve")
}

func Set(key string, value []byte) {
	c.set(key, value)
}

func Get(key string) []byte {
	return c.get(key)
}
