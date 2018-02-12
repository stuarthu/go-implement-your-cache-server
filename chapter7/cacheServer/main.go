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

func main() {
	l, e := net.Listen("tcp", ":12346")
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
