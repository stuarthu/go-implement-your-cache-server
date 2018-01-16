package main

import (
	"../cache"
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"strings"
)

func main() {
	l, e := net.Listen("tcp", ":12345")
	if e != nil {
		panic(e)
	}
	for {
		c, e := l.Accept()
		if e != nil {
			panic(e)
		}
		go process(c)
	}
}

func reply(c net.Conn, ch chan chan []byte, cl chan bool) {
	for {
		select {
		case r := <-ch:
			c.Write(<-r)
		case <-cl:
			return
		}
	}
}

func process(c net.Conn) {
	defer c.Close()
	r := bufio.NewReader(c)
	ch := make(chan chan []byte, 5000)
	cl := make(chan bool)
	defer func() {
		cl <- true
	}()
	go reply(c, ch, cl)
	for {
		b, e := r.ReadByte()
		if e != nil {
			if e != io.EOF {
				log.Println("close connection due to error:", e)
			}
			return
		}
		if b == 'S' {
			e = set(r)
		} else if b == 'G' {
			e = get(ch, r)
		} else {
			log.Println("close connection due to invalid operation:", b)
			return
		}
		if e != nil {
			log.Println("close connection due to error:", e)
			return
		}
	}
}

func readLen(r *bufio.Reader) int {
	tmp, e := r.ReadString(' ')
	if e != nil {
		log.Println(e)
		return 0
	}
	l, e := strconv.Atoi(strings.TrimSpace(tmp))
	if e != nil {
		log.Println(tmp, e)
		return 0
	}
	return l
}

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
	go cache.Set(string(k), v)
	return nil
}

func get(ch chan chan []byte, r *bufio.Reader) error {
	klen := readLen(r)
	k := make([]byte, klen)
	_, e := io.ReadFull(r, k)
	if e != nil {
		log.Println(e)
		return e
	}
    c := make(chan []byte)
	go func() {
		v := cache.Get(string(k))
		b := []byte(fmt.Sprintf("%d ", len(v)) + string(v))
		c <- b
	}()
    ch <- c
	return nil
}
