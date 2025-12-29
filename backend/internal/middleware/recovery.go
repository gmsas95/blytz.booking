package middleware

import (
	"github.com/gin-gonic/gin"
)

func Recovery() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		c.JSON(500, gin.H{
			"error":   "Internal Server Error",
			"code":    "INTERNAL_ERROR",
			"message": "An unexpected error occurred",
		})
	})
}
