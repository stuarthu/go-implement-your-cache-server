package tcp

import (
	"../cache"
)

type server struct {
	cache.Cache
}

func NewServer(c cache.Cache) Server {
	return &server{c}
}
