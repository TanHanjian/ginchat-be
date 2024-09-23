package myredis

import (
	"github.com/go-redis/redis/v8"
)

var Client *redis.Client

func InitRedis() {
	Client = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
}
