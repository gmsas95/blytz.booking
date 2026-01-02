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
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
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
	Slug        *string `json:"slug"`
	Vertical    *string `json:"vertical"`
	Description *string `json:"description"`
	ThemeColor  *string `json:"theme_color"`
}

// Service DTOs

type ServiceResponse struct {
	ID            string  `json:"id"`
	BusinessID    string  `json:"business_id"`
	Name          string  `json:"name"`
	Description   string  `json:"description"`
	DurationMin   int     `json:"duration_min"`
	TotalPrice    float64 `json:"total_price"`
	DepositAmount float64 `json:"deposit_amount"`
	CreatedAt     string  `json:"created_at"`
	UpdatedAt     string  `json:"updated_at"`
}

type CreateServiceRequest struct {
	BusinessID    string  `json:"business_id" binding:"required,uuid"`
	Name          string  `json:"name" binding:"required"`
	Description   string  `json:"description"`
	DurationMin   int     `json:"duration_min" binding:"required,min=1"`
	TotalPrice    float64 `json:"total_price" binding:"required,gt=0"`
	DepositAmount float64 `json:"deposit_amount" binding:"required,gte=0"`
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
	ID          string          `json:"id"`
	BusinessID  string          `json:"business_id"`
	ServiceID   string          `json:"service_id"`
	SlotID      string          `json:"slot_id"`
	ServiceName string          `json:"service_name"`
	SlotTime    string          `json:"slot_time"`
	Customer    CustomerDetails `json:"customer"`
	Status      string          `json:"status"`
	DepositPaid float64         `json:"deposit_paid"`
	TotalPrice  float64         `json:"total_price"`
	CreatedAt   string          `json:"created_at"`
	UpdatedAt   string          `json:"updated_at"`
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

// Error Response DTO

type ErrorResponse struct {
	Error   string `json:"error"`
	Details string `json:"details,omitempty"`
}
