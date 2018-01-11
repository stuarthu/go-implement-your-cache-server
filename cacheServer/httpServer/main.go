package main

import (
	"../cache"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func handler(w http.ResponseWriter, r *http.Request) {
	key := strings.Split(r.URL.EscapedPath(), "/")[1]
	if len(key) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	m := r.Method
	if m == http.MethodPut {
		b, _ := ioutil.ReadAll(r.Body)
		if len(b) != 0 {
			cache.Set(key, b)
		} else {
			log.Println(len(b))
			w.WriteHeader(http.StatusBadRequest)
		}
		return
	}
	if m == http.MethodGet {
		b := cache.Get(key)
		if len(b) != 0 {
			w.Write(b)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
		return
	}
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func main() {
	log.Fatal(http.ListenAndServe(":12345", http.HandlerFunc(handler)))
}
