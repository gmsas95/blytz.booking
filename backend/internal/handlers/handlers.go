package handlers

import (
	"net/http"

	"blytz.cloud/backend/internal/dto"
	"blytz.cloud/backend/internal/models"
	"blytz.cloud/backend/internal/repository"
	"blytz.cloud/backend/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	Repo            *repository.Repository
	AuthService     *services.AuthService
	BusinessService *services.BusinessService
	ServiceService  *services.ServiceService
	SlotService     *services.SlotService
	BookingService  *services.BookingService
}

func NewHandler(repo *repository.Repository) *Handler {
	return &Handler{
		Repo:            repo,
		AuthService:     services.NewAuthService(repo.DB),
		BusinessService: services.NewBusinessService(repo.DB),
		ServiceService:  services.NewServiceService(repo.DB),
		SlotService:     services.NewSlotService(repo.DB),
		BookingService:  services.NewBookingService(repo.DB),
	}
}

func (h *Handler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "healthy",
	})
}

// Business Handlers
func (h *Handler) ListBusinesses(c *gin.Context) {
	businesses, err := h.BusinessService.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "Failed to fetch businesses"})
		return
	}

	response := make([]dto.BusinessResponse, len(businesses))
	for i, b := range businesses {
		response[i] = dto.BusinessResponse{
			ID:          b.ID.String(),
			Name:        b.Name,
			Slug:        b.Slug,
			Vertical:    b.Vertical,
			Description: b.Description,
			ThemeColor:  b.ThemeColor,
			CreatedAt:   b.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:   b.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
	}

	c.JSON(http.StatusOK, response)
}

func (h *Handler) GetBusiness(c *gin.Context) {
	id := c.Param("businessId")
	businessID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "Invalid business ID"})
		return
	}

	business, err := h.BusinessService.GetByID(businessID)
	if err != nil {
		if err == services.ErrNotFound {
			c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "Business not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "Failed to fetch business"})
		return
	}

	c.JSON(http.StatusOK, dto.BusinessResponse{
		ID:          business.ID.String(),
		Name:        business.Name,
		Slug:        business.Slug,
		Vertical:    business.Vertical,
		Description: business.Description,
		ThemeColor:  business.ThemeColor,
		CreatedAt:   business.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   business.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	})
}

func (h *Handler) CreateBusiness(c *gin.Context) {
	var req dto.CreateBusinessRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	business := &models.Business{
		ID:          uuid.New(),
		Name:        req.Name,
		Slug:        req.Slug,
		Vertical:    req.Vertical,
		Description: req.Description,
		ThemeColor:  req.ThemeColor,
	}

	if err := h.BusinessService.Create(business); err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "Failed to create business"})
		return
	}

	c.JSON(http.StatusCreated, dto.BusinessResponse{
		ID:          business.ID.String(),
		Name:        business.Name,
		Slug:        business.Slug,
		Vertical:    business.Vertical,
		Description: business.Description,
		ThemeColor:  business.ThemeColor,
		CreatedAt:   business.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   business.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	})
}

func (h *Handler) UpdateBusiness(c *gin.Context) {
	id := c.Param("businessId")
	businessID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "Invalid business ID"})
		return
	}

	var req dto.UpdateBusinessRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	updates := &models.Business{}
	if req.Name != nil {
		updates.Name = *req.Name
	}
	if req.Vertical != nil {
		updates.Vertical = *req.Vertical
	}
	if req.Description != nil {
		updates.Description = *req.Description
	}
	if req.ThemeColor != nil {
		updates.ThemeColor = *req.ThemeColor
	}

	if err := h.BusinessService.Update(businessID, updates); err != nil {
		if err == services.ErrNotFound {
			c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "Business not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "Failed to update business"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Business updated successfully"})
}

// Service Handlers
func (h *Handler) GetServicesByBusiness(c *gin.Context) {
	businessID := c.Param("businessId")
	businessUUID, err := uuid.Parse(businessID)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "Invalid business ID"})
		return
	}

	services, err := h.ServiceService.GetByBusiness(businessUUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "Failed to fetch services"})
		return
	}

	response := make([]dto.ServiceResponse, len(services))
	for i, s := range services {
		response[i] = dto.ServiceResponse{
			ID:            s.ID.String(),
			BusinessID:    s.BusinessID.String(),
			Name:          s.Name,
			Description:   s.Description,
			DurationMin:   s.DurationMin,
			TotalPrice:    s.TotalPrice,
			DepositAmount: s.DepositAmount,
			CreatedAt:     s.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:     s.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
	}

	c.JSON(http.StatusOK, response)
}

func (h *Handler) CreateService(c *gin.Context) {
	businessIDStr := c.Param("businessId")
	businessID, err := uuid.Parse(businessIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "Invalid business ID"})
		return
	}

	var req dto.CreateServiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	service := &models.Service{
		ID:            uuid.New(),
		BusinessID:    businessID,
		Name:          req.Name,
		Description:   req.Description,
		DurationMin:   req.DurationMin,
		TotalPrice:    req.TotalPrice,
		DepositAmount: req.DepositAmount,
	}

	if err := h.ServiceService.Create(service); err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "Failed to create service"})
		return
	}

	c.JSON(http.StatusCreated, dto.ServiceResponse{
		ID:            service.ID.String(),
		BusinessID:    service.BusinessID.String(),
		Name:          service.Name,
		Description:   service.Description,
		DurationMin:   service.DurationMin,
		TotalPrice:    service.TotalPrice,
		DepositAmount: service.DepositAmount,
		CreatedAt:     service.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:     service.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	})
}

