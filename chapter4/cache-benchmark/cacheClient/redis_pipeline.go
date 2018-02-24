package cacheClient

import (
	"github.com/go-redis/redis"
)

func (r *redisClient) PipelinedRun(cmds []*Cmd) {
	if len(cmds) == 0 {
		return
	}
	pipe := r.client.Pipeline()
	cmders := make([]redis.Cmder, len(cmds))
	for i, c := range cmds {
		if c.Name == "get" {
			cmders[i] = pipe.Get(c.Key)
		} else if c.Name == "set" {
			cmders[i] = pipe.Set(c.Key, c.Value, 0)
		} else if c.Name == "del" {
			cmders[i] = pipe.Del(c.Key)
		} else {
			panic("unknown cmd name " + c.Name)
		}
	}
	_, e := pipe.Exec()
	if e != nil && e != redis.Nil {
		panic(e)
	}
	for i, c := range cmds {
		if c.Name == "get" {
			value, e := cmders[i].(*redis.StringCmd).Result()
			if e == redis.Nil {
				value, e = "", nil
			}
			c.Value, c.Error = value, e
		} else {
			c.Error = cmders[i].Err()
		}
	}
}
