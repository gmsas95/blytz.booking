package handlers

import (
	"net/http"

	"blytz.cloud/backend/internal/models"
	"blytz.cloud/backend/internal/repository"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	Repo *repository.Repository
}

func NewHandler(repo *repository.Repository) *Handler {
	return &Handler{Repo: repo}
}

func (h *Handler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "healthy",
	})
}

// Business Handlers
func (h *Handler) ListBusinesses(c *gin.Context) {
	var businesses []models.Business
	if err := h.Repo.DB.Find(&businesses).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch businesses"})
		return
	}

	c.JSON(http.StatusOK, businesses)
}

func (h *Handler) GetBusiness(c *gin.Context) {
	id := c.Param("businessId")

	var business models.Business
	if err := h.Repo.DB.Where("id = ?", id).First(&business).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Business not found"})
		return
	}

	c.JSON(http.StatusOK, business)
}

// Service Handlers
func (h *Handler) GetServicesByBusiness(c *gin.Context) {
	businessID := c.Param("businessId")

	var services []models.Service
	if err := h.Repo.DB.Where("business_id = ?", businessID).Find(&services).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch services"})
		return
	}

	c.JSON(http.StatusOK, services)
}

// Slot Handlers
func (h *Handler) GetSlotsByBusiness(c *gin.Context) {
	businessID := c.Param("businessId")

	var slots []models.Slot
	if err := h.Repo.DB.Where("business_id = ? AND is_booked = ?", businessID, false).Order("start_time").Find(&slots).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch slots"})
		return
	}

	c.JSON(http.StatusOK, slots)
}

// Booking Handlers
func (h *Handler) CreateBooking(c *gin.Context) {
	var booking models.Booking
	if err := c.ShouldBindJSON(&booking); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create booking
	if err := h.Repo.DB.Create(&booking).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create booking"})
		return
	}

	// Mark slot as booked
	h.Repo.DB.Model(&models.Slot{}).Where("id = ?", booking.SlotID).Update("is_booked", true)

	c.JSON(http.StatusCreated, booking)
}

func (h *Handler) ListBookings(c *gin.Context) {
	businessID := c.Param("businessId")

	var bookings []models.Booking
	if err := h.Repo.DB.Where("business_id = ?", businessID).Order("created_at DESC").Find(&bookings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch bookings"})
		return
	}

	c.JSON(http.StatusOK, bookings)
}
