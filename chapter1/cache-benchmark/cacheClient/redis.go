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

func newRedisClient(server string) *redisClient {
	client := redis.NewClient(&redis.Options{Addr: server + ":6379", ReadTimeout: -1})
	return &redisClient{client}
}
