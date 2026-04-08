package handlers

import (
	"net/http"
	"time"

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
	CustomerService *services.CustomerService
	VehicleService  *services.VehicleService
	JobService      *services.JobService
}

func getCurrentUserID(c *gin.Context) (uuid.UUID, error) {
	return uuid.Parse(c.GetString("user_id"))
}

func NewHandler(repo *repository.Repository) *Handler {
	return &Handler{
		Repo:            repo,
		AuthService:     services.NewAuthService(repo.DB),
		BusinessService: services.NewBusinessService(repo.DB),
		ServiceService:  services.NewServiceService(repo.DB),
		SlotService:     services.NewSlotService(repo.DB),
		BookingService:  services.NewBookingService(repo.DB),
		CustomerService: services.NewCustomerService(repo.DB),
		VehicleService:  services.NewVehicleService(repo.DB),
		JobService:      services.NewJobService(repo.DB),
	}
}

func currentBusinessID(c *gin.Context) (uuid.UUID, error) {
	businessID := c.GetString("business_id")
	if businessID == "" {
		businessID = c.Param("businessId")
	}
	return uuid.Parse(businessID)
}

func customerResponse(customer models.Customer) dto.CustomerResponse {
	return dto.CustomerResponse{
		ID:         customer.ID.String(),
		BusinessID: customer.BusinessID.String(),
		Name:       customer.Name,
		Email:      customer.Email,
		Phone:      customer.Phone,
		Notes:      customer.Notes,
		CreatedAt:  customer.CreatedAt.Format(time.RFC3339),
		UpdatedAt:  customer.UpdatedAt.Format(time.RFC3339),
	}
}

func vehicleResponse(vehicle models.Vehicle) dto.VehicleResponse {
	return dto.VehicleResponse{
		ID:           vehicle.ID.String(),
		BusinessID:   vehicle.BusinessID.String(),
		CustomerID:   vehicle.CustomerID.String(),
		Year:         vehicle.Year,
		Make:         vehicle.Make,
		Model:        vehicle.Model,
		Color:        vehicle.Color,
		LicensePlate: vehicle.LicensePlate,
		Customer:     customerResponse(vehicle.Customer),
		CreatedAt:    vehicle.CreatedAt.Format(time.RFC3339),
		UpdatedAt:    vehicle.UpdatedAt.Format(time.RFC3339),
	}
}

func jobResponse(job models.Job) dto.JobResponse {
	bookingID := ""
	if job.BookingID != nil {
		bookingID = job.BookingID.String()
	}
	return dto.JobResponse{
		ID:          job.ID.String(),
		BusinessID:  job.BusinessID.String(),
		CustomerID:  job.CustomerID.String(),
		VehicleID:   job.VehicleID.String(),
		BookingID:   bookingID,
		Title:       job.Title,
		Status:      string(job.Status),
		ScheduledAt: job.ScheduledAt.Format(time.RFC3339),
		Notes:       job.Notes,
		Customer:    customerResponse(job.Customer),
		Vehicle:     vehicleResponse(job.Vehicle),
		CreatedAt:   job.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   job.UpdatedAt.Format(time.RFC3339),
	}
}

func (h *Handler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "healthy",
	})
}

func (h *Handler) GetCurrentUser(c *gin.Context) {
	userID, err := getCurrentUserID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "Invalid user ID"})
		return
	}

	user, err := h.AuthService.GetByID(userID)
	if err != nil {
		if err == services.ErrNotFound {
			c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "Failed to fetch user"})
		return
	}

	memberships, err := h.AuthService.GetMemberships(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "Failed to fetch memberships"})
		return
	}

	membershipResponse := make([]dto.MembershipResponse, len(memberships))
	activeBusinessID := ""
	for i, membership := range memberships {
		membershipResponse[i] = dto.MembershipResponse{
			ID:         membership.ID.String(),
			BusinessID: membership.BusinessID.String(),
			Role:       string(membership.Role),
			Business: dto.MembershipBusinessResponse{
				ID:          membership.Business.ID.String(),
				Name:        membership.Business.Name,
				Slug:        membership.Business.Slug,
				Vertical:    membership.Business.Vertical,
				Description: membership.Business.Description,
				ThemeColor:  membership.Business.ThemeColor,
			},
		}
		if activeBusinessID == "" {
			activeBusinessID = membership.BusinessID.String()
		}
	}

	c.JSON(http.StatusOK, dto.CurrentUserResponse{
		User: dto.UserResponse{
			ID:        user.ID.String(),
			Email:     user.Email,
			Name:      user.Name,
			CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		},
		Memberships:      membershipResponse,
		ActiveBusinessID: activeBusinessID,
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
			ID:                 s.ID.String(),
			BusinessID:         s.BusinessID.String(),
			Name:               s.Name,
			Description:        s.Description,
			DurationMin:        s.DurationMin,
			TotalPriceMinor:    s.TotalPriceMinor,
			DepositAmountMinor: s.DepositAmountMinor,
			CurrencyCode:       s.CurrencyCode,
			CreatedAt:          s.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:          s.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
	}

	c.JSON(http.StatusOK, response)
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
		if err == services.ErrConflict {
			c.JSON(http.StatusConflict, dto.ErrorResponse{Error: "Selected slot is no longer available"})
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
		Status:           string(booking.Status),
		DepositPaidMinor: booking.DepositPaidMinor,
		TotalPriceMinor:  booking.TotalPriceMinor,
		CurrencyCode:     booking.CurrencyCode,
		CreatedAt:        booking.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:        booking.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	})
}

