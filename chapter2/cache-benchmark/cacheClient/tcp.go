package cacheClient

import (
	"bufio"
    "errors"
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

func (r *tcpClient) sendGet(key string) {
	klen := len(key)
	r.c.Write([]byte(fmt.Sprintf("G%d %s", klen, key)))
}

func (r *tcpClient) sendSet(key, value string) {
	klen := len(key)
	vlen := len(value)
	r.c.Write([]byte(fmt.Sprintf("S%d %d %s%s", klen, vlen, key, value)))
}

func (r *tcpClient) sendDel(key string) {
	klen := len(key)
	r.c.Write([]byte(fmt.Sprintf("D%d %s", klen, key)))
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

func (r *tcpClient) recvResponse() (string, error) {
	vlen := readLen(r.r)
	if vlen == 0 {
		return "", nil
	}
	if vlen < 0 {
		err := make([]byte, -vlen)
		_, e := io.ReadFull(r.r, err)
		if e != nil {
			return "", e
		}
		return "", errors.New(string(err))
	}
	value := make([]byte, vlen)
	_, e := io.ReadFull(r.r, value)
	if e != nil {
		return "", e
	}
	return string(value), nil
}

func (r *tcpClient) Run(c *Cmd) {
	if c.Name == "get" {
		r.sendGet(c.Key)
		c.Value, c.Error = r.recvResponse()
		return
	}
	if c.Name == "set" {
		r.sendSet(c.Key, c.Value)
		_, c.Error = r.recvResponse()
		return
	}
	if c.Name == "del" {
		r.sendDel(c.Key)
		_, c.Error = r.recvResponse()
		return
	}
	panic("unknown cmd name " + c.Name)
}

func newTCPClient(server string) *tcpClient {
	c, e := net.Dial("tcp", server+":12346")
	if e != nil {
		panic(e)
	}
	r := bufio.NewReader(c)
	return &tcpClient{c, r}
}
