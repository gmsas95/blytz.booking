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

	businessService := services.NewBusinessService(
		repository.NewBusinessRepository(repo.DB),
		repo.DB,
	)

	serviceService := services.NewServiceService(
		repository.NewServiceRepository(repo.DB),
		repository.NewBusinessRepository(repo.DB),
		repo.DB,
	)

	slotService := services.NewSlotService(
		repository.NewSlotRepository(repo.DB),
		repository.NewBusinessRepository(repo.DB),
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
	businessController := controllers.NewBusinessController(businessService)
	serviceController := controllers.NewServiceController(serviceService)
	slotController := controllers.NewSlotController(slotService)
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

		publicBookings := public.Group("/bookings")
		{
			publicBookings.POST("", bookingController.CreateBooking)
			publicBookings.GET("/:id", bookingController.GetBooking)
		}

		publicBusinesses := public.Group("/businesses")
		{
			publicBusinesses.GET("", businessController.ListBusinesses)
			publicBusinesses.GET("/health", func(c *gin.Context) { c.JSON(200, gin.H{"status": "ok"}) })
			publicBusinesses.GET("/slug/:slug", businessController.GetBySlug)
			publicBusinesses.GET("/slug/:slug/services", serviceController.ListServicesBySlug)
			publicBusinesses.GET("/slug/:slug/slots", slotController.ListAvailableSlotsBySlug)
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

			businesses := protected.Group("/businesses")
			{
				business := businesses.Group("/:id")
				business.Use(middleware.RequireBusinessOwner())
				{
					business.GET("", func(c *gin.Context) {
						c.JSON(200, gin.H{"id": c.Param("id"), "name": "Business"})
					})
					business.PATCH("", businessController.UpdateBusiness)
					business.DELETE("", businessController.DeleteBusiness)
					business.GET("/settings", businessController.GetSettings)
					business.PATCH("/settings", businessController.UpdateSettings)

					services := business.Group("/services")
					{
						services.POST("", serviceController.CreateService)
						services.GET("", serviceController.ListServices)
						services.GET("/:serviceId", serviceController.GetService)
						services.PATCH("/:serviceId", serviceController.UpdateService)
						services.DELETE("/:serviceId", serviceController.DeleteService)
					}

					slots := business.Group("/slots")
					{
						slots.GET("", slotController.ListAvailableSlots)
						slots.POST("", slotController.CreateSlots)
						slots.POST("/recurring", slotController.CreateRecurringSchedule)
						slots.DELETE("/:id", slotController.DeleteSlot)
						slots.DELETE("/recurring/:id", slotController.DeleteRecurringSchedule)
					}

					bookings := business.Group("/bookings")
					{
						bookings.GET("", bookingController.ListBookings)
						bookings.PATCH("/:id/status", bookingController.UpdateBookingStatus)
						bookings.DELETE("/:id", bookingController.CancelBooking)
					}
				}
			}

			protectedBookings := protected.Group("/bookings")
			{
				protectedBookings.PATCH("/:id/status", bookingController.UpdateBookingStatus)
				protectedBookings.DELETE("/:id", bookingController.CancelBooking)
			}
		}
	}

	log.Printf("Starting server on port %s...", cfg.Server.Port)
	if err := r.Run(":" + cfg.Server.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
