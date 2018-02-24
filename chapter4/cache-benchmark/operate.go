package main

import (
	"./cacheClient"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

func pipeline(client cacheClient.Client, cmds []*cacheClient.Cmd, r *result) {
	expect := make([]string, len(cmds))
	for i, c := range cmds {
		if c.Name == "get" {
			expect[i] = c.Value
		}
	}
	start := time.Now()
	client.PipelinedRun(cmds)
	d := time.Now().Sub(start)
	for i, c := range cmds {
		resultType := c.Name
		if resultType == "get" {
			if c.Value == "" {
				resultType = "miss"
			} else if c.Value != expect[i] {
				fmt.Println(expect[i])
				panic(c.Value)
			}
		}
		r.addDuration(d, resultType)
	}
}

func operate(id, count int, ch chan *result) {
	client := cacheClient.NewCacheClient(typ, server)
	cmds := make([]*cacheClient.Cmd, 0)
	valuePrefix := strings.Repeat("a", valueSize)
	r := &result{0, 0, 0, make([]statistic, 0)}
	for i := 0; i < count; i++ {
		var tmp int
		if keyspacelen > 0 {
			tmp = rand.Intn(keyspacelen)
		} else {
			tmp = id*count + i
		}
		key := fmt.Sprintf("%d", tmp)
		value := fmt.Sprintf("%s%d", valuePrefix, tmp)
		name := operation
		if operation == "mixed" {
			if rand.Intn(2) == 1 {
				name = "set"
			} else {
				name = "get"
			}
		}
		c := &cacheClient.Cmd{name, key, value, nil}
		if pipelen > 1 {
			cmds = append(cmds, c)
			if len(cmds) == pipelen {
				pipeline(client, cmds, r)
				cmds = make([]*cacheClient.Cmd, 0)
			}
		} else {
			run(client, c, r)
		}
	}
	if len(cmds) != 0 {
		pipeline(client, cmds, r)
	}
	ch <- r
}
