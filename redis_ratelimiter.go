package ratelimiter

import (
	"context"
	"time"

	"github.com/axiaoxin-com/logging"
	"github.com/go-redis/redis/v8"
)

var (
	// RedisRatelimiterCacheExpiration redis ratelimiter 缓存过期时间
	RedisRatelimiterCacheExpiration = time.Minute * 60
)

// RedisRatelimiter redis limiter
type RedisRatelimiter struct {
	*redis.Client
	script *redis.Script
}

// NewRedisRatelimiter 根据配置创建 redis limiter
func NewRedisRatelimiter(rdb *redis.Client) *RedisRatelimiter {
	return &RedisRatelimiter{
		Client: rdb,
		script: tokenBucketRedisLuaIsLimitedScript,
	}
}

// Allow 判断给定 key 是否被允许
func (r *RedisRatelimiter) Allow(ctx context.Context, key string, tokenFillInterval time.Duration, bucketSize int) bool {
	// 构造 lua 脚本参数
	keys := []string{key}
	args := []interface{}{
		bucketSize,
		1, // lua 脚本支持调整每次放入 token 的个数，这里全部统一使用每次放一个 token
		tokenFillInterval.Microseconds(),
		RedisRatelimiterCacheExpiration.Seconds(),
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
