package main

import (
	"log"
	"net/http"
	"strings"
	"time"

	"blytz.cloud/backend/config"
	"blytz.cloud/backend/internal/auth"
	"blytz.cloud/backend/internal/handlers"
	"blytz.cloud/backend/internal/middleware"
	"blytz.cloud/backend/internal/repository"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()
	if cfg.JWT.Secret == "" || cfg.JWT.Secret == "change-me-in-production" {
		log.Fatal("JWT_SECRET must be explicitly configured")
	}
	if cfg.Database.Password == "" {
		log.Fatal("DB_PASSWORD must be explicitly configured")
	}
	auth.SetJWTSecret(cfg.JWT.Secret)
	auth.SetCookieName(cfg.JWT.CookieName)
	handlers.SetForceSecureCookies(cfg.JWT.ForceSecure)

	// Set Gin mode
	if cfg.Server.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize repository
	repo, err := repository.NewRepository(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize repository: %v", err)
	}

	if cfg.Startup.AutoMigrate {
		if err := repo.AutoMigrate(); err != nil {
			log.Fatalf("Failed to run migrations: %v", err)
		}
	}

	if cfg.Startup.BackfillMoney {
		if err := repo.BackfillMoneyToMinorUnits("USD"); err != nil {
			log.Fatalf("Failed to backfill money fields: %v", err)
		}
	}

	if cfg.Startup.SeedData {
		if err := repo.SeedData(); err != nil {
			log.Printf("Warning: Failed to seed data: %v", err)
		}
	}

	// Initialize handlers
	handler := handlers.NewHandler(repo)

	// Setup Gin router
	r := gin.Default()
	if err := r.SetTrustedProxies(cfg.JWT.TrustedProxies); err != nil {
		log.Fatalf("Failed to configure trusted proxies: %v", err)
	}

	allowedOrigins := make(map[string]struct{}, len(cfg.CORS.AllowedOrigins))
	for _, origin := range cfg.CORS.AllowedOrigins {
		allowedOrigins[origin] = struct{}{}
	}

	// CORS middleware
	r.Use(func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		if origin != "" {
			if _, ok := allowedOrigins[origin]; ok {
				c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
				c.Writer.Header().Set("Vary", "Origin")
				c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
			} else if c.Request.Method == http.MethodOptions {
				c.AbortWithStatus(http.StatusForbidden)
				return
			}
		}

		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == http.MethodOptions {
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
		authRoutes := v1.Group("/auth")
		authRoutes.Use(middleware.RequireAllowedOrigin(cfg.CORS.AllowedOrigins), middleware.RateLimitByIP(30, time.Minute), middleware.RateLimitByIPAndEmail(10, time.Minute))
		authRoutes.POST("/register", handler.Register)
		authRoutes.POST("/login", handler.Login)
		v1.POST("/auth/logout", middleware.RequireAllowedOrigin(cfg.CORS.AllowedOrigins), auth.AuthMiddleware(handler.AuthService), handler.Logout)

		// Protected routes
		v1.GET("/auth/me", auth.AuthMiddleware(handler.AuthService), handler.GetCurrentUser)

		// Businesses
		v1.GET("/businesses", handler.ListBusinesses)
		v1.GET("/businesses/:businessId", handler.GetBusiness)

		// Services
		v1.GET("/businesses/:businessId/services", handler.GetServicesByBusiness)

		// Slots
		v1.GET("/businesses/:businessId/slots", handler.GetSlotsByBusiness)

		// Bookings
		v1.POST("/bookings", handler.CreateBooking)

		operator := v1.Group("/businesses/:businessId")
		operator.Use(auth.AuthMiddleware(handler.AuthService), middleware.RequireBusinessMembership(handler.AuthService))
		{
			operator.GET("/bookings", handler.ListBookings)
			operator.GET("/customers", handler.ListCustomers)
			operator.POST("/customers", middleware.RequireAllowedOrigin(cfg.CORS.AllowedOrigins), handler.CreateCustomer)
			operator.GET("/vehicles", handler.ListVehicles)
			operator.POST("/vehicles", middleware.RequireAllowedOrigin(cfg.CORS.AllowedOrigins), handler.CreateVehicle)
			operator.GET("/jobs", handler.ListJobs)
			operator.POST("/jobs", middleware.RequireAllowedOrigin(cfg.CORS.AllowedOrigins), handler.CreateJob)
		}
	}

	// Start server
	log.Printf("Allowed CORS origins: %s", strings.Join(cfg.CORS.AllowedOrigins, ", "))
	log.Printf("Startup flags: auto_migrate=%t seed_data=%t backfill_money=%t", cfg.Startup.AutoMigrate, cfg.Startup.SeedData, cfg.Startup.BackfillMoney)
	log.Printf("Starting server on port %s...", cfg.Server.Port)
	if err := r.Run(":" + cfg.Server.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
