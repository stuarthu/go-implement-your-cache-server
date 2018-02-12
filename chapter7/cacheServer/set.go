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

func set(r *bufio.Reader) error {
	klen := readLen(r)
	vlen := readLen(r)
	k := make([]byte, klen)
	v := make([]byte, vlen)
	_, e := io.ReadFull(r, k)
	if e != nil {
		log.Println(e)
		return e
	}
	_, e = io.ReadFull(r, v)
	if e != nil {
		log.Println(e)
		return e
	}
	key := string(k)
	if !node.ShouldProcess(key) {
		_, e = c.Write([]byte("6 reject"))
	} else {
		cache.Set(key, v)
		_, e = c.Write([]byte("0 "))
	}
	return e
}
