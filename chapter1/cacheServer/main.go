package main

import (
	myHttp "./http"
	"log"
	"net/http"
)

func main() {
	log.Fatal(http.ListenAndServe(":12345", myHttp.NewHandler(ca)))
}
