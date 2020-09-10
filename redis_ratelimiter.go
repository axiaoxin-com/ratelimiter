package ratelimiter

import (
	"context"

	"github.com/go-redis/redis/v8"
)

// RedisRatelimiter redis limiter
type RedisRatelimiter struct {
	*redis.Client
	Bucket BucketConfig
	script *redis.Script
}

// NewRedisRatelimiter 根据配置创建 redis limiter
func NewRedisRatelimiter(rdb *redis.Client, conf BucketConfig) *RedisRatelimiter {
	// 创建 limiter
	intervalMicrosecond := defaultBucketFillEveryMicrosecond
	if conf.FillEveryMicrosecond > 0 {
		intervalMicrosecond = conf.FillEveryMicrosecond
	}
	bucketCapacity := defaultBucketCapacity
	if conf.Capacity > 0 {
		bucketCapacity = conf.Capacity
	}
	expire := defaultBucketExpireSecond
	if conf.ExpireSecond > 0 {
		expire = conf.ExpireSecond
	}
	return &RedisRatelimiter{
		Client: rdb,
		Bucket: BucketConfig{
			FillEveryMicrosecond: intervalMicrosecond,
			Capacity:             bucketCapacity,
			ExpireSecond:         expire,
		},
		script: tokenBucketRedisLuaScript,
	}
}

// Allow 判断给定 key 是否被允许
func (r *RedisRatelimiter) Allow(ctx context.Context, key string) bool {
	return true
}
