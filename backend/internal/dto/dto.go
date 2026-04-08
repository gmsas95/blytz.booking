package dto

// Auth DTOs

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Name     string `json:"name" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type AuthResponse struct {
	User UserResponse `json:"user"`
}

type MembershipBusinessResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Vertical    string `json:"vertical"`
	Description string `json:"description"`
	ThemeColor  string `json:"theme_color"`
}

type MembershipResponse struct {
	ID         string                     `json:"id"`
	BusinessID string                     `json:"business_id"`
	Role       string                     `json:"role"`
	Business   MembershipBusinessResponse `json:"business"`
}

type CurrentUserResponse struct {
	User             UserResponse         `json:"user"`
	Memberships      []MembershipResponse `json:"memberships"`
	ActiveBusinessID string               `json:"active_business_id,omitempty"`
}

type UserResponse struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
}

// Business DTOs

type BusinessResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Vertical    string `json:"vertical"`
	Description string `json:"description"`
	ThemeColor  string `json:"theme_color"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

type CreateBusinessRequest struct {
	Name        string `json:"name" binding:"required"`
	Slug        string `json:"slug" binding:"required"`
	Vertical    string `json:"vertical" binding:"required"`
	Description string `json:"description"`
	ThemeColor  string `json:"theme_color"`
}

type UpdateBusinessRequest struct {
	Name        *string `json:"name"`
	Vertical    *string `json:"vertical"`
	Description *string `json:"description"`
	ThemeColor  *string `json:"theme_color"`
}

// Service DTOs

type ServiceResponse struct {
	ID                 string `json:"id"`
	BusinessID         string `json:"business_id"`
	Name               string `json:"name"`
	Description        string `json:"description"`
	DurationMin        int    `json:"duration_min"`
	TotalPriceMinor    int64  `json:"total_price_minor"`
	DepositAmountMinor int64  `json:"deposit_amount_minor"`
	CurrencyCode       string `json:"currency_code"`
	CreatedAt          string `json:"created_at"`
	UpdatedAt          string `json:"updated_at"`
}

type CreateServiceRequest struct {
	BusinessID         string `json:"business_id" binding:"required,uuid"`
	Name               string `json:"name" binding:"required"`
	Description        string `json:"description"`
	DurationMin        int    `json:"duration_min" binding:"required,min=1"`
	TotalPriceMinor    int64  `json:"total_price_minor" binding:"required,gt=0"`
	DepositAmountMinor int64  `json:"deposit_amount_minor" binding:"required,gte=0"`
	CurrencyCode       string `json:"currency_code" binding:"required,len=3"`
}

// Slot DTOs

type SlotResponse struct {
	ID         string `json:"id"`
	BusinessID string `json:"business_id"`
	StartTime  string `json:"start_time"`
	EndTime    string `json:"end_time"`
	IsBooked   bool   `json:"is_booked"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}

type CreateSlotRequest struct {
	BusinessID string `json:"business_id" binding:"required,uuid"`
	StartTime  string `json:"start_time" binding:"required"`
	EndTime    string `json:"end_time" binding:"required,gtfield=StartTime"`
}

// Booking DTOs

type CustomerDetails struct {
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required,email"`
	Phone string `json:"phone" binding:"required"`
}

type BookingResponse struct {
	ID               string          `json:"id"`
	BusinessID       string          `json:"business_id"`
	ServiceID        string          `json:"service_id"`
	SlotID           string          `json:"slot_id"`
	ServiceName      string          `json:"service_name"`
	SlotTime         string          `json:"slot_time"`
	Customer         CustomerDetails `json:"customer"`
	Status           string          `json:"status"`
	DepositPaidMinor int64           `json:"deposit_paid_minor"`
	TotalPriceMinor  int64           `json:"total_price_minor"`
	CurrencyCode     string          `json:"currency_code"`
	CreatedAt        string          `json:"created_at"`
	UpdatedAt        string          `json:"updated_at"`
}

type CreateBookingRequest struct {
	BusinessID string          `json:"business_id" binding:"required,uuid"`
	ServiceID  string          `json:"service_id" binding:"required,uuid"`
	SlotID     string          `json:"slot_id" binding:"required,uuid"`
	Customer   CustomerDetails `json:"customer" binding:"required"`
}

type UpdateBookingStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=PENDING CONFIRMED COMPLETED CANCELLED"`
}

// Customer DTOs

type CustomerResponse struct {
	ID         string `json:"id"`
	BusinessID string `json:"business_id"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	Phone      string `json:"phone"`
	Notes      string `json:"notes"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}

type CreateCustomerRequest struct {
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required,email"`
	Phone string `json:"phone" binding:"required"`
	Notes string `json:"notes"`
}

// Vehicle DTOs

type VehicleResponse struct {
	ID           string           `json:"id"`
	BusinessID   string           `json:"business_id"`
	CustomerID   string           `json:"customer_id"`
	Year         int              `json:"year"`
	Make         string           `json:"make"`
	Model        string           `json:"model"`
	Color        string           `json:"color"`
	LicensePlate string           `json:"license_plate"`
	Customer     CustomerResponse `json:"customer"`
	CreatedAt    string           `json:"created_at"`
	UpdatedAt    string           `json:"updated_at"`
}

type CreateVehicleRequest struct {
	CustomerID   string `json:"customer_id" binding:"required,uuid"`
	Year         int    `json:"year" binding:"required,min=1900,max=2100"`
	Make         string `json:"make" binding:"required"`
	Model        string `json:"model" binding:"required"`
	Color        string `json:"color"`
	LicensePlate string `json:"license_plate"`
}

// Job DTOs

type JobResponse struct {
	ID          string           `json:"id"`
	BusinessID  string           `json:"business_id"`
	CustomerID  string           `json:"customer_id"`
	VehicleID   string           `json:"vehicle_id"`
	BookingID   string           `json:"booking_id,omitempty"`
	Title       string           `json:"title"`
	Status      string           `json:"status"`
	ScheduledAt string           `json:"scheduled_at"`
	Notes       string           `json:"notes"`
	Customer    CustomerResponse `json:"customer"`
	Vehicle     VehicleResponse  `json:"vehicle"`
	CreatedAt   string           `json:"created_at"`
	UpdatedAt   string           `json:"updated_at"`
}

type CreateJobRequest struct {
	CustomerID  string `json:"customer_id" binding:"required,uuid"`
	VehicleID   string `json:"vehicle_id" binding:"required,uuid"`
	BookingID   string `json:"booking_id,omitempty" binding:"omitempty,uuid"`
	Title       string `json:"title" binding:"required"`
	Status      string `json:"status" binding:"omitempty,oneof=SCHEDULED IN_PROGRESS READY DELIVERED"`
	ScheduledAt string `json:"scheduled_at" binding:"required"`
	Notes       string `json:"notes"`
}

// Error Response DTO

type ErrorResponse struct {
	Error   string `json:"error"`
	Details string `json:"details,omitempty"`
}
