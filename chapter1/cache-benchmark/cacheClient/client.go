package cacheClient

type Cmd struct {
	Name  string
	Key   string
	Value string
    Error error
}

type Client interface {
	Run(*Cmd)
}
