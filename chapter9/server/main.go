package main

import (
	"github.com/stuarthu/go-implement-your-cache-server/chapter9/server/cache"
	"github.com/stuarthu/go-implement-your-cache-server/chapter9/server/cluster"
	"github.com/stuarthu/go-implement-your-cache-server/chapter9/server/http"
	"github.com/stuarthu/go-implement-your-cache-server/chapter9/server/tcp"
	"flag"
	"log"
)

func main() {
	typ := flag.String("type", "inmemory", "cache type")
	ttl := flag.Int("ttl", 30, "cache time to live")
	node := flag.String("node", "127.0.0.1", "node address")
	clus := flag.String("cluster", "", "cluster address")
	flag.Parse()
	log.Println("type is", *typ)
	log.Println("ttl is", *ttl)
	log.Println("node is", *node)
	log.Println("cluster is", *clus)
	c := cache.New(*typ, *ttl)
	n, e := cluster.New(*node, *clus)
	if e != nil {
		panic(e)
	}
	go tcp.New(c, n).Listen()
	http.New(c, n).Listen()
}
