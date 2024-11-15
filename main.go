package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis_rate/v10"
	"github.com/redis/go-redis/v9"
)

func main() {
	opts := &redis.Options{
		Addr: "redis:6379",
		DB:   0,
	}
	rdb := redis.NewClient(opts)

	if _, pingErr := rdb.Ping(context.Background()).Result(); pingErr != nil {
		panic(fmt.Errorf("failed to ping redis db; %v", pingErr))
	}

	limiter := redis_rate.NewLimiter(rdb)
	rateMiddleware := func() gin.HandlerFunc {
		return func(c *gin.Context) {
			res, err := limiter.Allow(c, "pingRequest", redis_rate.PerMinute(5))
			if err != nil {
				panic(fmt.Errorf("rate limiter failed; %v", err))
			}
			fmt.Printf("Allowed %v; remaining %v; [ts: %v]", res.Allowed, res.Remaining, time.Now())
			if res.Remaining == 0 {
				c.AbortWithStatus(http.StatusTooManyRequests)
				return
			}
			c.Next()
		}
	}

	r := gin.Default()
	r.Use(rateMiddleware())
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.Run(":8888")
}
