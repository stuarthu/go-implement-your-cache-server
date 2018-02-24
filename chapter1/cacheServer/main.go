package main

import (
	"./cache"
	"./http"
)

func main() {
	ca := cache.New("inmemory")
	http.New(ca).Listen()
}
