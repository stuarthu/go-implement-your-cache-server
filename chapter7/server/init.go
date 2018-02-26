package main

import (
	"./cache"
	"./cluster"
	"flag"
	"log"
)

var ca cache.Cache
var node cluster.Node

func init() {
	var typ, ip, cl string
	flag.StringVar(&typ, "type", "inmemory", "cache type")
	flag.StringVar(&ip, "node", "", "node address")
	flag.StringVar(&cl, "cluster", "", "cluster address")
	flag.Parse()
	log.Println("type is", typ)
	log.Println("node is", ip)
	log.Println("cluster is", cl)
	ca = cache.New(typ)
	var e error
	node, e = cluster.New(ip, cl)
	if e != nil {
		panic(e)
	}
}
