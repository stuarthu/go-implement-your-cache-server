package client

type Client interface {
	Set(string, string)
	Get(string) string
}

func NewCacheClient(typ, server string) Client {
	if typ == "http" {
		return NewHTTPClient(server)
	}
	if typ == "redis" {
		return NewRedisClient(server)
	}
	if typ == "tcp" {
		return NewTCPClient(server)
	}
	panic("invalid type " + typ)
}
