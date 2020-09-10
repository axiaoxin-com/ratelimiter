package ratelimiter

import (
	"context"

	"github.com/axiaoxin-com/logging"
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
		script: tokenBucketRedisLuaIsLimitedScript,
	}
}

// Allow 判断给定 key 是否被允许
func (r *RedisRatelimiter) Allow(ctx context.Context, key string) bool {
	// 构造 lua 脚本参数
	keys := []string{key}
	args := []interface{}{
		r.Bucket.Capacity,
		1, // lua 脚本支持调整每次放入 token 的个数，这里全部统一使用每次放一个 token
		r.Bucket.FillEveryMicrosecond,
		r.Bucket.ExpireSecond,
	}
	// 在 redis 中执行 lua 脚本计算当前 key 是否被限频
	// Run 会自动使用 evalsha 优化带宽
	v, err := r.script.Run(ctx, r.Client, keys, args...).Result()
	if err != nil {
		// 有 err 默认放行
		logging.Error(ctx, "RedisRatelimiter run script error:"+err.Error())
		return true
	}
	isLimited, ok := v.(int64)
	if !ok {
		logging.Error(ctx, "RedisRatelimiter assert script result error")
		return true
	}
	// 1 表示被限频，返回 false 不允许本次操作
	if isLimited == 1 {
		return false
	}
	return true
}
