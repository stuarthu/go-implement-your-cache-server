package http

import (
	"../cache"
	"encoding/json"
	"net/http"
)

type statusHandler struct {
	cache.Cache
}

func (s *statusHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	stat := s.GetStat()
	b, _ := json.Marshal(stat)
	w.Write(b)
}

func (s *Server) statusHandler() http.Handler {
	return &statusHandler{s}
}
