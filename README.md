# ratelimiter

[![Build Status](https://travis-ci.org/axiaoxin-com/ratelimiter.svg?branch=master)](https://travis-ci.org/axiaoxin-com/ratelimiter)
[![go report card](https://goreportcard.com/badge/github.com/axiaoxin-com/ratelimiter)](https://goreportcard.com/report/github.com/axiaoxin-com/ratelimiter)

[中文 README](./README_CN.md)

Simple version implementation of token bucket request frequency limiting.

ratelimiter library that supports in-memory and distributed eventually consistent redis stores (includes Gin middleware)

- [lua-ngx-ratelimiter](./lua-ngx-ratelimiter): a token bucket frequency limiting implementation of lua + nginx + redis
- [MemRatelimiter](./mem_ratelimiter.go): a process memory limiter implemented with [rate](https://github.com/golang/time/tree/master/rate) + [go-cache](https://github.com/patrickmn/go-cache)
- [RedisRatelimiter](./redis_ratelimiter.go): a distributed limiter implemented with redis + lua
- [GinMemRatelimiter](./gin_mem_ratelimiter.go): encapsulating the MemRatelimiter as a gin middleware
- [GinRedisRatelimiter](./gin_redis_ratelimiter.go): encapsulating the RedisRatelimiter as a gin middleware

## Go PKG Installation

```
go get -u github.com/axiaoxin-com/ratelimiter
```

## Gin Middleware Example

**[GinMemRatelimiter](./example/gin_mem_ratelimiter.go)**

```
package main

import (
	"time"

	"github.com/axiaoxin-com/ratelimiter"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.New()
	// Put a token into the token bucket every 1s
	// Maximum 1 request allowed per second
	r.Use(ratelimiter.GinMemRatelimiter(ratelimiter.GinRatelimiterConfig{
		// config: how to generate a limit key
		LimitKey: func(c *gin.Context) string {
			return c.ClientIP()
		},
		// config: how to respond when limiting
		LimitedHandler: func(c *gin.Context) {
			c.JSON(200, "too many requests!!!")
			c.Abort()
			return
		},
		// config: return ratelimiter token fill interval and bucket size
		TokenBucketConfig: func(*gin.Context) (time.Duration, int) {
			return time.Second * 1, 1
		},
	}))
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, "hi")
	})
	r.Run()
}
```

**[GinRedisRatelimiter](./example/gin_redis_ratelimiter.go)**

```
package main

import (
	"time"

	"github.com/axiaoxin-com/goutils"
	"github.com/axiaoxin-com/ratelimiter"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

func main() {
	r := gin.New()
	// Put a token into the token bucket every 1s
	// Maximum 1 request allowed per second
	rdb, err := goutils.NewRedisClient(&redis.Options{})
	if err != nil {
		panic(err)
	}
	r.Use(ratelimiter.GinRedisRatelimiter(rdb, ratelimiter.GinRatelimiterConfig{
		// config: how to generate a limit key
		LimitKey: func(c *gin.Context) string {
			return c.ClientIP()
		},
		// config: how to respond when limiting
		LimitedHandler: func(c *gin.Context) {
			c.JSON(200, "too many requests!!!")
			c.Abort()
			return
		},
		// config: return ratelimiter token fill interval and bucket size
		TokenBucketConfig: func(*gin.Context) (time.Duration, int) {
			return time.Second * 1, 1
		},
	}))
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, "hi")
	})
	r.Run()
}
```

## Ratelimiter can be directly used in golang program. Examples:

**[MemRatelimiter](./example/mem_ratelimiter.go)**

```
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
```

**[RedisRatelimiter](./example/redis_ratelimiter.go)**

```
package main

import (
	"context"
	"fmt"
	"time"

	"github.com/axiaoxin-com/goutils"
	"github.com/axiaoxin-com/ratelimiter"
	"github.com/go-redis/redis/v8"
)

func main() {
	rdb, err := goutils.NewRedisClient(&redis.Options{})
	if err != nil {
		panic(err)
	}

	limiter := ratelimiter.NewRedisRatelimiter(rdb)
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
```
