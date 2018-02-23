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
	return replyError(s.Del(k), conn)
}
