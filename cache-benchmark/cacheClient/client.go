package cacheClient

type Cmd struct {
	Name  string
	Key   string
	Value string
	Error error
}

type Client interface {
	Run(*Cmd)
	PipelinedRun([]*Cmd)
}

func New(typ, server string) Client {
	if typ == "redis" {
		return newRedisClient(server)
	}
	if typ == "http" {
		return newHTTPClient(server)
	}
	if typ == "tcp" {
		return newTCPClient(server)
	}
	panic("unknown client type " + typ)
}
