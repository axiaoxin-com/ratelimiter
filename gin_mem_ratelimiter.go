package ratelimiter

import (
	"time"

	"github.com/gin-gonic/gin"
)

// GinMemRatelimiter 按配置信息生成进程内存限频中间件
func GinMemRatelimiter(conf GinRatelimiterConfig) gin.HandlerFunc {
	limiter := NewMemRatelimiter()

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

		if !limiter.Allow(c, limitKey, tokenFillInterval, bucketSize) {
			limitedHandler(c)
			return
		}
		c.Next()
	}
}
