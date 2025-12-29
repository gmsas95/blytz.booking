package middleware

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"blytz.cloud/backend/internal/errors"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type RateLimitMiddleware struct {
	redisClient *redis.Client
}

func NewRateLimitMiddleware(redisClient *redis.Client) *RateLimitMiddleware {
	return &RateLimitMiddleware{redisClient: redisClient}
}

func (rlm *RateLimitMiddleware) Limit(requests int, window time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.Background()
		key := fmt.Sprintf("ratelimit:%s", c.ClientIP())

		current, err := rlm.redisClient.Incr(ctx, key).Result()
		if err != nil {
			c.Next()
			return
		}

		if current == 1 {
			rlm.redisClient.Expire(ctx, key, window)
		}

		c.Header("X-RateLimit-Limit", strconv.Itoa(requests))
		remaining := max(0, requests-int(current))
		c.Header("X-RateLimit-Remaining", strconv.Itoa(remaining))

		if current > int64(requests) {
			errors.HandleError(c, errors.RateLimitExceeded("Rate limit exceeded"))
			c.Abort()
			return
		}

		c.Next()
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
