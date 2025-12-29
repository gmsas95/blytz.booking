# Recommended Backend Architecture

## Current Structure

```
backend/
├── cmd/
│   └── server/
│       └── main.go              # ✅ Entry point
├── config/
│   └── config.go              # ✅ Configuration
├── internal/
│   ├── handlers/
│   │   └── handlers.go        # ⚠️ Too much logic, direct DB access
│   ├── models/
│   │   └── models.go         # ✅ Data models
│   └── repository/
│       └── repository.go       # ⚠️ No abstraction, basic CRUD only
├── Dockerfile
├── go.mod
└── go.sum
```

**Issues:**
- No service layer (business logic in handlers)
- No middleware (CORS inline in main.go)
- No validation (only JSON binding)
- No error handling (generic errors)
- No testing
- No separation of concerns

---

## Recommended Structure

```
backend/
├── cmd/
│   └── server/
│       └── main.go                  # Application entry point
│
├── config/
│   ├── config.go                   # Configuration loader
│   ├── database.go                 # DB connection setup
│   └── redis.go                   # Redis connection setup
│
├── internal/
│   ├── api/
│   │   └── router.go              # Route definitions grouped by resource
│   │
│   ├── controllers/
│   │   ├── auth_controller.go      # Auth HTTP handlers
│   │   ├── business_controller.go   # Business HTTP handlers
│   │   ├── service_controller.go   # Service HTTP handlers
│   │   ├── slot_controller.go      # Slot HTTP handlers
│   │   ├── booking_controller.go   # Booking HTTP handlers
│   │   ├── customer_controller.go  # Customer HTTP handlers
│   │   ├── payment_controller.go   # Payment HTTP handlers
│   │   ├── analytics_controller.go # Analytics HTTP handlers
│   │   ├── webhook_controller.go   # Webhook HTTP handlers
│   │   └── admin_controller.go    # Admin HTTP handlers
│   │
│   ├── services/
│   │   ├── auth_service.go         # Business logic for auth
│   │   ├── business_service.go    # Business logic for businesses
│   │   ├── service_service.go     # Business logic for services
│   │   ├── slot_service.go        # Business logic for slots
│   │   ├── booking_service.go     # Business logic for bookings
│   │   ├── customer_service.go    # Business logic for customers
│   │   ├── payment_service.go     # Business logic for payments
│   │   ├── stripe_service.go      # Stripe API integration
│   │   ├── analytics_service.go   # Business logic for analytics
│   │   ├── webhook_service.go      # Webhook delivery
│   │   └── notification_service.go # Email/SMS notifications
│   │
│   ├── repository/
│   │   ├── repository.go          # Repository interface definitions
│   │   ├── user_repository.go     # User data access
│   │   ├── business_repository.go  # Business data access
│   │   ├── service_repository.go  # Service data access
│   │   ├── slot_repository.go     # Slot data access
│   │   ├── booking_repository.go  # Booking data access
│   │   ├── customer_repository.go # Customer data access
│   │   ├── payment_repository.go  # Payment data access
│   │   └── subscription_repository.go # Subscription data access
│   │
│   ├── models/
│   │   ├── models.go             # All database models
│   │   ├── enums.go             # Enums (status, roles, etc.)
│   │   └── dto/                # Data Transfer Objects
│   │       ├── business_dto.go
│   │       ├── service_dto.go
│   │       ├── booking_dto.go
│   │       └── analytics_dto.go
│   │
│   ├── middleware/
│   │   ├── auth.go               # JWT authentication
│   │   ├── authorization.go       # RBAC/permission checking
│   │   ├── tenant.go             # Subdomain/tenant extraction
│   │   ├── rate_limit.go         # Rate limiting
│   │   ├── cors.go               # CORS configuration
│   │   ├── request_id.go         # Correlation ID injection
│   │   ├── logger.go            # Request/response logging
│   │   ├── recovery.go           # Panic recovery
│   │   └── security_headers.go   # Security headers
│   │
│   ├── validators/
│   │   ├── validator.go          # Main validator setup
│   │   ├── business_validator.go # Business-specific validation
│   │   ├── service_validator.go  # Service-specific validation
│   │   ├── slot_validator.go     # Slot-specific validation
│   │   ├── booking_validator.go  # Booking-specific validation
│   │   └── user_validator.go    # User auth validation
│   │
│   ├── errors/
│   │   ├── errors.go             # Custom error types
│   │   ├── app_error.go          # Application error struct
│   │   └── error_handler.go      # Error response formatter
│   │
│   ├── types/
│   │   ├── context.go            # Custom context types
│   │   └── pagination.go         # Pagination types
│   │
│   └── utils/
│       ├── crypto.go             # Password hashing
│       ├── jwt.go               # JWT token generation/validation
│       ├── time.go              # Time utilities
│       └── currency.go           # Currency formatting
│
├── pkg/
│   ├── logger/
│   │   ├── logger.go            # Structured logger interface
│   │   └── zap.go              # Zap implementation
│   └── metrics/
│       ├── metrics.go            # Metrics interface
│       └── prometheus.go        # Prometheus implementation
│
├── migrations/
│   ├── 001_init_schema.up.sql   # Database migrations
│   ├── 001_init_schema.down.sql
│   └── ...                   # Additional migrations
│
├── scripts/
│   ├── migrate.sh               # Run migrations
│   └── seed.sh                 # Seed test data
│
├── test/
│   ├── integration/             # Integration tests
│   │   ├── auth_test.go
│   │   ├── booking_test.go
│   │   └── ...
│   └── mocks/                  # Mock implementations
│       ├── repository_mock.go
│       └── stripe_mock.go
│
├── .env.example
├── Dockerfile
├── docker-compose.yml
├── go.mod
└── go.sum
```

