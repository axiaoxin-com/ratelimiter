package ratelimiter

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

// GinRedisRatelimiter 按配置信息生成 redis 限频中间件
func GinRedisRatelimiter(rdb *redis.Client, conf GinRatelimiterConfig) gin.HandlerFunc {
	limiter := NewRedisRatelimiter(rdb)
	return func(c *gin.Context) {
		// 获取 limit key
		var limitKey string
		if conf.LimitKey != nil {
			limitKey = conf.LimitKey(c)
		} else {
			limitKey = defaultGinLimitKey(c)
		}

		limitedHandler := defaultGinLimitedHandler
		if conf.LimitedHandler != nil {
			limitedHandler = conf.LimitedHandler
		}

		var tokenFillInterval time.Duration
		var bucketSize int
		if conf.TokenBucketConfig != nil {
			tokenFillInterval, bucketSize = conf.TokenBucketConfig(c)
		}

		// 在 redis 中执行 lua 脚本
		if !limiter.Allow(c, limitKey, tokenFillInterval, bucketSize) {
			limitedHandler(c)
			return
		}
		c.Next()
	}
}
