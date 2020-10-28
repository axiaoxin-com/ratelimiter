package main

import (
	"context"
	"fmt"
	"time"

	"github.com/axiaoxin-com/ratelimiter"
)

func main() {
	limiter := ratelimiter.NewMemRatelimiter()
	limitKey := "uniq_limit_key"
	tokenFillInterval := time.Second * 1
	bucketSize := 1
	for i := 0; i < 3; i++ {
		// 1st and 3nd is allowed
		if i == 2 {
			time.Sleep(time.Second * 1)
		}
		isAllow := limiter.Allow(context.TODO(), limitKey, tokenFillInterval, bucketSize)
		fmt.Println(i, time.Now(), isAllow)
	}
}
