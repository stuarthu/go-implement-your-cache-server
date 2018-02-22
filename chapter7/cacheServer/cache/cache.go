package cache

type Cache interface {
	Set(string, []byte)
	Get(string) []byte
}
