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
	v, e := s.Get(k)
	if e != nil {
		return replyError(e, conn)
	}
	vlen := fmt.Sprintf("%d ", len(v))
	_, e = conn.Write(append([]byte(vlen), v...))
	return e
}

func replyError(err error, conn net.Conn) error {
	if err != nil {
		errString := err.Error()
		tmp := fmt.Sprintf("-%d ", len(errString)) + errString
		_, e := conn.Write([]byte(tmp))
		return e
	}
	_, e := conn.Write([]byte("0 "))
	return e
}
