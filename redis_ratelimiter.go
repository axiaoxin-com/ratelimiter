package ratelimiter

import "github.com/go-redis/redis/v8"

// RedisRatelimiter redis limiter
type RedisRatelimiter struct {
	rdb *redis.Client
}

// NewRedisRatelimiter 根据配置创建 redis limiter
func NewRedisRatelimiter() *RedisRatelimiter {

	return nil
}
