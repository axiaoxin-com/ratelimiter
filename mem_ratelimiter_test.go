package ratelimiter

import (
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
	assert.True(t, limiter.Allow("key"))
	assert.False(t, limiter.Allow("key"))
	time.Sleep(1 * time.Second)
	assert.True(t, limiter.Allow("key"))
}
