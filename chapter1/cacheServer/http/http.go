package http

import (
	"../cache"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type Server struct {
	cache.Cache
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	res := strings.Split(r.URL.EscapedPath(), "/")[1]
	if res == "status" {
		s.ServeStatus(w, r)
		return
	}
	if res == "cache" {
		s.ServeCache(w, r)
		return
	}
	w.WriteHeader(http.StatusNotFound)
}

func (s *Server) ServeCache(w http.ResponseWriter, r *http.Request) {
	key := strings.Split(r.URL.EscapedPath(), "/")[2]
	if len(key) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	m := r.Method
	if m == http.MethodPut {
		b, _ := ioutil.ReadAll(r.Body)
		if len(b) != 0 {
			s.Set(key, b)
		} else {
			log.Println(len(b))
			w.WriteHeader(http.StatusBadRequest)
		}
		return
	}
	if m == http.MethodGet {
		b := s.Get(key)
		if len(b) != 0 {
			w.Write(b)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
		return
	}
	if m == http.MethodDelete {
		s.Del(key)
		return
	}
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func (s *Server) ServeStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	stat := s.GetStat()
	b, _ := json.Marshal(stat)
	w.Write(b)
}
