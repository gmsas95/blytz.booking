package main

import (
	"log"
	"time"

	"blytz.cloud/backend/config"
	"blytz.cloud/backend/internal/controllers"
	"blytz.cloud/backend/internal/middleware"
	"blytz.cloud/backend/internal/repository"
	"blytz.cloud/backend/internal/services"
	"blytz.cloud/backend/internal/utils"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.LoadConfig()

	if cfg.Server.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	repo, err := repository.NewRepository(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize repository: %v", err)
	}

	if err := repo.AutoMigrate(); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	jwtManager := utils.NewJWTManager(cfg.JWT.Secret)

	authMiddleware := middleware.NewAuthMiddleware(jwtManager)
	rateLimitMiddleware := middleware.NewRateLimitMiddleware(repo.Redis)

	authService := services.NewAuthService(
		repository.NewUserRepository(repo.DB),
		repository.NewRefreshTokenRepository(repo.DB),
		jwtManager,
		repo.DB,
	)

	bookingService := services.NewBookingService(
		repository.NewBookingRepository(repo.DB),
		repository.NewSlotRepository(repo.DB),
		repository.NewServiceRepository(repo.DB),
		repository.NewCustomerRepository(repo.DB),
		repository.NewBookingHistoryRepository(repo.DB),
		repo.DB,
	)

	authController := controllers.NewAuthController(authService)
	bookingController := controllers.NewBookingController(bookingService)

	r := gin.New()

	r.Use(middleware.Recovery())
	r.Use(middleware.RequestID())
	r.Use(middleware.Logger())
	r.Use(middleware.CORS())
	r.Use(rateLimitMiddleware.Limit(100, time.Minute))

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "healthy", "version": "1.0.0"})
	})

	public := r.Group("/api/v1")
	{
		auth := public.Group("/auth")
		{
			auth.POST("/register", authController.Register)
			auth.POST("/login", authController.Login)
		}

		bookings := public.Group("/bookings")
		{
			bookings.POST("", bookingController.CreateBooking)
		}

		protected := public.Group("")
		protected.Use(authMiddleware.RequireAuth())
		{
			auth := protected.Group("/auth")
			{
				auth.POST("/logout", authController.Logout)
				auth.POST("/refresh", authController.RefreshToken)
				auth.GET("/me", authController.GetMe)
			}

			bookings := protected.Group("/bookings")
			{
				bookings.GET("/:id", bookingController.GetBooking)
				bookings.PATCH("/:id/status", bookingController.UpdateBookingStatus)
				bookings.DELETE("/:id", bookingController.CancelBooking)

				businessBookings := protected.Group("/businesses/:businessId/bookings")
				businessBookings.Use(middleware.RequireBusinessOwner())
				{
					businessBookings.GET("", bookingController.ListBookings)
				}
			}
		}
	}

	log.Printf("Starting server on port %s...", cfg.Server.Port)
	if err := r.Run(":" + cfg.Server.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
