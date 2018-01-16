package cacheClient

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

func (r *tcpClient) sendGet(key string) {
	klen := len(key)
	c := fmt.Sprintf("G%d %s", klen, key)
	_, e := r.c.Write([]byte(c))
	if e != nil {
		panic(e)
	}
}

func (r *tcpClient) recvValue() string {
	vlen := readLen(r.r)
	value := make([]byte, vlen)
    _, e := io.ReadFull(r.r, value)
	if e != nil {
		panic(e)
	}
	return string(value)
}

func (r *tcpClient) sendSet(key, value string) {
	klen := len(key)
	vlen := len(value)
	c := fmt.Sprintf("S%d %d %s%s", klen, vlen, key, value)
	_, e := r.c.Write([]byte(c))
	if e != nil {
		panic(e)
	}
}

func (r *tcpClient) Run(c *Cmd) {
    if c.Name == "get" {
        r.sendGet(c.Key)
        c.Value = r.recvValue()
        return
    }
    if c.Name == "set" {
        r.sendSet(c.Key, c.Value)
        return
    }
    panic("unknown cmd name " + c.Name)
}

func (r *tcpClient) PipelinedRun(cmds []*Cmd) {
    if len(cmds) == 0 {
        return
    }
    for _, c := range cmds {
        if c.Name == "get" {
            r.sendGet(c.Key)
        }
        if c.Name == "set" {
            r.sendSet(c.Key, c.Value)
        }
    }
    for _, c := range cmds {
        if c.Name == "get" {
            c.Value = r.recvValue()
        }
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
