package tcp

import (
	"bufio"
	"io"
	"log"
	"net"
)

type result struct {
	v []byte
	e error
}

func reply(conn net.Conn, resultCh chan chan *result, closeCh chan bool) {
	for {
		select {
		case c := <-resultCh:
			r := <-c
			sendResponse(r.v, r.e, conn)
		case <-closeCh:
			conn.Close()
			return
		}
	}
}

func (s *server) process(conn net.Conn) {
	r := bufio.NewReader(conn)
	resultCh := make(chan chan *result, 5000)
	closeCh := make(chan bool)
	defer func() {
		closeCh <- true
	}()
	go reply(conn, resultCh, closeCh)
	for {
		op, e := r.ReadByte()
		if e != nil {
			if e != io.EOF {
				log.Println("close connection due to error:", e)
			}
			return
		}
		if op == 'S' {
			e = s.set(resultCh, r)
		} else if op == 'G' {
			e = s.get(resultCh, r)
		} else if op == 'D' {
			e = s.del(resultCh, r)
		} else {
			log.Println("close connection due to invalid operation:", op)
			return
		}
		if e != nil {
			log.Println("close connection due to error:", e)
			return
		}
	}
}