---

## Architecture Principles

### 1. Layered Architecture

```
┌─────────────────────────────────────┐
│   Controllers (HTTP Handlers)     │  ← Parse HTTP request/response
├─────────────────────────────────────┤
│   Services (Business Logic)       │  ← Domain logic, validation
├─────────────────────────────────────┤
│   Repository (Data Access)       │  ← Database queries
├─────────────────────────────────────┤
│   Models (Data Structures)       │  ← Database entities
└─────────────────────────────────────┘
```

**Rules:**
- Controllers **NEVER** access database directly
- Controllers call services
- Services call repositories
- Repositories only do data access (no business logic)
- Models are passive data structures

### 2. Dependency Inversion

Services depend on **interfaces**, not concrete implementations:

```go
// internal/repository/repository.go
type BookingRepository interface {
    Create(booking *models.Booking) error
    GetByID(id uuid.UUID) (*models.Booking, error)
    GetByBusinessID(businessID uuid.UUID, opts ...QueryOption) ([]*models.Booking, error)
    Update(booking *models.Booking) error
    Delete(id uuid.UUID) error
    // ...
}

// internal/services/booking_service.go
type BookingService struct {
    bookingRepo repository.BookingRepository
    slotRepo    repository.SlotRepository
    logger      logger.Logger
}

// Dependency injection
func NewBookingService(
    bookingRepo repository.BookingRepository,
    slotRepo repository.SlotRepository,
    log logger.Logger,
) *BookingService {
    return &BookingService{
        bookingRepo: bookingRepo,
        slotRepo:    slotRepo,
        logger:      log,
    }
}
```

### 3. Middleware Chain

```
Request → Request ID → Security Headers → CORS → Auth → Rate Limit → Tenant → Recovery → Handler
                                                                                       ↓
                                                                                     Response
```

### 4. Error Handling

Standardized error responses:

