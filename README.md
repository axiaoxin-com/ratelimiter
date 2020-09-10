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

**GinMemRatelimiter**

```
package main

import (
	"github.com/axiaoxin-com/ratelimiter"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.New()
	// Put a token into the token bucket every 1000 * 1000 microseconds
	// Maximum 1 request allowed per second
	r.Use(ratelimiter.GinMemRatelimiter(1000*1000, 1))
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, "hi")
	})
	r.Run()
}
```

**GinRedisRatelimiter**

```
package main

import (
	"github.com/axiaoxin-com/ratelimiter"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.New()
	rdb, err := goutils.NewRedisClient(&redis.Options{})
	// Put a token into the token bucket every 1000 * 1000 microseconds
	// Maximum 1 request allowed per second
	r.Use(ratelimiter.GinRedisRatelimiter(rdb, 1000*1000, 1))
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, "hi")
	})
	r.Run()
}
```

## Ratelimiter can be directly used in golang program. Examples:

**MemRatelimiter**

```
conf := ratelimiter.BucketConfig{
    Capacity:             1,
    FillEveryMicrosecond: 1000 * 1000,
    ExpireSecond:         60,
}
limiter := ratelimiter.NewMemRatelimiter(conf)
ctx := context.TODO()
if limiter.Allow(ctx, "somekey") {
    // do something
}
```

**RedisRatelimiter**

```
conf := BucketConfig{
    Capacity:             1,
    FillEveryMicrosecond: 1000 * 1000,
    ExpireSecond:         60 * 5,
}
rdb, _ := goutils.NewRedisClient(&redis.Options{})
limiter := NewRedisRatelimiter(rdb, conf)
ctx := context.Background()
if limiter.Allow(ctx, "key") {
    // do something
}
```
