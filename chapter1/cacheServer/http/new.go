package http

import (
	"../cache"
	"net/http"
)

type handler struct {
	cache.Cache
}

func NewHandler(c cache.Cache) http.Handler {
	return &handler{c}
}
