package ratelimiter

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// defaultGinLimitKey 使用客户端 IP 生成默认的限频 key
func defaultGinLimitKey(c *gin.Context) string {
	return fmt.Sprintf("pink-lady:ratelimiter:%s", c.ClientIP())
}

// defaultGinLimitedHandler 限频触发返回 429
func defaultGinLimitedHandler(c *gin.Context) {
	c.AbortWithStatus(http.StatusTooManyRequests)
}

// GinRatelimiterConfig Gin Ratelimiter 中间件的配置信息
type GinRatelimiterConfig struct {
	// LimitKey 生成限频 key 的函数，不传使用默认的对 IP 维度进行限制
	LimitKey func(*gin.Context) string
	// LimitedHandler 触发限频时的 handler
	LimitedHandler func(*gin.Context)
	// Bucket 令牌桶配置
	Bucket BucketConfig
}
