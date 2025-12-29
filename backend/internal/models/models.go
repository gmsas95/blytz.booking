package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BookingStatus string
type UserRole string
type PaymentStatus string
type SubscriptionPlan string
type SubscriptionStatus string

const (
	BookingStatusPending   BookingStatus = "PENDING"
	BookingStatusConfirmed BookingStatus = "CONFIRMED"
	BookingStatusCompleted BookingStatus = "COMPLETED"
	BookingStatusCancelled BookingStatus = "CANCELLED"
	BookingStatusNoShow    BookingStatus = "NO_SHOW"
)

const (
	UserRoleOwner      UserRole = "owner"
	UserRoleAdmin      UserRole = "admin"
	UserRoleStaff      UserRole = "staff"
	UserRoleSuperAdmin UserRole = "superadmin"
)

const (
	PaymentStatusPending   PaymentStatus = "pending"
	PaymentStatusPartial   PaymentStatus = "partial"
	PaymentStatusSucceeded PaymentStatus = "succeeded"
	PaymentStatusFailed    PaymentStatus = "failed"
	PaymentStatusRefunded  PaymentStatus = "refunded"
)

const (
	SubscriptionPlanFree       SubscriptionPlan = "free"
	SubscriptionPlanPro        SubscriptionPlan = "pro"
	SubscriptionPlanEnterprise SubscriptionPlan = "enterprise"
)

const (
	SubscriptionStatusActive    SubscriptionStatus = "active"
	SubscriptionStatusTrial     SubscriptionStatus = "trial"
	SubscriptionStatusPastDue   SubscriptionStatus = "past_due"
	SubscriptionStatusCancelled SubscriptionStatus = "cancelled"
)

type Business struct {
	ID          uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Name        string     `json:"name" gorm:"not null"`
	Slug        string     `json:"slug" gorm:"uniqueIndex;not null"`
	Vertical    string     `json:"vertical" gorm:"not null"`
	Description string     `json:"description"`
	ThemeColor  string     `json:"theme_color" gorm:"default:'blue'"`
	LogoURL     *string    `json:"logo_url" gorm:"type:text"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty" gorm:"index"`

	Users     []User            `json:"users,omitempty" gorm:"foreignKey:BusinessID"`
	Services  []Service         `json:"services,omitempty" gorm:"foreignKey:BusinessID"`
	Slots     []Slot            `json:"slots,omitempty" gorm:"foreignKey:BusinessID"`
	Bookings  []Booking         `json:"bookings,omitempty" gorm:"foreignKey:BusinessID"`
	Settings  *BusinessSettings `json:"settings,omitempty" gorm:"foreignKey:BusinessID"`
	Customers []Customer        `json:"customers,omitempty" gorm:"foreignKey:BusinessID"`
}

type BusinessSettings struct {
	ID                      uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	BusinessID              uuid.UUID `json:"business_id" gorm:"type:uuid;not null;uniqueIndex"`
	RequireDeposit          bool      `json:"require_deposit" gorm:"default:false"`
	CancellationPolicyHours int       `json:"cancellation_policy_hours" gorm:"default:24"`
	ConfirmationEmail       bool      `json:"confirmation_email" gorm:"default:true"`
	ReminderEmail           bool      `json:"reminder_email" gorm:"default:true"`
	ReminderHoursBefore     int       `json:"reminder_hours_before" gorm:"default:24"`
	Timezone                string    `json:"timezone" gorm:"default:'UTC'"`
	BookingBufferMinutes    int       `json:"booking_buffer_minutes" gorm:"default:0"`
	CreatedAt               time.Time `json:"created_at"`
	UpdatedAt               time.Time `json:"updated_at"`

	Business Business `json:"-" gorm:"foreignKey:BusinessID"`
}