```go
// internal/errors/errors.go
type ErrorCode string

const (
    ErrCodeValidation    ErrorCode = "VALIDATION_ERROR"
    ErrCodeNotFound     ErrorCode = "NOT_FOUND"
    ErrCodeUnauthorized ErrorCode = "UNAUTHORIZED"
    ErrCodeForbidden    ErrorCode = "FORBIDDEN"
    ErrCodeConflict     ErrorCode = "CONFLICT"
    ErrCodeInternal     ErrorCode = "INTERNAL_ERROR"
)

type AppError struct {
    Code       ErrorCode   `json:"code"`
    Message    string      `json:"message"`
    Details    interface{} `json:"details,omitempty"`
    StatusCode int         `json:"-"`
}

func (e *AppError) Error() string {
    return e.Message
}

// Usage in services
func (s *BookingService) CreateBooking(ctx context.Context, req *dto.CreateBookingRequest) (*models.Booking, error) {
    // Validate slot availability
    slot, err := s.slotRepo.GetByID(ctx, req.SlotID)
    if err != nil {
        return nil, &errors.AppError{
            Code:       errors.ErrCodeNotFound,
            Message:    "Slot not found",
            StatusCode: http.StatusNotFound,
        }
    }

    if slot.IsBooked {
        return nil, &errors.AppError{
            Code:       errors.ErrCodeConflict,
            Message:    "Slot is no longer available",
            StatusCode: http.StatusConflict,
        }
    }

    // Create booking with transaction
    booking := &models.Booking{...}
    if err := s.bookingRepo.Create(ctx, booking); err != nil {
        return nil, errors.Internal("Failed to create booking", err)
    }

    return booking, nil
}
```

### 5. Validation Layer

Separate validation from business logic:

```go
// internal/validators/booking_validator.go
type BookingValidator struct {
    validator *validator.Validate
}

func (v *BookingValidator) ValidateCreateRequest(req *dto.CreateBookingRequest) error {
    if err := v.validator.Struct(req); err != nil {
        return errors.ValidationError(err)
    }

    // Custom business rules
    if req.StartTime.Before(time.Now().Add(1 * time.Hour)) {
        return errors.ValidationError(fmt.Errorf("booking must be at least 1 hour in advance"))
    }

    if !isValidPhone(req.Customer.Phone) {
        return errors.ValidationError(fmt.Errorf("invalid phone number format"))
    }

    return nil
}
```

### 6. Repository Pattern

Interface-based data access with transaction support:

```go
// internal/repository/booking_repository.go
type BookingRepository interface {
    Create(ctx context.Context, booking *models.Booking) error
    GetByID(ctx context.Context, id uuid.UUID) (*models.Booking, error)
    GetByBusinessID(ctx context.Context, businessID uuid.UUID, opts ...QueryOption) ([]*models.Booking, error)
    Update(ctx context.Context, booking *models.Booking) error
    UpdateStatus(ctx context.Context, id uuid.UUID, status models.BookingStatus) error
    Delete(ctx context.Context, id uuid.UUID) error
    WithTx(tx *gorm.DB) BookingRepository
}

type bookingRepository struct {
    db *gorm.DB
}

func NewBookingRepository(db *gorm.DB) BookingRepository {
    return &bookingRepository{db: db}
}

func (r *bookingRepository) WithTx(tx *gorm.DB) BookingRepository {
    return &bookingRepository{db: tx}
}

func (r *bookingRepository) Create(ctx context.Context, booking *models.Booking) error {
    return r.db.WithContext(ctx).Create(booking).Error
}
```

---

## Example Flow: Create Booking

### Current (Broken)

```go
// handlers.go - Direct DB access, no validation
func (h *Handler) CreateBooking(c *gin.Context) {
    var booking models.Booking
    if err := c.ShouldBindJSON(&booking); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }

    // ❌ No transaction - race condition
    h.Repo.DB.Create(&booking)                              // Step 1
    h.Repo.DB.Model(&models.Slot{}).Where(...).Update(...)    // Step 2 (separate!)

    c.JSON(201, booking)
}
```

### Recommended (Correct)

