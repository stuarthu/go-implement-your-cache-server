package tcp

import "bufio"

func (s *server) set(ch chan chan *result, r *bufio.Reader) error {
	k, v, e := readKeyAndValue(r)
	if e != nil {
		return e
	}
	c := make(chan *result)
	go func() {
		c <- &result{nil, s.Set(k, v)}
	}()
	ch <- c
	return nil
}
