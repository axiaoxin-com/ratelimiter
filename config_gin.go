package ratelimiter

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// DefaultGinLimitKey 使用客户端 IP 生成默认的限频 key
func DefaultGinLimitKey(c *gin.Context) string {
	return fmt.Sprintf("pink-lady:ratelimiter:%s:%s", c.ClientIP(), c.FullPath())
}

// DefaultGinLimitedHandler 限频触发返回 429
func DefaultGinLimitedHandler(c *gin.Context) {
	c.AbortWithStatus(http.StatusTooManyRequests)
}

// GinRatelimiterConfig Gin Ratelimiter 中间件的配置信息
type GinRatelimiterConfig struct {
	// LimitKey 生成限频 key 的函数，不传使用默认的对 IP 维度进行限制
	LimitKey func(*gin.Context) string
	// LimitedHandler 触发限频时的 handler
	LimitedHandler func(*gin.Context)
	// TokenBucketConfig 获取 token bucket 每次放入一个token的时间间隔和桶大小配置
	TokenBucketConfig func(*gin.Context) (time.Duration, int)
}
