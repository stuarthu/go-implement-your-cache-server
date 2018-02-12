package cacheClient

import (
	"github.com/go-redis/redis"
	"log"
)

type redisClient struct {
	client *redis.Client
}

func (r *redisClient) get(key string) string {
	res, e := r.client.Get(key).Result()
	if e == redis.Nil {
		return ""
	}
	if e != nil {
		log.Println(key)
		panic(e)
	}
	return res
}

func (r *redisClient) set(key, value string) {
	e := r.client.Set(key, value, 0).Err()
	if e != nil {
		panic(e)
	}
}

func (r *redisClient) Run(c *Cmd) {
	if c.Name == "get" {
		c.Value = r.get(c.Key)
		return
	}
	if c.Name == "set" {
		r.set(c.Key, c.Value)
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
				value = ""
			} else if e != nil {
				log.Println(c.Key)
				panic(e)
			}
			c.Value = value
		}
	}
}

func NewRedisClient(server string) *redisClient {
	client := redis.NewClient(&redis.Options{Addr: server + ":6379", ReadTimeout: -1})
	return &redisClient{client}
}
