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
	CreatedAt string `json:"createdAt"`
}

// Business DTOs

type BusinessResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Vertical    string `json:"vertical"`
	Description string `json:"description"`
	ThemeColor  string `json:"themeColor"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
}

type CreateBusinessRequest struct {
	Name        string `json:"name" binding:"required"`
	Slug        string `json:"slug" binding:"required"`
	Vertical    string `json:"vertical" binding:"required"`
	Description string `json:"description"`
	ThemeColor  string `json:"themeColor"`
}

type UpdateBusinessRequest struct {
	Name        *string `json:"name"`
	Slug        *string `json:"slug"`
	Vertical    *string `json:"vertical"`
	Description *string `json:"description"`
	ThemeColor  *string `json:"themeColor"`
}

// Service DTOs

type ServiceResponse struct {
	ID            string  `json:"id"`
	BusinessID    string  `json:"businessId"`
	Name          string  `json:"name"`
	Description   string  `json:"description"`
	DurationMin   int     `json:"durationMin"`
	TotalPrice    float64 `json:"totalPrice"`
	DepositAmount float64 `json:"depositAmount"`
	CreatedAt     string  `json:"createdAt"`
	UpdatedAt     string  `json:"updatedAt"`
}

type CreateServiceRequest struct {
	Name          string  `json:"name" binding:"required"`
	Description   string  `json:"description"`
	DurationMin   int     `json:"durationMin" binding:"required,min=1"`
	TotalPrice    float64 `json:"totalPrice" binding:"required,gt=0"`
	DepositAmount float64 `json:"depositAmount" binding:"required,gte=0"`
}

type UpdateServiceRequest struct {
	Name          *string  `json:"name"`
	Description   *string  `json:"description"`
	DurationMin   *int     `json:"durationMin"`
	TotalPrice    *float64 `json:"totalPrice"`
	DepositAmount *float64 `json:"depositAmount"`
}

// Slot DTOs

type SlotResponse struct {
	ID         string `json:"id"`
	BusinessID string `json:"businessId"`
	StartTime  string `json:"startTime"`
	EndTime    string `json:"endTime"`
	IsBooked   bool   `json:"isBooked"`
	CreatedAt  string `json:"createdAt"`
	UpdatedAt  string `json:"updatedAt"`
}

type CreateSlotRequest struct {
	BusinessID string `json:"businessId" binding:"required,uuid"`
	StartTime  string `json:"startTime" binding:"required"`
	EndTime    string `json:"endTime" binding:"required,gtfield=StartTime"`
}

// Booking DTOs

type CustomerDetails struct {
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required,email"`
	Phone string `json:"phone" binding:"required"`
}

type BookingResponse struct {
	ID          string          `json:"id"`
	BusinessID  string          `json:"businessId"`
	ServiceID   string          `json:"serviceId"`
	SlotID      string          `json:"slotId"`
	ServiceName string          `json:"serviceName"`
	SlotTime    string          `json:"slotTime"`
	Customer    CustomerDetails `json:"customer"`
	Status      string          `json:"status"`
	DepositPaid float64         `json:"depositPaid"`
	TotalPrice  float64         `json:"totalPrice"`
	CreatedAt   string          `json:"createdAt"`
	UpdatedAt   string          `json:"updatedAt"`
}

type CreateBookingRequest struct {
	BusinessID string          `json:"businessId" binding:"required,uuid"`
	ServiceID  string          `json:"serviceId" binding:"required,uuid"`
	SlotID     string          `json:"slotId" binding:"required,uuid"`
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
