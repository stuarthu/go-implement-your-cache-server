package main

import (
	"./cache"
	"./cluster"
	"flag"
	"fmt"
)

var ca cache.Cache
var node cluster.Node

func init() {
	var typ, ip, cl string
	flag.StringVar(&typ, "type", "memory", "cache type")
	flag.StringVar(&ip, "node", "", "node address")
	flag.StringVar(&cl, "cluster", "", "cluster address")
	flag.Parse()
	fmt.Println("type is", typ)
	fmt.Println("node is", ip)
	fmt.Println("cluster is", cl)
	ca = cache.New(typ)
	var e error
	node, e = cluster.New(ip, cl)
	if e != nil {
		panic(e)
	}
}
