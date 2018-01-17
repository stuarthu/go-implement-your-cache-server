package client

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"strings"
)

type tcpClient struct {
	c net.Conn
	r *bufio.Reader
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

func (r *tcpClient) Get(key string) chan string {
	klen := len(key)
	c := fmt.Sprintf("G%d %s", klen, key)
	_, e := r.c.Write([]byte(c))
	if e != nil {
		panic(e)
	}
	ch := make(chan string)
	go func() {
		vlen := readLen(r.r)
		value := make([]byte, vlen)
		_, e = io.ReadFull(r.r, value)
		if e != nil {
			panic(e)
		}
		ch <- value
	}()
	return ch
}

func (r *tcpClient) Set(key, value string) {
	klen := len(key)
	vlen := len(value)
	c := fmt.Sprintf("S%d %d %s%s", klen, vlen, key, value)
	_, e := r.c.Write([]byte(c))
	if e != nil {
		panic(e)
	}
}

func NewTCPClient(server string) *tcpClient {
	c, e := net.Dial("tcp", server)
	if e != nil {
		panic(e)
	}
	r := bufio.NewReader(c)
	return &tcpClient{c, r}
}
