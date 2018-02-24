package tcp

import "bufio"

func (s *server) get(ch chan chan *result, r *bufio.Reader) error {
	k, e := readKey(r)
	if e != nil {
		return e
	}
	c := make(chan *result)
	go func() {
		v, e := s.Get(k)
		c <- &result{v, e}
	}()
	ch <- c
	return nil
}
