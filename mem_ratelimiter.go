package ratelimiter

import (
	"context"
	"time"

	"github.com/axiaoxin-com/logging"
	"github.com/patrickmn/go-cache"
	"golang.org/x/time/rate"
)

// MemRatelimiter  进程内存 limiter
type MemRatelimiter struct {
	*rate.Limiter
	*cache.Cache
	Expire time.Duration
}

var (
	// MemRatelimiterCacheExpiration MemRatelimiter key 的过期时间
	MemRatelimiterCacheExpiration = time.Minute * 60
	// MemRatelimiterCacheCleanInterval MemRatelimiter 过期 key 的清理时间间隔
	MemRatelimiterCacheCleanInterval = time.Minute * 60
)

// NewMemRatelimiter 根据配置信息创建 mem limiter
func NewMemRatelimiter() *MemRatelimiter {
	// 创建 mem cache
	memCache := cache.New(MemRatelimiterCacheExpiration, MemRatelimiterCacheCleanInterval)
	return &MemRatelimiter{
		Cache: memCache,
	}
}

// Allow 使用 time/rate 的 token bucket 算法判断给定 key 和对应的限制速率下是否被允许
// tokenFillInterval 每隔多长时间往桶中放一个 Token
// bucketSize 代表 Token 桶的容量大小
func (r *MemRatelimiter) Allow(ctx context.Context, key string, tokenFillInterval time.Duration, bucketSize int) bool {
	// 参数小于 0 时不限制
	if tokenFillInterval.Seconds() <= 0 || bucketSize <= 0 {
		return true
	}

	tokenRate := rate.Every(tokenFillInterval)
	limiterI, exists := r.Cache.Get(key)
	if !exists {
		limiter := rate.NewLimiter(tokenRate, bucketSize)
		limiter.Allow()
		r.Cache.Set(key, limiter, MemRatelimiterCacheExpiration)
		return true
	}

	if limiter, ok := limiterI.(*rate.Limiter); ok {
		isAllow := limiter.Allow()
		r.Cache.Set(key, limiter, MemRatelimiterCacheExpiration)
		return isAllow
	}

	logging.Error(nil, "MemRatelimiter assert limiter error")
	return true

}
