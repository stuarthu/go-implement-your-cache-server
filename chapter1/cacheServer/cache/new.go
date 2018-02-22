package cache

import "log"

func New(typ string) Cache {
	defer log.Println(typ, "ready to serve")

	if typ == "inmemory" {
		return NewInMemoryCache()
	}

	panic("unknown cache type " + typ)
}
