package tcp

import (
	"bufio"
	"fmt"
	"net"
)

func (s *server) get(conn net.Conn, r *bufio.Reader) error {
	k, e := readKey(r)
	if e != nil {
		return e
	}
	v := s.Get(k)
	vlen := fmt.Sprintf("%d ", len(v))
	_, e = conn.Write(append([]byte(vlen), v...))
	return e
}
