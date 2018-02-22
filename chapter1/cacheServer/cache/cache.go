package cache

type Stat struct {
	Count     int64
	KeySize   int64
	ValueSize int64
}

type Cache interface {
	Set(string, []byte)
	Get(string) []byte
	Del(string)
	GetStat() Stat
}

func (s *Stat) add(k string, v []byte) {
	s.Count += 1
	s.KeySize += int64(len(k))
	s.ValueSize += int64(len(v))
}

func (s *Stat) del(k string, v []byte) {
	s.Count -= 1
	s.KeySize -= int64(len(k))
	s.ValueSize -= int64(len(v))
}
