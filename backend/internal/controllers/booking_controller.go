package controllers

import (
	"blytz.cloud/backend/internal/errors"
	"blytz.cloud/backend/internal/models"
	"blytz.cloud/backend/internal/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type BookingController struct {
	bookingService *services.BookingService
}

func NewBookingController(bookingService *services.BookingService) *BookingController {
	return &BookingController{bookingService: bookingService}
}

func (ctrl *BookingController) CreateBooking(c *gin.Context) {
	var req services.CreateBookingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errors.HandleError(c, errors.Validation("Invalid request", err))
		return
	}

	booking, err := ctrl.bookingService.CreateBooking(c.Request.Context(), &req)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"booking": booking})
}

func (ctrl *BookingController) GetBooking(c *gin.Context) {
	id := c.Param("id")
	booking, err := ctrl.bookingService.GetBooking(c.Request.Context(), id)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, booking)
}

func (ctrl *BookingController) ListBookings(c *gin.Context) {
	businessID := c.Param("businessId")

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	offset := (page - 1) * limit

	bookings, total, err := ctrl.bookingService.ListBookings(c.Request.Context(), businessID, offset, limit)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": bookings,
		"meta": gin.H{
			"page":        page,
			"limit":       limit,
			"total":       total,
			"total_pages": (total + int64(limit) - 1) / int64(limit),
		},
	})
}

func (ctrl *BookingController) UpdateBookingStatus(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Status string `json:"status" validate:"required,oneof=PENDING CONFIRMED COMPLETED CANCELLED NO_SHOW"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		errors.HandleError(c, errors.Validation("Invalid request", err))
		return
	}

	booking, err := ctrl.bookingService.UpdateBookingStatus(c.Request.Context(), id, models.BookingStatus(req.Status), "")
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, booking)
}

func (ctrl *BookingController) CancelBooking(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Reason string `json:"reason"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		errors.HandleError(c, errors.Validation("Invalid request", err))
		return
	}

	if err := ctrl.bookingService.CancelBooking(c.Request.Context(), id, req.Reason, ""); err != nil {
		errors.HandleError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}
