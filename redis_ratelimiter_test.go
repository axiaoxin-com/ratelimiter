package ratelimiter

import (
	"context"
	"testing"
	"time"

	"github.com/axiaoxin-com/goutils"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRedisRatelimiter(t *testing.T) {
	conf := BucketConfig{
		Capacity:             1,
		FillEveryMicrosecond: 1000 * 1000,
		ExpireSecond:         60 * 5,
	}
	rdb, err := goutils.NewRedisClient(&redis.Options{})
	require.Nil(t, err)
	limiter := NewRedisRatelimiter(rdb, conf)
	ctx := context.Background()
	assert.True(t, limiter.Allow(ctx, "key"))
	assert.False(t, limiter.Allow(ctx, "key"))
	time.Sleep(1 * time.Second)
	assert.True(t, limiter.Allow(ctx, "key"))
}
