package main

import (
	"flag"
	"log"

	"github.com/stuarthu/go-implement-your-cache-server/chapter7/server/cache"
	"github.com/stuarthu/go-implement-your-cache-server/chapter7/server/cluster"
	"github.com/stuarthu/go-implement-your-cache-server/chapter7/server/http"
	"github.com/stuarthu/go-implement-your-cache-server/chapter7/server/tcp"
)

func main() {
	typ := flag.String("type", "inmemory", "cache type")
	node := flag.String("node", "127.0.0.1", "node address")
	clus := flag.String("cluster", "", "cluster address")
	flag.Parse()
	log.Println("type is", *typ)
	log.Println("node is", *node)
	log.Println("cluster is", *clus)
	c := cache.New(*typ)
	n, e := cluster.New(*node, *clus)
	if e != nil {
		panic(e)
	}
	go tcp.New(c, n).Listen()
	http.New(c, n).Listen()
}
