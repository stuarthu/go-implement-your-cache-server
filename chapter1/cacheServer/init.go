package main

import (
	"./cache"
)

var ca cache.Cache

func init() {
	ca = cache.New("inmemory")
}
