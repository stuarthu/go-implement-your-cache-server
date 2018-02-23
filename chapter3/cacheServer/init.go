package main

import (
	"./cache"
	"flag"
	"log"
)

var ca cache.Cache

func init() {
	typ := flag.String("type", "inmemory", "cache type")
	flag.Parse()
	log.Println("type is", *typ)
	ca = cache.New(*typ)
}
