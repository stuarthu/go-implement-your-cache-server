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
	return sendResponse(nil, s.Set(k, v), conn)
}