type Service struct {
	ID            uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	BusinessID    uuid.UUID  `json:"business_id" gorm:"type:uuid;not null;index"`
	Name          string     `json:"name" gorm:"not null"`
	Description   string     `json:"description"`
	DurationMin   int        `json:"duration_min" gorm:"not null"`
	TotalPrice    float64    `json:"total_price" gorm:"not null"`
	DepositAmount float64    `json:"deposit_amount" gorm:"not null;default:0"`
	IsActive      bool       `json:"is_active" gorm:"default:true"`
	MaxCapacity   int        `json:"max_capacity" gorm:"default:1"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	DeletedAt     *time.Time `json:"deleted_at,omitempty" gorm:"index"`

	Business Business `json:"business,omitempty" gorm:"foreignKey:BusinessID"`
}

type Slot struct {
	ID          uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	BusinessID  uuid.UUID  `json:"business_id" gorm:"type:uuid;not null;index"`
	ServiceID   *uuid.UUID `json:"service_id" gorm:"type:uuid;index"`
	StartTime   time.Time  `json:"start_time" gorm:"not null;index"`
	EndTime     time.Time  `json:"end_time" gorm:"not null"`
	IsBooked    bool       `json:"is_booked" gorm:"default:false"`
	Capacity    int        `json:"capacity" gorm:"default:1"`
	BookedCount int        `json:"booked_count" gorm:"default:0"`
	CreatedBy   *uuid.UUID `json:"created_by,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty" gorm:"index"`

	Business Business `json:"business,omitempty" gorm:"foreignKey:BusinessID"`
	Service  *Service `json:"service,omitempty" gorm:"foreignKey:ServiceID"`
}

type Customer struct {
	ID            uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	BusinessID    uuid.UUID  `json:"business_id" gorm:"type:uuid;not null;index"`
	Name          string     `json:"name" gorm:"not null"`
	Email         string     `json:"email" gorm:"not null;index"`
	Phone         string     `json:"phone" gorm:"not null"`
	Notes         string     `json:"notes"`
	TotalBookings int        `json:"total_bookings" gorm:"default:0"`
	TotalSpent    float64    `json:"total_spent" gorm:"default:0"`
	LastBookingAt *time.Time `json:"last_booking_at,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	DeletedAt     *time.Time `json:"deleted_at,omitempty" gorm:"index"`
}
type Booking struct {
	ID              uuid.UUID       `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	BusinessID      uuid.UUID       `json:"business_id" gorm:"type:uuid;not null;index"`
	ServiceID       uuid.UUID       `json:"service_id" gorm:"type:uuid;not null;index"`
	SlotID          uuid.UUID       `json:"slot_id" gorm:"type:uuid;not null;index"`
	CustomerID      *uuid.UUID      `json:"customer_id,omitempty" gorm:"type:uuid;index"`
	ServiceName     string          `json:"service_name" gorm:"not null"`
	SlotTime        time.Time       `json:"slot_time" gorm:"not null"`
	CustomerDetails CustomerDetails `json:"customer" gorm:"embedded"`
	Status          BookingStatus   `json:"status" gorm:"not null;default:'PENDING'"`
	DepositPaid     float64         `json:"deposit_paid" gorm:"not null"`
	TotalPrice      float64         `json:"total_price" gorm:"not null"`
	PaymentStatus   PaymentStatus   `json:"payment_status" gorm:"default:'pending'"`
	PaymentID       *uuid.UUID      `json:"payment_id,omitempty"`
	Notes           string          `json:"notes"`
	CancelledAt     *time.Time      `json:"cancelled_at,omitempty"`
	CancelReason    string          `json:"cancel_reason,omitempty"`
	CancelledBy     *uuid.UUID      `json:"cancelled_by,omitempty"`
	NoShowAt        *time.Time      `json:"no_show_at,omitempty"`
	CreatedAt       time.Time       `json:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at"`
	DeletedAt       *time.Time      `json:"deleted_at,omitempty" gorm:"index"`

	Business Business         `json:"business,omitempty" gorm:"foreignKey:BusinessID"`
	Service  Service          `json:"service,omitempty" gorm:"foreignKey:ServiceID"`
	Slot     Slot             `json:"slot,omitempty" gorm:"foreignKey:SlotID"`
	Customer *Customer        `json:"customer_obj,omitempty" gorm:"foreignKey:CustomerID"`
	History  []BookingHistory `json:"history,omitempty" gorm:"foreignKey:BookingID"`
}

type BookingHistory struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	BookingID   uuid.UUID `json:"booking_id" gorm:"type:uuid;not null;index"`
	Action      string    `json:"action"`
	Previous    string    `json:"previous"`
	Current     string    `json:"current"`
	PerformedBy uuid.UUID `json:"performed_by" gorm:"type:uuid;index"`
	Timestamp   time.Time `json:"timestamp"`

	Booking Booking `json:"-" gorm:"foreignKey:BookingID"`
}

type User struct {
	ID            uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Email         string     `json:"email" gorm:"uniqueIndex;not null"`
	PasswordHash  string     `json:"-" gorm:"not null"`
	Name          string     `json:"name"`
	Role          UserRole   `json:"role" gorm:"not null;default:'owner'"`
	BusinessID    *uuid.UUID `json:"business_id,omitempty" gorm:"type:uuid;index"`
	IsActive      bool       `json:"is_active" gorm:"default:true"`
	EmailVerified bool       `json:"email_verified" gorm:"default:false"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	DeletedAt     *time.Time `json:"deleted_at,omitempty" gorm:"index"`

	Business *Business `json:"business,omitempty" gorm:"foreignKey:BusinessID"`
}

type Payment struct {
	ID                    uuid.UUID     `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	BookingID             uuid.UUID     `json:"booking_id" gorm:"type:uuid;not null;uniqueIndex"`
	Amount                float64       `json:"amount" gorm:"not null"`
	Currency              string        `json:"currency" gorm:"default:'USD'"`
	Status                PaymentStatus `json:"status" gorm:"not null;default:'pending'"`
	PaymentMethod         string        `json:"payment_method"`
	StripePaymentIntentID string        `json:"stripe_payment_intent_id" gorm:"index"`
	StripeReceiptURL      string        `json:"stripe_receipt_url"`
	CreatedAt             time.Time     `json:"created_at"`
	UpdatedAt             time.Time     `json:"updated_at"`
}

