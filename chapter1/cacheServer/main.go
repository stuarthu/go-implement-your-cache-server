package main

import (
	myHttp "./http"
	"log"
	"net/http"
)

func main() {
	httpServer := &myHttp.Server{ca}
	//    http.Handle("/cache/", httpServer)
	//    http.Handle("/status/", httpServer)
	log.Fatal(http.ListenAndServe(":12345", httpServer))
}