func (h *Handler) DeleteService(c *gin.Context) {
	businessID := c.Param("businessId")
	serviceID := c.Param("serviceId")

	businessUUID, err := uuid.Parse(businessID)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "Invalid business ID"})
		return
	}

	serviceUUID, err := uuid.Parse(serviceID)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "Invalid service ID"})
		return
	}

	// Verify service belongs to business
	service, err := h.ServiceService.GetByID(serviceUUID)
	if err != nil {
		if err == services.ErrNotFound {
			c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "Service not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "Failed to fetch service"})
		return
	}

	if service.BusinessID != businessUUID {
		c.JSON(http.StatusForbidden, dto.ErrorResponse{Error: "Service does not belong to this business"})
		return
	}

	if err := h.ServiceService.Delete(serviceUUID); err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "Failed to delete service"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Service deleted successfully"})
}

// Slot Handlers
func (h *Handler) GetSlotsByBusiness(c *gin.Context) {
	businessID := c.Param("businessId")
	businessUUID, err := uuid.Parse(businessID)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "Invalid business ID"})
		return
	}

	slots, err := h.SlotService.GetAvailableByBusiness(businessUUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "Failed to fetch slots"})
		return
	}

	response := make([]dto.SlotResponse, len(slots))
	for i, s := range slots {
		response[i] = dto.SlotResponse{
			ID:         s.ID.String(),
			BusinessID: s.BusinessID.String(),
			StartTime:  s.StartTime.Format("2006-01-02T15:04:05Z07:00"),
			EndTime:    s.EndTime.Format("2006-01-02T15:04:05Z07:00"),
			IsBooked:   s.IsBooked,
			CreatedAt:  s.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:  s.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
	}

	c.JSON(http.StatusOK, response)
}

// Booking Handlers
func (h *Handler) CreateBooking(c *gin.Context) {
	var req dto.CreateBookingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	booking := &models.Booking{
		BusinessID: uuid.MustParse(req.BusinessID),
		ServiceID:  uuid.MustParse(req.ServiceID),
		SlotID:     uuid.MustParse(req.SlotID),
		Customer: models.CustomerDetails{
			Name:  req.Customer.Name,
			Email: req.Customer.Email,
			Phone: req.Customer.Phone,
		},
	}

	if err := h.BookingService.Create(booking); err != nil {
		if err == services.ErrBadRequest {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "Invalid booking request"})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "Failed to create booking"})
		return
	}

	c.JSON(http.StatusCreated, dto.BookingResponse{
		ID:          booking.ID.String(),
		BusinessID:  booking.BusinessID.String(),
		ServiceID:   booking.ServiceID.String(),
		SlotID:      booking.SlotID.String(),
		ServiceName: booking.ServiceName,
		SlotTime:    booking.SlotTime.Format("2006-01-02T15:04:05Z07:00"),
		Customer: dto.CustomerDetails{
			Name:  booking.Customer.Name,
			Email: booking.Customer.Email,
			Phone: booking.Customer.Phone,
		},
		Status:      string(booking.Status),
		DepositPaid: booking.DepositPaid,
		TotalPrice:  booking.TotalPrice,
		CreatedAt:   booking.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   booking.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	})
}

func (h *Handler) ListBookings(c *gin.Context) {
	businessID := c.Param("businessId")
	businessUUID, err := uuid.Parse(businessID)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "Invalid business ID"})
		return
	}

	bookings, err := h.BookingService.GetByBusiness(businessUUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "Failed to fetch bookings"})
		return
	}

	response := make([]dto.BookingResponse, len(bookings))
	for i, b := range bookings {
		response[i] = dto.BookingResponse{
			ID:          b.ID.String(),
			BusinessID:  b.BusinessID.String(),
			ServiceID:   b.ServiceID.String(),
			SlotID:      b.SlotID.String(),
			ServiceName: b.ServiceName,
			SlotTime:    b.SlotTime.Format("2006-01-02T15:04:05Z07:00"),
			Customer: dto.CustomerDetails{
				Name:  b.Customer.Name,
				Email: b.Customer.Email,
				Phone: b.Customer.Phone,
			},
			Status:      string(b.Status),
			DepositPaid: b.DepositPaid,
			TotalPrice:  b.TotalPrice,
			CreatedAt:   b.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:   b.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
	}

	c.JSON(http.StatusOK, response)
}

// Auth Handlers

func (h *Handler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	user, token, err := h.AuthService.Register(req.Email, req.Name, req.Password)
	if err != nil {
		if err == services.ErrConflict {
			c.JSON(http.StatusConflict, dto.ErrorResponse{Error: "Email already registered"})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, dto.AuthResponse{
		Token: token,
		User: dto.UserResponse{
			ID:        user.ID.String(),
			Email:     user.Email,
			Name:      user.Name,
			CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		},
	})
}

func (h *Handler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	user, token, err := h.AuthService.Login(req.Email, req.Password)
	if err != nil {
		if err == services.ErrUnauthorized {
			c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "Invalid credentials"})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "Failed to login"})
		return
	}

	c.JSON(http.StatusOK, dto.AuthResponse{
		Token: token,
		User: dto.UserResponse{
			ID:        user.ID.String(),
			Email:     user.Email,
			Name:      user.Name,
			CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		},
	})
}
