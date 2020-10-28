package ratelimiter

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMemRatelimiter(t *testing.T) {
	limiter := NewMemRatelimiter()
	ctx := context.Background()
	assert.True(t, limiter.Allow(ctx, "key", time.Second*1, 1))
	assert.False(t, limiter.Allow(ctx, "key", time.Second*1, 1))
	time.Sleep(1 * time.Second)
	assert.True(t, limiter.Allow(ctx, "key", time.Second*1, 1))
}
