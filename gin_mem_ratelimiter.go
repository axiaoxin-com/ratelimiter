package ratelimiter

import (
	"github.com/gin-gonic/gin"
)

// GinMemRatelimiter 按配置信息生成进程内存限频中间件
func GinMemRatelimiter(conf GinRatelimiterConfig) gin.HandlerFunc {
	if conf.TokenBucketConfig == nil {
		panic("GinRatelimiterConfig must implement the TokenBucketConfig callback function")
	}
	limiter := NewMemRatelimiter()

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

		if !limiter.Allow(c, limitKey, tokenFillInterval, bucketSize) {
			limitedHandler(c)
			return
		}
		c.Next()
	}
}
