package main

import (
	"./http"
	"./tcp"
)

func main() {
	go tcp.New(ca, node).Listen()
	http.New(ca, node).Listen()
}
