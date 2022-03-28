package redis

import (
	"github.com/go-redis/redis/v8"
)

var r *redis.Client

func init() {
	r = redis.NewClient(redisConfig)
}
