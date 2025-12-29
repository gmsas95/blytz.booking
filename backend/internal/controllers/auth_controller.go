package controllers

import (
	"blytz.cloud/backend/internal/errors"
	"blytz.cloud/backend/internal/services"
	"blytz.cloud/backend/internal/types"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	authService *services.AuthService
}

func NewAuthController(authService *services.AuthService) *AuthController {
	return &AuthController{authService: authService}
}

func (ctrl *AuthController) Register(c *gin.Context) {
	var req services.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errors.HandleError(c, errors.Validation("Invalid request", err))
		return
	}

	response, err := ctrl.authService.Register(c.Request.Context(), &req)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, response)
}

func (ctrl *AuthController) Login(c *gin.Context) {
	var req services.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errors.HandleError(c, errors.Validation("Invalid request", err))
		return
	}

	response, err := ctrl.authService.Login(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Unauthorized",
			"message": err.Error(),
			"code":    errors.ErrCodeUnauthorized,
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (ctrl *AuthController) Logout(c *gin.Context) {
	userID, err := types.GetUserIDFromContext(c.Request.Context())
	if err != nil {
		errors.HandleError(c, errors.Unauthorized("User ID not found in context"))
		return
	}

	if err := ctrl.authService.Logout(c.Request.Context(), userID); err != nil {
		errors.HandleError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (ctrl *AuthController) RefreshToken(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" validate:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		errors.HandleError(c, errors.Validation("Invalid request", err))
		return
	}

	response, err := ctrl.authService.RefreshToken(c.Request.Context(), req.RefreshToken)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}

func (ctrl *AuthController) GetMe(c *gin.Context) {
	userID, err := types.GetUserIDFromContext(c.Request.Context())
	if err != nil {
		errors.HandleError(c, errors.Unauthorized("User ID not found in context"))
		return
	}

	user, business, err := ctrl.authService.GetMe(c.Request.Context(), userID)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	response := gin.H{
		"user": user,
	}

	if business != nil {
		response["business"] = business
	}

	c.JSON(http.StatusOK, response)
}
