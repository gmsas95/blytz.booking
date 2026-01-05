package main

import (
	"log"
	"net/http"

	"blytz.cloud/backend/config"
	"blytz.cloud/backend/internal/auth"
	"blytz.cloud/backend/internal/dto"
	"blytz.cloud/backend/internal/email"
	"blytz.cloud/backend/internal/handlers"
	"blytz.cloud/backend/internal/middleware"
	"blytz.cloud/backend/internal/repository"
	"blytz.cloud/backend/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Set Gin mode
	if cfg.Server.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize repository
	repo, err := repository.NewRepository(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize repository: %v", err)
	}

	// Run migrations and seed data
	if err := repo.AutoMigrate(); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	if err := repo.SeedData(); err != nil {
		log.Printf("Warning: Failed to seed data: %v", err)
	}

	// Initialize handlers
	handler := handlers.NewHandler(repo, email.EmailConfig{
		From:     cfg.Email.From,
		Host:     cfg.Email.Host,
		Port:     cfg.Email.Port,
		Username: cfg.Email.Username,
		Password: cfg.Email.Password,
	})

	businessService := services.NewBusinessService(repo.DB)
	subdomainMiddleware := middleware.NewSubdomainMiddleware(businessService, &middleware.SubdomainConfig{
		BaseDomain: cfg.Server.BaseDomain,
	})

	// Setup Gin router
	r := gin.Default()

	// Subdomain middleware (extracts business from subdomain)
	r.Use(subdomainMiddleware.ExtractAndValidate())

	// CORS middleware
	r.Use(func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		allowedOrigins := []string{
			"https://blytz.cloud",
			"http://localhost:3000",
			"http://localhost:8080",
		}

		// Check if origin is allowed
		allowed := false
		for _, allowedOrigin := range allowedOrigins {
			if origin == allowedOrigin {
				allowed = true
				break
			}
		}

		if allowed {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
			c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
			c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")
		}

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Health check
	r.GET("/health", handler.HealthCheck)

	// API v1 routes
	v1 := r.Group("/api/v1")
	{
		// Auth routes (public)
		v1.POST("/auth/register", handler.Register)
		v1.POST("/auth/login", handler.Login)
		v1.POST("/auth/forgot-password", handler.ForgotPassword)
		v1.POST("/auth/reset-password", handler.ResetPassword)

		// Protected routes
		v1.GET("/auth/me", auth.AuthMiddleware(), func(c *gin.Context) {
			userID := c.GetString("user_id")
			userUUID, err := uuid.Parse(userID)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
				return
			}

			user, err := handler.AuthService.GetByID(userUUID)
			if err != nil {
				c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"id":         user.ID.String(),
				"email":      user.Email,
				"name":       user.Name,
				"created_at": user.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			})
		})

		// Businesses
		v1.GET("/businesses", handler.ListBusinesses)
		v1.POST("/businesses", handler.CreateBusiness)
		v1.GET("/businesses/:businessId", handler.GetBusiness)
		v1.PUT("/businesses/:businessId", handler.UpdateBusiness)
		v1.GET("/business/by-subdomain", func(c *gin.Context) {
			slug := c.Query("slug")
			if slug == "" {
				slug = c.GetString("subdomain")
			}
			business, err := businessService.GetBySlug(slug)
			if err != nil {
				if err == services.ErrNotFound {
					c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "Business not found"})
					return
				}
				c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "Failed to fetch business"})
				return
			}
			c.JSON(http.StatusOK, dto.BusinessResponse{
				ID:              business.ID.String(),
				Name:            business.Name,
				Slug:            business.Slug,
				Vertical:        business.Vertical,
				Description:     business.Description,
				ThemeColor:      business.ThemeColor,
				SlotDurationMin: business.SlotDurationMin,
				MaxBookings:     business.MaxBookings,
				CreatedAt:       business.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
				UpdatedAt:       business.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
			})
		})

		// Services
		v1.GET("/businesses/:businessId/services", handler.GetServicesByBusiness)
		v1.POST("/businesses/:businessId/services", handler.CreateService)
		v1.PUT("/businesses/:businessId/services/:serviceId", handler.UpdateService)
		v1.DELETE("/businesses/:businessId/services/:serviceId", handler.DeleteService)

		// Slots
		v1.GET("/businesses/:businessId/slots", handler.GetSlotsByBusiness)
		v1.POST("/businesses/:businessId/slots", handler.CreateSlot)
		v1.DELETE("/businesses/:businessId/slots/:slotId", handler.DeleteSlot)

		// Availability
		v1.GET("/businesses/:businessId/availability", handler.GetAvailability)
		v1.POST("/businesses/:businessId/availability", handler.SetAvailability)
		v1.POST("/businesses/:businessId/slots/generate", handler.GenerateSlots)

		// Bookings
		v1.POST("/bookings", handler.CreateBooking)
		v1.GET("/businesses/:businessId/bookings", handler.ListBookings)
		v1.DELETE("/businesses/:businessId/bookings/:bookingId", handler.CancelBooking)
	}

	// Start server
	log.Printf("Starting server on port %s...", cfg.Server.Port)
	if err := r.Run(":" + cfg.Server.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
