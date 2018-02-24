package cacheClient

import (
	"github.com/go-redis/redis"
)

type redisClient struct {
	client *redis.Client
}

func (r *redisClient) get(key string) (string, error) {
	res, e := r.client.Get(key).Result()
	if e == redis.Nil {
		return "", nil
	}
	return res, e
}

func (r *redisClient) set(key, value string) error {
	return r.client.Set(key, value, 0).Err()
}

func (r *redisClient) del(key string) error {
	return r.client.Del(key).Err()
}

func (r *redisClient) Run(c *Cmd) {
	if c.Name == "get" {
		c.Value, c.Error = r.get(c.Key)
		return
	}
	if c.Name == "set" {
		c.Error = r.set(c.Key, c.Value)
		return
	}
	if c.Name == "del" {
		c.Error = r.del(c.Key)
		return
	}
	panic("unknown cmd name " + c.Name)
}

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

func newRedisClient(server string) *redisClient {
	client := redis.NewClient(&redis.Options{Addr: server + ":6379", ReadTimeout: -1})
	return &redisClient{client}
}
