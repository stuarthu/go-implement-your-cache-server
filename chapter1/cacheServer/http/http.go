package http

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	res := strings.Split(r.URL.EscapedPath(), "/")[1]
	if res == "status" {
		h.serveStatus(w, r)
		return
	}
	if res == "cache" {
		h.serveCache(w, r)
		return
	}
	w.WriteHeader(http.StatusNotFound)
}

func (h *handler) serveCache(w http.ResponseWriter, r *http.Request) {
	key := strings.Split(r.URL.EscapedPath(), "/")[2]
	if len(key) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	m := r.Method
	if m == http.MethodPut {
		b, _ := ioutil.ReadAll(r.Body)
		if len(b) != 0 {
			h.Set(key, b)
		} else {
			log.Println(len(b))
			w.WriteHeader(http.StatusBadRequest)
		}
		return
	}
	if m == http.MethodGet {
		b := h.Get(key)
		if len(b) != 0 {
			w.Write(b)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
		return
	}
	if m == http.MethodDelete {
		h.Del(key)
		return
	}
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func (h *handler) serveStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	stat := h.GetStat()
	b, _ := json.Marshal(stat)
	w.Write(b)
}
