package tcp

import (
	"bufio"
	"net"
)

func (s *server) del(conn net.Conn, r *bufio.Reader) error {
	k, e := readKey(r)
	if e != nil {
		return e
	}
	s.Del(k)
	_, e = conn.Write([]byte("0 "))
	return e
}
