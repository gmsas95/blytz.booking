package controllers

import (
	"blytz.cloud/backend/internal/errors"
	"blytz.cloud/backend/internal/services"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type SlotController struct {
	slotService *services.SlotService
}

func NewSlotController(slotService *services.SlotService) *SlotController {
	return &SlotController{slotService: slotService}
}

func (ctrl *SlotController) ListAvailableSlots(c *gin.Context) {
	businessID := c.Param("businessId")
	serviceID := c.Query("service_id")

	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		errors.HandleError(c, errors.Validation("Invalid start_date format", err))
		return
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		errors.HandleError(c, errors.Validation("Invalid end_date format", err))
		return
	}

	slots, err := ctrl.slotService.ListAvailableSlots(c.Request.Context(), businessID, serviceID, startDate, endDate)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": slots})
}

func (ctrl *SlotController) GetSlot(c *gin.Context) {
	id := c.Param("id")
	slot, err := ctrl.slotService.GetSlot(c.Request.Context(), id)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, slot)
}

func (ctrl *SlotController) DeleteSlot(c *gin.Context) {
	id := c.Param("id")
	if err := ctrl.slotService.DeleteSlot(c.Request.Context(), id); err != nil {
		errors.HandleError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (ctrl *SlotController) CreateSlots(c *gin.Context) {
	var req services.CreateSlotsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errors.HandleError(c, errors.Validation("Invalid request", err))
		return
	}

	slots, err := ctrl.slotService.CreateSlots(c.Request.Context(), &req)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"created_count": len(slots),
		"slots":         slots,
	})
}

func (ctrl *SlotController) CreateRecurringSchedule(c *gin.Context) {
	var req services.RecurringScheduleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errors.HandleError(c, errors.Validation("Invalid request", err))
		return
	}

	schedule, err := ctrl.slotService.CreateRecurringSchedule(c.Request.Context(), &req)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, schedule)
}

func (ctrl *SlotController) DeleteRecurringSchedule(c *gin.Context) {
	id := c.Param("id")
	if err := ctrl.slotService.DeleteRecurringSchedule(c.Request.Context(), id); err != nil {
		errors.HandleError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (ctrl *SlotController) ListAvailableSlotsBySlug(c *gin.Context) {
	slug := c.Param("slug")
	serviceID := c.Query("service_id")

	var startDate, endDate *time.Time

	if startDateStr := c.Query("start_date"); startDateStr != "" {
		if parsed, err := time.Parse("2006-01-02", startDateStr); err == nil {
			startDate = &parsed
		} else {
			errors.HandleError(c, errors.Validation("Invalid start_date format", err))
			return
		}
	}

	if endDateStr := c.Query("end_date"); endDateStr != "" {
		if parsed, err := time.Parse("2006-01-02", endDateStr); err == nil {
			endDate = &parsed
		} else {
			errors.HandleError(c, errors.Validation("Invalid end_date format", err))
			return
		}
	}

	slots, err := ctrl.slotService.ListAvailableSlotsBySlug(c.Request.Context(), slug, serviceID, startDate, endDate)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": slots})
}
