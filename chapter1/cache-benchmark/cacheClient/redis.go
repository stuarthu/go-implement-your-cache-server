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

func newRedisClient(server string) *redisClient {
	client := redis.NewClient(&redis.Options{Addr: server + ":6379", ReadTimeout: -1})
	return &redisClient{client}
}
