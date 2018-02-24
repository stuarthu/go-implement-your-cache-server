package cacheClient

type Client interface {
	Run(*Cmd)
	PipelinedRun([]*Cmd)
}
