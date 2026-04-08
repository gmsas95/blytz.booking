package auth

import (
	"net/http"
	"strings"

	"blytz.cloud/backend/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type SessionValidator interface {
	ValidateUserSession(userID uuid.UUID, tokenVersion int) (*models.User, error)
}

var cookieName = "blytz_session"

func SetCookieName(name string) {
	if name != "" {
		cookieName = name
	}
}

func CookieName() string {
	return cookieName
}

func AuthMiddleware(authService SessionValidator) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := ""
		if cookie, err := c.Cookie(cookieName); err == nil {
			tokenString = cookie
		}
		if tokenString == "" {
			authHeader := c.GetHeader("Authorization")
			if authHeader != "" {
				tokenString = strings.TrimPrefix(authHeader, "Bearer ")
				if tokenString == authHeader {
					tokenString = ""
				}
			}
		}
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
			c.Abort()
			return
		}

		claims, err := ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		userID, err := uuid.Parse(claims.UserID)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}
		if _, err := authService.ValidateUserSession(userID, claims.TokenVersion); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("email", claims.Email)
		c.Next()
	}
}
