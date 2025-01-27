package cache

import (
	"context"

	"github.com/redis/go-redis/v9"
)

var (
	redisClient *redis.Client
)

func init() {
	redisClient = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
}

func CacheRedisSet(key string, value string) error {
	return redisClient.Set(context.Background(), key, value, 0).Err()
}

func CacheRedisGet(key string) (string, error) {
	return redisClient.Get(context.Background(), key).Result()
}
