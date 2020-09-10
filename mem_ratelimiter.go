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
	Bucket BucketConfig
	*cache.Cache
}

// NewMemRatelimiter 根据配置信息创建 mem limiter
func NewMemRatelimiter(conf BucketConfig) *MemRatelimiter {
	// 创建 limiter
	every := defaultBucketFillEveryMicrosecond
	if conf.FillEveryMicrosecond > 0 {
		every = conf.FillEveryMicrosecond
	}
	limit := rate.Every(time.Duration(every) * time.Microsecond)
	burst := defaultBucketCapacity
	if conf.Capacity > 0 {
		burst = conf.Capacity
	}
	expire := defaultBucketExpireSecond
	if conf.ExpireSecond > 0 {
		expire = conf.ExpireSecond
	}
	limiter := rate.NewLimiter(limit, burst)
	// 创建 mem cache
	defaultExpiration := time.Duration(defaultBucketExpireSecond)
	if conf.ExpireSecond > 0 {
		defaultExpiration = time.Duration(conf.ExpireSecond)
	}
	cleanupInterval := time.Duration(defaultExpiration + 5)
	memCache := cache.New(defaultExpiration*time.Second, cleanupInterval*time.Second)
	return &MemRatelimiter{
		Limiter: limiter,
		Bucket: BucketConfig{
			FillEveryMicrosecond: every,
			Capacity:             burst,
			ExpireSecond:         expire,
		},
		Cache: memCache,
	}
}

// Allow 判断给定 key 是否被允许
func (r *MemRatelimiter) Allow(ctx context.Context, key string) bool {
	limiterI, exists := r.Cache.Get(key)
	if !exists {
		r.Limiter.Allow()
		r.Cache.Set(key, r.Limiter, time.Duration(r.Bucket.ExpireSecond)*time.Second)
		return true
	}

	if limiter, ok := limiterI.(*rate.Limiter); ok {
		return limiter.Allow()
	}
	logging.Error(nil, "MemRatelimiter assert limiter error")
	return true

}
