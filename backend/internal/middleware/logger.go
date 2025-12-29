package middleware

import (
	"log"
	"time"

	"blytz.cloud/backend/internal/types"
	"github.com/gin-gonic/gin"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		duration := time.Since(start)
		requestID := types.GetRequestIDFromContext(c.Request.Context())

		statusCode := c.Writer.Status()
		method := c.Request.Method
		path := c.Request.URL.Path

		log.Printf("[%s] %s %s %d %v %s",
			time.Now().Format(time.RFC3339),
			method,
			path,
			statusCode,
			duration,
			requestID,
		)
	}
}
