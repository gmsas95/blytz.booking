package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BookingStatus string
type MembershipRole string
type JobStatus string

const (
	BookingStatusPending   BookingStatus = "PENDING"
	BookingStatusConfirmed BookingStatus = "CONFIRMED"
	BookingStatusCompleted BookingStatus = "COMPLETED"
	BookingStatusCancelled BookingStatus = "CANCELLED"

	MembershipRoleOwner MembershipRole = "OWNER"
	MembershipRoleStaff MembershipRole = "STAFF"

	JobStatusScheduled  JobStatus = "SCHEDULED"
	JobStatusInProgress JobStatus = "IN_PROGRESS"
	JobStatusReady      JobStatus = "READY"
	JobStatusDelivered  JobStatus = "DELIVERED"
)

type Business struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Name        string    `json:"name" gorm:"not null"`
	Slug        string    `json:"slug" gorm:"uniqueIndex;not null"`
	Vertical    string    `json:"vertical" gorm:"not null"`
	Description string    `json:"description"`
	ThemeColor  string    `json:"theme_color" gorm:"default:'blue'"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Service struct {
	ID                 uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	BusinessID         uuid.UUID `json:"business_id" gorm:"type:uuid;not null"`
	Name               string    `json:"name" gorm:"not null"`
	Description        string    `json:"description"`
	DurationMin        int       `json:"duration_min" gorm:"not null"`
	TotalPriceMinor    int64     `json:"total_price_minor" gorm:"not null;default:0"`
	DepositAmountMinor int64     `json:"deposit_amount_minor" gorm:"not null;default:0"`
	CurrencyCode       string    `json:"currency_code" gorm:"size:3;not null;default:'USD'"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
	Business           Business  `json:"business" gorm:"foreignKey:BusinessID"`
}

type Slot struct {
	ID         uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	BusinessID uuid.UUID `json:"business_id" gorm:"type:uuid;not null;index"`
	StartTime  time.Time `json:"start_time" gorm:"not null;index"`
	EndTime    time.Time `json:"end_time" gorm:"not null"`
	IsBooked   bool      `json:"is_booked" gorm:"default:false"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	Business   Business  `json:"business" gorm:"foreignKey:BusinessID"`
}

type CustomerDetails struct {
	Name  string `json:"name" gorm:"not null"`
	Email string `json:"email" gorm:"not null"`
	Phone string `json:"phone" gorm:"not null"`
}

type Booking struct {
	ID               uuid.UUID       `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	BusinessID       uuid.UUID       `json:"business_id" gorm:"type:uuid;not null;index"`
	ServiceID        uuid.UUID       `json:"service_id" gorm:"type:uuid;not null"`
	SlotID           uuid.UUID       `json:"slot_id" gorm:"type:uuid;not null;index"`
	ServiceName      string          `json:"service_name" gorm:"not null"`
	SlotTime         time.Time       `json:"slot_time" gorm:"not null"`
	Customer         CustomerDetails `json:"customer" gorm:"embedded"`
	Status           BookingStatus   `json:"status" gorm:"not null;default:'PENDING'"`
	DepositPaidMinor int64           `json:"deposit_paid_minor" gorm:"not null;default:0"`
	TotalPriceMinor  int64           `json:"total_price_minor" gorm:"not null;default:0"`
	CurrencyCode     string          `json:"currency_code" gorm:"size:3;not null;default:'USD'"`
	CreatedAt        time.Time       `json:"created_at"`
	UpdatedAt        time.Time       `json:"updated_at"`
	Business         Business        `json:"business" gorm:"foreignKey:BusinessID"`
	Service          Service         `json:"service" gorm:"foreignKey:ServiceID"`
	Slot             Slot            `json:"slot" gorm:"foreignKey:SlotID"`
}

// User model for operators
type User struct {
	ID           uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Email        string    `json:"email" gorm:"uniqueIndex;not null"`
	Name         string    `json:"name"`
	PasswordHash string    `json:"-" gorm:"not null"`
	TokenVersion int       `json:"-" gorm:"not null;default:1"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type Membership struct {
	ID         uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID     uuid.UUID      `json:"user_id" gorm:"type:uuid;not null;uniqueIndex:idx_user_business_membership"`
	BusinessID uuid.UUID      `json:"business_id" gorm:"type:uuid;not null;uniqueIndex:idx_user_business_membership;index"`
	Role       MembershipRole `json:"role" gorm:"not null;default:'STAFF'"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	User       User           `json:"user" gorm:"foreignKey:UserID"`
	Business   Business       `json:"business" gorm:"foreignKey:BusinessID"`
}

type Customer struct {
	ID         uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	BusinessID uuid.UUID `json:"business_id" gorm:"type:uuid;not null;index"`
	Name       string    `json:"name" gorm:"not null"`
	Email      string    `json:"email" gorm:"not null;index"`
	Phone      string    `json:"phone" gorm:"not null"`
	Notes      string    `json:"notes"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	Business   Business  `json:"business" gorm:"foreignKey:BusinessID"`
}

type Vehicle struct {
	ID           uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	BusinessID   uuid.UUID `json:"business_id" gorm:"type:uuid;not null;index"`
	CustomerID   uuid.UUID `json:"customer_id" gorm:"type:uuid;not null;index"`
	Year         int       `json:"year"`
	Make         string    `json:"make" gorm:"not null"`
	Model        string    `json:"model" gorm:"not null"`
	Color        string    `json:"color"`
	LicensePlate string    `json:"license_plate"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Business     Business  `json:"business" gorm:"foreignKey:BusinessID"`
	Customer     Customer  `json:"customer" gorm:"foreignKey:CustomerID"`
}

type Job struct {
	ID          uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	BusinessID  uuid.UUID  `json:"business_id" gorm:"type:uuid;not null;index"`
	CustomerID  uuid.UUID  `json:"customer_id" gorm:"type:uuid;not null;index"`
	VehicleID   uuid.UUID  `json:"vehicle_id" gorm:"type:uuid;not null;index"`
	BookingID   *uuid.UUID `json:"booking_id" gorm:"type:uuid;index"`
	Title       string     `json:"title" gorm:"not null"`
	Status      JobStatus  `json:"status" gorm:"not null;default:'SCHEDULED'"`
	ScheduledAt time.Time  `json:"scheduled_at" gorm:"not null;index"`
	Notes       string     `json:"notes"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	Business    Business   `json:"business" gorm:"foreignKey:BusinessID"`
	Customer    Customer   `json:"customer" gorm:"foreignKey:CustomerID"`
	Vehicle     Vehicle    `json:"vehicle" gorm:"foreignKey:VehicleID"`
	Booking     *Booking   `json:"booking,omitempty" gorm:"foreignKey:BookingID"`
}

// BeforeCreate hook for GORM
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

func (m *Membership) BeforeCreate(tx *gorm.DB) error {
	if m.ID == uuid.Nil {
		m.ID = uuid.New()
	}
	return nil
}

func (c *Customer) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return nil
}

func (v *Vehicle) BeforeCreate(tx *gorm.DB) error {
	if v.ID == uuid.Nil {
		v.ID = uuid.New()
	}
	return nil
}

func (j *Job) BeforeCreate(tx *gorm.DB) error {
	if j.ID == uuid.Nil {
		j.ID = uuid.New()
	}
	return nil
}
