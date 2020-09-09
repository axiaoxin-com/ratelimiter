package ratelimiter

import (
	"testing"
	"time"

	"github.com/axiaoxin-com/goutils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGinMemRatelimiter(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(GinMemRatelimiter(1000*1000, 1))
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
