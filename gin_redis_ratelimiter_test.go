package ratelimiter

import (
	"testing"
	"time"

	"github.com/axiaoxin-com/goutils"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGinRedisRatelimiter(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	rdb, err := goutils.NewRedisClient(&redis.Options{})
	require.Nil(t, err)
	r.Use(GinRedisRatelimiter(rdb, GinRatelimiterConfig{
		TokenBucketConfig: func(c *gin.Context) (time.Duration, int) {
			return 1 * time.Second, 1
		},
	}))
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, "hi")
	})
	time.Sleep(1 * time.Second)
	recorder, err := goutils.RequestHTTPHandler(r, "GET", "/", nil, nil)
	assert.Nil(t, err)
	assert.Equal(t, recorder.Code, 200)
	recorder, err = goutils.RequestHTTPHandler(r, "GET", "/", nil, nil)
	assert.Nil(t, err)
	assert.Equal(t, recorder.Code, 429)
	time.Sleep(1 * time.Second)
	recorder, err = goutils.RequestHTTPHandler(r, "GET", "/", nil, nil)
	assert.Nil(t, err)
	assert.Equal(t, recorder.Code, 200)
}
