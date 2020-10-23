package http

import (
	"net/http"

	"github.com/stuarthu/go-implement-your-cache-server/chapter8/server/cache"
	"github.com/stuarthu/go-implement-your-cache-server/chapter8/server/cluster"
)

type Server struct {
	cache.Cache
	cluster.Node
}

func (s *Server) Listen() {
	http.Handle("/cache/", s.cacheHandler())
	http.Handle("/status", s.statusHandler())
	http.Handle("/cluster", s.clusterHandler())
	http.Handle("/rebalance", s.rebalanceHandler())
	http.ListenAndServe(s.Addr()+":12345", nil)
}

func New(c cache.Cache, n cluster.Node) *Server {
	return &Server{c, n}
}
