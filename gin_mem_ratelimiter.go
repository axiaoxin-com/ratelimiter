package ratelimiter

import (
	"github.com/gin-gonic/gin"
)

// GinMemRatelimiter gin 进程内存级别的请求频率限制
func GinMemRatelimiter(everyMicrosecond, bucketCapacity int) gin.HandlerFunc {
	conf := GinRatelimiterConfig{
		LimitKey:       defaultGinLimitKey,
		LimitedHandler: defaultGinLimitedHandler,
		Bucket: BucketConfig{
			Capacity:             bucketCapacity,
			FillEveryMicrosecond: everyMicrosecond,
			ExpireSecond:         defaultBucketExpireSecond,
		},
	}
	return GinMemRatelimiterWithConfig(conf)
}

// GinMemRatelimiterWithConfig 按配置信息生成进程内存限频中间件
func GinMemRatelimiterWithConfig(conf GinRatelimiterConfig) gin.HandlerFunc {
	limiter := NewMemRatelimiter(conf.Bucket)

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

		if !limiter.Allow(c, limitKey) {
			limitedHandler(c)
			return
		}
		c.Next()
	}
}
