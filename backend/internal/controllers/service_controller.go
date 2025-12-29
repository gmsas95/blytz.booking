package controllers

import (
	"blytz.cloud/backend/internal/errors"
	"blytz.cloud/backend/internal/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ServiceController struct {
	serviceService *services.ServiceService
}

func NewServiceController(serviceService *services.ServiceService) *ServiceController {
	return &ServiceController{serviceService: serviceService}
}

func (ctrl *ServiceController) CreateService(c *gin.Context) {
	businessID := c.Param("businessId")
	var req services.CreateServiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errors.HandleError(c, errors.Validation("Invalid request", err))
		return
	}

	req.BusinessID = businessID

	service, err := ctrl.serviceService.CreateService(c.Request.Context(), &req)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, service)
}

func (ctrl *ServiceController) GetService(c *gin.Context) {
	id := c.Param("serviceId")
	service, err := ctrl.serviceService.GetService(c.Request.Context(), id)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, service)
}

func (ctrl *ServiceController) ListServices(c *gin.Context) {
	businessID := c.Param("businessId")

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	offset := (page - 1) * limit

	services, total, err := ctrl.serviceService.ListServices(c.Request.Context(), businessID, offset, limit)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": services,
		"meta": gin.H{
			"page":        page,
			"limit":       limit,
			"total":       total,
			"total_pages": (total + int64(limit) - 1) / int64(limit),
		},
	})
}

func (ctrl *ServiceController) UpdateService(c *gin.Context) {
	id := c.Param("serviceId")
	var req services.UpdateServiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errors.HandleError(c, errors.Validation("Invalid request", err))
		return
	}

	service, err := ctrl.serviceService.UpdateService(c.Request.Context(), id, &req)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, service)
}

func (ctrl *ServiceController) DeleteService(c *gin.Context) {
	id := c.Param("serviceId")
	if err := ctrl.serviceService.DeleteService(c.Request.Context(), id); err != nil {
		errors.HandleError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (ctrl *ServiceController) ListServicesBySlug(c *gin.Context) {
	slug := c.Param("slug")

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	offset := (page - 1) * limit

	services, total, err := ctrl.serviceService.ListServicesBySlug(c.Request.Context(), slug, offset, limit)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": services,
		"meta": gin.H{
			"page":        page,
			"limit":       limit,
			"total":       total,
			"total_pages": (total + int64(limit) - 1) / int64(limit),
		},
	})
}
