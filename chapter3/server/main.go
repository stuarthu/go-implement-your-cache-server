package main

import (
	"github.com/stuarthu/go-implement-your-cache-server/chapter3/server/cache"
	"github.com/stuarthu/go-implement-your-cache-server/chapter3/server/http"
	"github.com/stuarthu/go-implement-your-cache-server/chapter3/server/tcp"
	
	"flag"
	"log"
)

func main() {
	typ := flag.String("type", "inmemory", "cache type")
	flag.Parse()
	log.Println("type is", *typ)
	c := cache.New(*typ)
	go tcp.New(c).Listen()
	http.New(c).Listen()
}
