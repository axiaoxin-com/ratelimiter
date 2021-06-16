package ratelimiter

import (
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

// GinRedisRatelimiter 按配置信息生成 redis 限频中间件
func GinRedisRatelimiter(rdb *redis.Client, conf GinRatelimiterConfig) gin.HandlerFunc {
	if conf.TokenBucketConfig == nil {
		panic("GinRatelimiterConfig must implement the TokenBucketConfig callback function")
	}
	limiter := NewRedisRatelimiter(rdb)
	return func(c *gin.Context) {
		// 获取 limit key
		limitKey := DefaultGinLimitKey(c)
		if conf.LimitKey != nil {
			limitKey = conf.LimitKey(c)
		}

		limitedHandler := DefaultGinLimitedHandler
		if conf.LimitedHandler != nil {
			limitedHandler = conf.LimitedHandler
		}

		tokenFillInterval, bucketSize := conf.TokenBucketConfig(c)

		// 在 redis 中执行 lua 脚本
		if !limiter.Allow(c, limitKey, tokenFillInterval, bucketSize) {
			limitedHandler(c)
			return
		}
		c.Next()
	}
}