func (h *Handler) ListBookings(c *gin.Context) {
	businessUUID, err := currentBusinessID(c)
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
			Status:           string(b.Status),
			DepositPaidMinor: b.DepositPaidMinor,
			TotalPriceMinor:  b.TotalPriceMinor,
			CurrencyCode:     b.CurrencyCode,
			CreatedAt:        b.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:        b.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
	}

	c.JSON(http.StatusOK, response)
}

func (h *Handler) ListCustomers(c *gin.Context) {
	businessID, err := currentBusinessID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "Invalid business ID"})
		return
	}

	customers, err := h.CustomerService.GetByBusiness(businessID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "Failed to fetch customers"})
		return
	}

	response := make([]dto.CustomerResponse, len(customers))
	for i, customer := range customers {
		response[i] = customerResponse(customer)
	}
	c.JSON(http.StatusOK, response)
}

func (h *Handler) CreateCustomer(c *gin.Context) {
	businessID, err := currentBusinessID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "Invalid business ID"})
		return
	}

	var req dto.CreateCustomerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	customer := &models.Customer{BusinessID: businessID, Name: req.Name, Email: req.Email, Phone: req.Phone, Notes: req.Notes}
	if err := h.CustomerService.Create(customer); err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "Failed to create customer"})
		return
	}
	c.JSON(http.StatusCreated, customerResponse(*customer))
}

func (h *Handler) ListVehicles(c *gin.Context) {
	businessID, err := currentBusinessID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "Invalid business ID"})
		return
	}

	vehicles, err := h.VehicleService.GetByBusiness(businessID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "Failed to fetch vehicles"})
		return
	}

	response := make([]dto.VehicleResponse, len(vehicles))
	for i, vehicle := range vehicles {
		response[i] = vehicleResponse(vehicle)
	}
	c.JSON(http.StatusOK, response)
}

func (h *Handler) CreateVehicle(c *gin.Context) {
	businessID, err := currentBusinessID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "Invalid business ID"})
		return
	}

	var req dto.CreateVehicleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	customerID, err := uuid.Parse(req.CustomerID)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "Invalid customer ID"})
		return
	}

	vehicle := &models.Vehicle{BusinessID: businessID, CustomerID: customerID, Year: req.Year, Make: req.Make, Model: req.Model, Color: req.Color, LicensePlate: req.LicensePlate}
	if err := h.VehicleService.Create(vehicle); err != nil {
		if err == services.ErrBadRequest {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "Customer does not belong to this workshop"})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "Failed to create vehicle"})
		return
	}

	h.VehicleService.DB.Preload("Customer", "business_id = ?", businessID).First(vehicle, "id = ? AND business_id = ?", vehicle.ID, businessID)
	c.JSON(http.StatusCreated, vehicleResponse(*vehicle))
}

func (h *Handler) ListJobs(c *gin.Context) {
	businessID, err := currentBusinessID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "Invalid business ID"})
		return
	}

	jobs, err := h.JobService.GetByBusiness(businessID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "Failed to fetch jobs"})
		return
	}

	response := make([]dto.JobResponse, len(jobs))
	for i, job := range jobs {
		response[i] = jobResponse(job)
	}
	c.JSON(http.StatusOK, response)
}

func (h *Handler) CreateJob(c *gin.Context) {
	businessID, err := currentBusinessID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "Invalid business ID"})
		return
	}

	var req dto.CreateJobRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	customerID, err := uuid.Parse(req.CustomerID)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "Invalid customer ID"})
		return
	}
	vehicleID, err := uuid.Parse(req.VehicleID)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "Invalid vehicle ID"})
		return
	}
	scheduledAt, err := time.Parse(time.RFC3339, req.ScheduledAt)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "Invalid scheduled_at"})
		return
	}
	status := models.JobStatusScheduled
	if req.Status != "" {
		status = models.JobStatus(req.Status)
	}
	var bookingID *uuid.UUID
	if req.BookingID != "" {
		parsedBookingID, err := uuid.Parse(req.BookingID)
		if err != nil {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "Invalid booking ID"})
			return
		}
		bookingID = &parsedBookingID
	}

	job := &models.Job{BusinessID: businessID, CustomerID: customerID, VehicleID: vehicleID, BookingID: bookingID, Title: req.Title, Status: status, ScheduledAt: scheduledAt, Notes: req.Notes}
	if err := h.JobService.Create(job); err != nil {
		if err == services.ErrBadRequest {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "Customer, vehicle, or booking does not belong to this workshop"})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "Failed to create job"})
		return
	}

	h.JobService.DB.Preload("Customer", "business_id = ?", businessID).Preload("Vehicle", "business_id = ?", businessID).First(job, "id = ? AND business_id = ?", job.ID, businessID)
	c.JSON(http.StatusCreated, jobResponse(*job))
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
