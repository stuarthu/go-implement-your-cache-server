package main

import (
	"./cache"
	"./http"
)

func main() {
	c := cache.New("inmemory")
	http.New(c).Listen()
}
