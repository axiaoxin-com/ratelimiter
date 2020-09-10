package ratelimiter

import (
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

// GinRedisRatelimiter gin redis 分布式请求频率限制
func GinRedisRatelimiter(rdb *redis.Client, everyMicrosecond, bucketCapacity int) gin.HandlerFunc {
	conf := GinRatelimiterConfig{
		LimitKey:       defaultGinLimitKey,
		LimitedHandler: defaultGinLimitedHandler,
		Bucket: BucketConfig{
			Capacity:             bucketCapacity,
			FillEveryMicrosecond: everyMicrosecond,
			ExpireSecond:         defaultBucketExpireSecond,
		},
	}
	return GinRedisRatelimiterWithConfig(rdb, conf)
}

// GinRedisRatelimiterWithConfig 按配置信息生成 redis 限频中间件
func GinRedisRatelimiterWithConfig(rdb *redis.Client, conf GinRatelimiterConfig) gin.HandlerFunc {
	limiter := NewRedisRatelimiter(rdb, conf.Bucket)
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

		// 在 redis 中执行 lua 脚本
		if !limiter.Allow(c, limitKey) {
			limitedHandler(c)
			return
		}
		c.Next()
	}
}
