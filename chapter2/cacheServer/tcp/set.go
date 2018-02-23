package tcp

import (
	"bufio"
	"net"
)

func (s *server) set(conn net.Conn, r *bufio.Reader) error {
	k, v, e := readKeyAndValue(r)
	if e != nil {
		return e
	}
	s.Set(k, v)
	_, e = conn.Write([]byte("0 "))
	return e
}
