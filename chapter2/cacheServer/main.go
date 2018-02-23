package main

import (
	myHttp "./http"
	"./tcp"
	"log"
	"net/http"
)

func main() {
	go tcp.NewServer(ca).Listen()
	log.Fatal(http.ListenAndServe(":12345", myHttp.NewHandler(ca)))
}
