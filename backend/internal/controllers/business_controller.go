package controllers

import (
	"blytz.cloud/backend/internal/errors"
	"blytz.cloud/backend/internal/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type BusinessController struct {
	businessService *services.BusinessService
}

func NewBusinessController(businessService *services.BusinessService) *BusinessController {
	return &BusinessController{businessService: businessService}
}

func (ctrl *BusinessController) GetBySlug(c *gin.Context) {
	slug := c.Param("slug")
	business, err := ctrl.businessService.GetBySlug(c.Request.Context(), slug)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, business)
}

func (ctrl *BusinessController) UpdateBusiness(c *gin.Context) {
	id := c.Param("id")
	var req services.UpdateBusinessRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errors.HandleError(c, errors.Validation("Invalid request", err))
		return
	}

	business, err := ctrl.businessService.UpdateBusiness(c.Request.Context(), id, &req)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, business)
}

func (ctrl *BusinessController) DeleteBusiness(c *gin.Context) {
	id := c.Param("id")
	if err := ctrl.businessService.DeleteBusiness(c.Request.Context(), id); err != nil {
		errors.HandleError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (ctrl *BusinessController) GetSettings(c *gin.Context) {
	id := c.Param("id")
	settings, err := ctrl.businessService.UpdateSettings(c.Request.Context(), id, nil)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, settings)
}

func (ctrl *BusinessController) UpdateSettings(c *gin.Context) {
	id := c.Param("id")
	var req services.BusinessSettingsUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errors.HandleError(c, errors.Validation("Invalid request", err))
		return
	}

	settings, err := ctrl.businessService.UpdateSettings(c.Request.Context(), id, &req)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, settings)
}

func (ctrl *BusinessController) ListBusinesses(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	offset := (page - 1) * limit

	businesses, total, err := ctrl.businessService.ListBusinesses(c.Request.Context(), offset, limit)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": businesses,
		"meta": gin.H{
			"page":        page,
			"limit":       limit,
			"total":       total,
			"total_pages": (total + int64(limit) - 1) / int64(limit),
		},
	})
}
