package middleware

import (
	"net/http"

	"blytz.cloud/backend/internal/validator"

	"github.com/gin-gonic/gin"
)

type ValidationErrors map[string]string

func ValidateRequest(req interface{}) ValidationErrors {
	errors := make(ValidationErrors)

	// Type assertion to handle different request types
	switch r := req.(type) {
	case *RegisterRequest:
		if !validator.ValidateEmail(r.Email) {
			errors["email"] = "Invalid email format"
		}
		if !validator.ValidateName(r.Name) {
			errors["name"] = "Name must be between 2 and 100 characters"
		}
		if !validator.ValidatePassword(r.Password) {
			errors["password"] = "Password must be at least 6 characters with letters and numbers"
		}
	case *LoginRequest:
		if !validator.ValidateEmail(r.Email) {
			errors["email"] = "Invalid email format"
		}
		if !validator.ValidatePassword(r.Password) {
			errors["password"] = "Password is required"
		}
	}

	if len(errors) > 0 {
		return errors
	}
	return nil
}

type RegisterRequest struct {
	Email    string
	Name     string
	Password string
}

type LoginRequest struct {
	Email    string
	Password string
}

func CustomValidator() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip if no body
		if c.Request.ContentLength == 0 {
			c.Next()
			return
		}

		// Check content type
		contentType := c.GetHeader("Content-Type")
		if contentType != "application/json" && contentType != "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Content-Type must be application/json"})
			c.Abort()
			return
		}

		c.Next()
	}
}
