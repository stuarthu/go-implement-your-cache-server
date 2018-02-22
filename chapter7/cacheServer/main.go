package main

import (
	"bufio"
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

func reply(conn net.Conn, replyCh chan chan []byte, closeCh chan bool) {
	for {
		select {
		case r := <-replyCh:
			conn.Write(<-r)
		case <-closeCh:
			conn.Close()
			return
		}
	}
}

func process(conn net.Conn) {
	r := bufio.NewReader(conn)
	replyCh := make(chan chan []byte, 5000)
	closeCh := make(chan bool)
	defer func() {
		closeCh <- true
	}()
	go reply(conn, replyCh, closeCh)
	for {
		b, e := r.ReadByte()
		if e != nil {
			if e != io.EOF {
				log.Println("close connection due to error:", e)
			}
			return
		}
		if b == 'S' {
			e = set(conn, r)
		} else if b == 'G' {
			e = get(replyCh, r)
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
