package cacheClient

type Cmd struct {
	Name  string
	Key   string
	Value string
}

type Client interface {
	Run(*Cmd)
	PipelinedRun([]*Cmd)
}

func NewCacheClient(typ, server string) Client {
	if typ == "redis" {
		return NewRedisClient(server)
	}
	if typ == "tcp" {
		return NewTCPClient(server)
	}
	panic("invalid type " + typ)
}
