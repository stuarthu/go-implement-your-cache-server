package main

import (
	"github.com/stuarthu/go-implement-your-cache-server/chapter1/server/cache"
	"github.com/stuarthu/go-implement-your-cache-server/chapter1/server/http"
)

func main() {
	c := cache.New("inmemory")
	http.New(c).Listen()
}
