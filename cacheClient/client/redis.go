package client

import (
	"github.com/go-redis/redis"
	"log"
)

type redisClient struct {
	client *redis.Client
}

func (r *redisClient) Get(key string) string {
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

func (r *redisClient) Set(key, value string) {
	e := r.client.Set(key, value, 0).Err()
	if e != nil {
		panic(e)
	}
}

func NewRedisClient(server string) *redisClient {
	client := redis.NewClient(&redis.Options{Addr: server, ReadTimeout: -1, DialTimeout: -1})
	return &redisClient{client}
}