```go
// controllers/booking_controller.go
type BookingController struct {
    service *services.BookingService
}

func (ctrl *BookingController) CreateBooking(c *gin.Context) {
    var req dto.CreateBookingRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, errors.ValidationResponse(err))
        return
    }

    // Extract context (with user from auth middleware)
    ctx := c.Request.Context()

    // Call service layer
    booking, err := ctrl.service.CreateBooking(ctx, &req)
    if err != nil {
        errors.HandleError(c, err)
        return
    }

    c.JSON(201, dto.BookingToResponse(booking))
}

// services/booking_service.go
func (s *BookingService) CreateBooking(ctx context.Context, req *dto.CreateBookingRequest) (*models.Booking, error) {
    // Validate request
    if err := s.validator.ValidateCreateRequest(req); err != nil {
        return nil, err
    }

    // Get business from context (tenant)
    businessID := types.GetBusinessIDFromContext(ctx)

    // Validate slot exists and is available
    slot, err := s.slotRepo.GetByID(ctx, req.SlotID)
    if err != nil {
        return nil, errors.NotFound("Slot not found", err)
    }

    if slot.BusinessID != businessID {
        return nil, errors.Forbidden("Slot belongs to different business")
    }

    if slot.IsBooked {
        return nil, errors.Conflict("Slot is no longer available")
    }

    // Validate service exists and belongs to business
    service, err := s.serviceRepo.GetByID(ctx, req.ServiceID)
    if err != nil {
        return nil, errors.NotFound("Service not found", err)
    }

    if service.BusinessID != businessID {
        return nil, errors.Forbidden("Service belongs to different business")
    }

    // Validate pricing
    if req.TotalPrice != service.TotalPrice {
        return nil, errors.Validation(fmt.Errorf("price mismatch"))
    }

    // Get or create customer
    customer, err := s.customerRepo.GetOrCreate(ctx, &req.Customer)
    if err != nil {
        return nil, errors.Internal("Failed to get/create customer", err)
    }

    // ✅ Transaction - all or nothing
    var booking *models.Booking
    err = s.repo.Transaction(func(tx *gorm.DB) error {
        // Create booking
        booking = &models.Booking{
            BusinessID: businessID,
            ServiceID:  req.ServiceID,
            SlotID:     req.SlotID,
            CustomerID: customer.ID,
            Status:      models.BookingStatusPending,
            // ... other fields
        }

        if err := s.bookingRepo.WithTx(tx).Create(ctx, booking); err != nil {
            return errors.Internal("Failed to create booking", err)
        }

        // Mark slot as booked (in same transaction)
        if err := s.slotRepo.WithTx(tx).MarkBooked(ctx, req.SlotID); err != nil {
            return errors.Internal("Failed to mark slot as booked", err)
        }

        // Create audit log
        history := &models.BookingHistory{
            BookingID:   booking.ID,
            Action:      "created",
            Previous:    "",
            Current:     bookingToString(booking),
            PerformedBy: types.GetUserIDFromContext(ctx),
        }
        return s.historyRepo.WithTx(tx).Create(ctx, history)
    })

    if err != nil {
        return nil, err
    }

    // Send confirmation email (async)
    s.notificationService.SendBookingConfirmation(ctx, booking)

    return booking, nil
}

// repository/repository.go
func (r *Repository) Transaction(fn func(tx *gorm.DB) error) error {
    return r.db.Transaction(fn)
}
```

---

## Middleware Chain Implementation

