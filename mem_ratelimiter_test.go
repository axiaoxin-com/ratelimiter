package ratelimiter

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMemRatelimiter(t *testing.T) {
	conf := BucketConfig{
		Capacity:             1,
		FillEveryMicrosecond: 1000 * 1000,
		ExpireSecond:         60,
	}
	limiter := NewMemRatelimiter(conf)
	ctx := context.Background()
	assert.True(t, limiter.Allow(ctx, "key"))
	assert.False(t, limiter.Allow(ctx, "key"))
	time.Sleep(1 * time.Second)
	assert.True(t, limiter.Allow(ctx, "key"))
}
