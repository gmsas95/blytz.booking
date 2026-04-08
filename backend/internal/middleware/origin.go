package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RequireAllowedOrigin(allowedOrigins []string) gin.HandlerFunc {
	allowed := make(map[string]struct{}, len(allowedOrigins))
	for _, origin := range allowedOrigins {
		allowed[origin] = struct{}{}
	}

	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		if origin == "" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Origin header required"})
			c.Abort()
			return
		}
		if _, ok := allowed[origin]; !ok {
			c.JSON(http.StatusForbidden, gin.H{"error": "Origin not allowed"})
			c.Abort()
			return
		}
		c.Next()
	}
}