```go
// api/router.go
func NewRouter(
    authCtrl *controllers.AuthController,
    businessCtrl *controllers.BusinessController,
    // ... other controllers
) *gin.Engine {
    r := gin.New()

    // Global middleware (applied to all routes)
    r.Use(
        middleware.RequestID(),           // Add correlation ID
        middleware.SecurityHeaders(),      // Security headers
        middleware.CORS(),              // CORS configuration
        middleware.Logger(),             // Request logging
        middleware.Recovery(),           // Panic recovery
    )

    // Public routes (no auth required)
    public := r.Group("/api/v1")
    {
        auth := public.Group("/auth")
        {
            auth.POST("/register", authCtrl.Register)
            auth.POST("/login", authCtrl.Login)
        }

        // Public booking endpoints (for customers)
        public.GET("/businesses/:slug", businessCtrl.GetBySlug) // Public access
        public.GET("/businesses/:slug/services", serviceCtrl.ListByBusinessSlug)
        public.GET("/businesses/:slug/slots", slotCtrl.ListAvailable)
        public.POST("/bookings", bookingCtrl.Create) // Public booking
    }

    // Protected routes (require authentication)
    protected := r.Group("/api/v1")
    protected.Use(middleware.Auth()) // JWT verification
    {
        // User profile
        protected.GET("/auth/me", authCtrl.Me)
        protected.POST("/auth/logout", authCtrl.Logout)
        protected.POST("/auth/refresh", authCtrl.Refresh)

        // Business management (owner/admin only)
        business := protected.Group("/businesses/:id")
        business.Use(middleware.Authorize("owner", "admin")) // RBAC
        {
            business.PATCH("", businessCtrl.Update)
            business.DELETE("", businessCtrl.Delete)
            business.GET("/settings", businessCtrl.GetSettings)
            business.PATCH("/settings", businessCtrl.UpdateSettings)
        }

        // Service management (owner/admin only)
        services := protected.Group("/businesses/:id/services")
        services.Use(middleware.Authorize("owner", "admin"))
        {
            services.POST("", serviceCtrl.Create)
            services.GET("", serviceCtrl.ListByBusiness)
            services.GET("/:service_id", serviceCtrl.GetByID)
            services.PATCH("/:service_id", serviceCtrl.Update)
            services.DELETE("/:service_id", serviceCtrl.Delete)
        }

        // Booking management (owner/admin/staff)
        bookings := protected.Group("/businesses/:id/bookings")
        bookings.Use(middleware.Authorize("owner", "admin", "staff"))
        {
            bookings.GET("", bookingCtrl.ListByBusiness)
        }

        // Protected booking operations
        protected.GET("/bookings/:id", bookingCtrl.GetByID)
        protected.PATCH("/bookings/:id", bookingCtrl.Update)
        protected.DELETE("/bookings/:id", bookingCtrl.Cancel)
        protected.POST("/bookings/:id/confirm", bookingCtrl.Confirm)
        protected.POST("/bookings/:id/complete", bookingCtrl.Complete)
        protected.POST("/bookings/:id/reschedule", bookingCtrl.Reschedule)
    }

    // Admin routes (superadmin only)
    admin := r.Group("/admin")
    admin.Use(middleware.Auth())
    admin.Use(middleware.RequireSuperAdmin())
    {
        admin.GET("/businesses", adminCtrl.ListBusinesses)
        admin.POST("/businesses/:id/suspend", adminCtrl.SuspendBusiness)
        admin.GET("/users", adminCtrl.ListUsers)
    }

    return r
}
```

---

## Testing Strategy

```
test/
├── unit/
│   ├── services/
│   │   ├── booking_service_test.go
│   │   ├── service_service_test.go
│   │   └── ...
│   ├── validators/
│   │   └── booking_validator_test.go
│   └── utils/
│       └── jwt_test.go
│
├── integration/
│   ├── handlers/
│   │   ├── auth_handler_test.go
│   │   └── booking_handler_test.go
│   └── repository/
│       └── booking_repository_test.go
│
└── e2e/
    └── booking_flow_test.go
```

**Unit Tests:** Test business logic in isolation (mock repositories)
**Integration Tests:** Test handlers and services together (real DB)
**E2E Tests:** Test complete user flows (HTTP requests)

---

## Summary

### Key Benefits

1. **Separation of Concerns**
   - Controllers handle HTTP
   - Services handle business logic
   - Repositories handle data access

2. **Testability**
   - Services can be unit tested with mocks
   - Controllers can be integration tested
   - Easy to mock dependencies

3. **Maintainability**
   - Changes in one layer don't affect others
   - Clear responsibilities
   - Easy to locate bugs

4. **Scalability**
   - Easy to add new features
   - Can swap implementations (e.g., different cache)
   - Clear extension points

5. **Type Safety**
   - Compile-time checks
   - IDE support
   - Refactoring confidence

### Migration Path

**Week 1-2:**
- Create directory structure
- Move existing handlers to controllers
- Create service layer stubs

**Week 3-4:**
- Implement middleware
- Implement auth service
- Update controllers to use services

**Week 5-6:**
- Implement validation layer
- Implement error handling
- Add tests

**Week 7+:**
- Implement missing endpoints
- Payment integration
- Analytics

---

**Total Refactoring Effort:** 6-8 weeks for full migration