type Subscription struct {
	ID                   uuid.UUID          `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	BusinessID           uuid.UUID          `json:"business_id" gorm:"type:uuid;not null;uniqueIndex"`
	Plan                 SubscriptionPlan   `json:"plan" gorm:"not null;default:'free'"`
	Status               SubscriptionStatus `json:"status" gorm:"not null;default:'trial'"`
	StripeCustomerID     string             `json:"stripe_customer_id"`
	StripeSubscriptionID string             `json:"stripe_subscription_id"`
	CurrentPeriodStart   time.Time          `json:"current_period_start"`
	CurrentPeriodEnd     time.Time          `json:"current_period_end"`
	CancelAtPeriodEnd    bool               `json:"cancel_at_period_end" gorm:"default:false"`
	CreatedAt            time.Time          `json:"created_at"`
	UpdatedAt            time.Time          `json:"updated_at"`

	Business Business `json:"-" gorm:"foreignKey:BusinessID"`
}

type RecurringSchedule struct {
	ID           uuid.UUID   `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	BusinessID   uuid.UUID   `json:"business_id" gorm:"type:uuid;not null;index"`
	Name         string      `json:"name" gorm:"not null"`
	DaysOfWeek   []int       `json:"days_of_week" gorm:"type:integer[];not null"`
	StartTime    string      `json:"start_time" gorm:"not null"`
	EndTime      string      `json:"end_time" gorm:"not null"`
	StartDate    time.Time   `json:"start_date" gorm:"not null"`
	EndDate      time.Time   `json:"end_date" gorm:"not null"`
	ExcludeDates []time.Time `json:"exclude_dates" gorm:"type:date[]"`
	IsActive     bool        `json:"is_active" gorm:"default:true"`
	CreatedAt    time.Time   `json:"created_at"`
	UpdatedAt    time.Time   `json:"updated_at"`
	DeletedAt    *time.Time  `json:"deleted_at,omitempty" gorm:"index"`

	Business Business `json:"-" gorm:"foreignKey:BusinessID"`
}

type Webhook struct {
	ID         uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	BusinessID uuid.UUID  `json:"business_id" gorm:"type:uuid;not null;index"`
	URL        string     `json:"url" gorm:"not null"`
	Events     []string   `json:"events" gorm:"type:text[];not null"`
	Secret     string     `json:"secret" gorm:"not null"`
	IsActive   bool       `json:"is_active" gorm:"default:true"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	DeletedAt  *time.Time `json:"deleted_at,omitempty" gorm:"index"`

	Business Business `json:"-" gorm:"foreignKey:BusinessID"`
}

type RefreshToken struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID    uuid.UUID `json:"user_id" gorm:"type:uuid;not null;index"`
	Token     string    `json:"token" gorm:"not null;uniqueIndex"`
	ExpiresAt time.Time `json:"expires_at" gorm:"not null;index"`
	CreatedAt time.Time `json:"created_at"`

	User User `json:"-" gorm:"foreignKey:UserID"`
}

func (b *Business) BeforeCreate(tx *gorm.DB) error {
	if b.ID == uuid.Nil {
		b.ID = uuid.New()
	}
	return nil
}

func (s *Service) BeforeCreate(tx *gorm.DB) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return nil
}

func (sl *Slot) BeforeCreate(tx *gorm.DB) error {
	if sl.ID == uuid.Nil {
		sl.ID = uuid.New()
	}
	return nil
}

func (b *Booking) BeforeCreate(tx *gorm.DB) error {
	if b.ID == uuid.Nil {
		b.ID = uuid.New()
	}
	return nil
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}

func (c *Customer) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return nil
}

func (p *Payment) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}

func (sub *Subscription) BeforeCreate(tx *gorm.DB) error {
	if sub.ID == uuid.Nil {
		sub.ID = uuid.New()
	}
	return nil
}

func (rs *RecurringSchedule) BeforeCreate(tx *gorm.DB) error {
	if rs.ID == uuid.Nil {
		rs.ID = uuid.New()
	}
	return nil
}

func (w *Webhook) BeforeCreate(tx *gorm.DB) error {
	if w.ID == uuid.Nil {
		w.ID = uuid.New()
	}
	return nil
}

func (rt *RefreshToken) BeforeCreate(tx *gorm.DB) error {
	if rt.ID == uuid.Nil {
		rt.ID = uuid.New()
	}
	return nil
}

func (rt *RefreshToken) IsExpired() bool {
	return time.Now().After(rt.ExpiresAt)
}

type CustomerDetails struct {
	Name  string `json:"name" gorm:"not null"`
	Email string `json:"email" gorm:"not null"`
	Phone string `json:"phone" gorm:"not null"`
}
