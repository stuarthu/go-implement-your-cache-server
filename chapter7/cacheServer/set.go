package main

import (
	"bufio"
	"io"
	"log"
	"net"
)

func set(conn net.Conn, r *bufio.Reader) error {
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
		_, e = conn.Write([]byte("6 reject"))
	} else {
		ca.Set(key, v)
		_, e = conn.Write([]byte("0 "))
	}
	return e
}
