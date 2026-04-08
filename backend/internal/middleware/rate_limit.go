package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type rateLimitEntry struct {
	count     int
	resetTime time.Time
}

func RateLimitByIP(limit int, window time.Duration) gin.HandlerFunc {
	var (
		mu      sync.Mutex
		entries = map[string]rateLimitEntry{}
	)

	return func(c *gin.Context) {
		ip := c.ClientIP()
		now := time.Now()

		mu.Lock()
		entry := entries[ip]
		if now.After(entry.resetTime) {
			entry = rateLimitEntry{count: 0, resetTime: now.Add(window)}
		}
		entry.count++
		entries[ip] = entry
		mu.Unlock()

		if entry.count > limit {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "Too many requests. Please try again later."})
			c.Abort()
			return
		}

		c.Next()
	}
}
