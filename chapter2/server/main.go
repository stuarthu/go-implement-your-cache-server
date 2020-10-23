package main

import (
	"github.com/stuarthu/go-implement-your-cache-server/chapter2/server/cache"
	"github.com/stuarthu/go-implement-your-cache-server/chapter2/server/http"
	"github.com/stuarthu/go-implement-your-cache-server/chapter2/server/tcp"
)

func main() {
	ca := cache.New("inmemory")
	go tcp.New(ca).Listen()
	http.New(ca).Listen()
}
