package http

import (
	"bytes"
	"net/http"
)

type rebalanceHandler struct {
	*Server
}

func (h *rebalanceHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
	go h.rebalance()
}

func (h *rebalanceHandler) rebalance() {
	s := h.NewScanner()
	c := &http.Client{}
	for s.Scan() {
		k := s.Key()
		n, ok := h.ShouldProcess(k)
		if !ok {
			r, _ := http.NewRequest(http.MethodPut, "http://"+n+":12345/cache/"+k, bytes.NewReader(s.Value()))
			c.Do(r)
			h.Del(k)
		}
	}
	s.Close()
}

func (s *Server) rebalanceHandler() http.Handler {
	return &rebalanceHandler{s}
}
