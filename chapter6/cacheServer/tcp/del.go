package tcp

import "bufio"

func (s *server) del(ch chan chan *result, r *bufio.Reader) error {
	k, e := readKey(r)
	if e != nil {
		return e
	}
	c := make(chan *result)
	go func() {
		c <- &result{nil, s.Del(k)}
	}()
	ch <- c
	return nil
}
