package middleware

import (
	"blytz.cloud/backend/internal/errors"
	"blytz.cloud/backend/internal/types"
	"blytz.cloud/backend/internal/utils"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AuthMiddleware struct {
	jwtManager *utils.JWTManager
}

func NewAuthMiddleware(jwtManager *utils.JWTManager) *AuthMiddleware {
	return &AuthMiddleware{jwtManager: jwtManager}
}

func (am *AuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			errors.HandleError(c, errors.Unauthorized("Missing authorization header"))
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			errors.HandleError(c, errors.Unauthorized("Invalid authorization header format"))
			c.Abort()
			return
		}

		tokenString := parts[1]
		claims, err := am.jwtManager.ValidateToken(tokenString)
		if err != nil {
			errors.HandleError(c, errors.Unauthorized("Invalid or expired token"))
			c.Abort()
			return
		}

		ctx := c.Request.Context()
		ctx = types.WithUserID(ctx, claims.UserID)
		if claims.BusinessID != uuid.Nil {
			ctx = types.WithBusinessID(ctx, claims.BusinessID)
		}
		ctx = types.WithUserRole(ctx, claims.Role)

		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

func RequireRoles(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		userRole, err := types.GetUserRoleFromContext(ctx)
		if err != nil {
			errors.HandleError(c, errors.Unauthorized("User role not found in context"))
			c.Abort()
			return
		}

		for _, role := range allowedRoles {
			if userRole == role {
				c.Next()
				return
			}
		}

		errors.HandleError(c, errors.Forbidden("Insufficient permissions"))
		c.Abort()
	}
}

func RequireSuperAdmin() gin.HandlerFunc {
	return RequireRoles("superadmin")
}

func RequireOwnerOrAdmin() gin.HandlerFunc {
	return RequireRoles("owner", "admin")
}

func RequireOwnerOrAdminOrStaff() gin.HandlerFunc {
	return RequireRoles("owner", "admin", "staff")
}

func RequireBusinessOwner() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		userRole, err := types.GetUserRoleFromContext(ctx)
		if err != nil {
			errors.HandleError(c, errors.Unauthorized("User role not found"))
			c.Abort()
			return
		}

		if userRole == "superadmin" {
			c.Next()
			return
		}

		if userRole != "owner" {
			errors.HandleError(c, errors.Forbidden("Only business owners can perform this action"))
			c.Abort()
			return
		}

		businessID, err := types.GetBusinessIDFromContext(ctx)
		if err != nil || businessID == uuid.Nil {
			errors.HandleError(c, errors.Forbidden("Business ID not found in context"))
			c.Abort()
			return
		}

		paramBusinessID := c.Param("business_id")
		if paramBusinessID != "" && paramBusinessID != businessID.String() {
			errors.HandleError(c, errors.Forbidden("You can only access your own business"))
			c.Abort()
			return
		}

		c.Next()
	}
}
