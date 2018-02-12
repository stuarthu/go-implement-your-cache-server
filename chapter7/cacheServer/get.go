package main

import (
	"./cache"
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"strings"
)

func get(ch chan chan []byte, r *bufio.Reader) error {
	klen := readLen(r)
	k := make([]byte, klen)
	_, e := io.ReadFull(r, k)
	if e != nil {
		log.Println(e)
		return e
	}
	key := string(k)
	c := make(chan []byte)
	if !node.ShouldProcess(key) {
		c <- []byte("6 reject")
	} else {
		go func() {
			v := cache.Get(key)
			b := []byte(fmt.Sprintf("0 %d ", len(v)) + string(v))
			c <- b
		}()
	}
	ch <- c
	return nil
}
