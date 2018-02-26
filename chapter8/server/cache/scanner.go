package cache

type Scanner interface {
	Scan() bool
	Key() string
	Value() []byte
	Close()
}
