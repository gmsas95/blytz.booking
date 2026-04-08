package middleware

import (
	"net/http"

	"blytz.cloud/backend/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func RequireBusinessMembership(authService *services.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := uuid.Parse(c.GetString("user_id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			c.Abort()
			return
		}

		businessID, err := uuid.Parse(c.Param("businessId"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid business ID"})
			c.Abort()
			return
		}

		hasAccess, err := authService.UserHasBusinessAccess(userID, businessID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify membership"})
			c.Abort()
			return
		}
		if !hasAccess {
			c.JSON(http.StatusForbidden, gin.H{"error": "You do not have access to this workshop"})
			c.Abort()
			return
		}

		c.Set("business_id", businessID.String())
		c.Next()
	}
}
