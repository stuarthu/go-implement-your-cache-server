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
	return sendResponse(nil, s.Del(k), conn)
}
