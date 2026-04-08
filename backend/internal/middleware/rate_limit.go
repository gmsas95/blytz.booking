package middleware

import (
	"bytes"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type rateLimitEntry struct {
	count     int
	resetTime time.Time
}

const maxAuthRateLimitBodyBytes = 4096

func RateLimitByIP(limit int, window time.Duration) gin.HandlerFunc {
	return RateLimitByKey(limit, window, func(c *gin.Context) string {
		return c.ClientIP()
	})
}

func RateLimitByIPAndEmail(limit int, window time.Duration) gin.HandlerFunc {
	return RateLimitByKey(limit, window, func(c *gin.Context) string {
		email := strings.TrimSpace(strings.ToLower(c.PostForm("email")))
		if email == "" && c.Request.Body != nil && strings.Contains(strings.ToLower(c.GetHeader("Content-Type")), "application/json") {
			body, _ := io.ReadAll(io.LimitReader(c.Request.Body, maxAuthRateLimitBodyBytes))
			c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
			var payload struct {
				Email string `json:"email"`
			}
			_ = c.ShouldBindJSON(&payload)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
			email = strings.TrimSpace(strings.ToLower(payload.Email))
		}
		if email == "" {
			email = "unknown"
		}
		return c.ClientIP() + ":" + email
	})
}

func RateLimitByKey(limit int, window time.Duration, keyFunc func(c *gin.Context) string) gin.HandlerFunc {
	var (
		mu      sync.Mutex
		entries = map[string]rateLimitEntry{}
	)

	return func(c *gin.Context) {
		now := time.Now()
		key := keyFunc(c)

		mu.Lock()
		for existingKey, entry := range entries {
			if now.After(entry.resetTime) {
				delete(entries, existingKey)
			}
		}

		entry := entries[key]
		if now.After(entry.resetTime) {
			entry = rateLimitEntry{count: 0, resetTime: now.Add(window)}
		}
		entry.count++
		entries[key] = entry
		mu.Unlock()

		if entry.count > limit {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "Too many requests. Please try again later."})
			c.Abort()
			return
		}

		c.Next()
	}
}
