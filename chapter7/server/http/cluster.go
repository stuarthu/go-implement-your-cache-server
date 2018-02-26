package http

import (
	"../cluster"
	"encoding/json"
	"net/http"
)

type clusterHandler struct {
	cluster.Node
}

func (c *clusterHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	m := c.Members()
	b, _ := json.Marshal(m)
	w.Write(b)
}

func (s *Server) clusterHandler() http.Handler {
	return &clusterHandler{s}
}
