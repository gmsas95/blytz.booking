package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BookingStatus string

const (
	BookingStatusPending   BookingStatus = "PENDING"
	BookingStatusConfirmed BookingStatus = "CONFIRMED"
	BookingStatusCompleted BookingStatus = "COMPLETED"
	BookingStatusCancelled BookingStatus = "CANCELLED"
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
	ID             uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	BusinessID     uuid.UUID `json:"business_id" gorm:"type:uuid;not null"`
	Name           string    `json:"name" gorm:"not null"`
	Description    string    `json:"description"`
	DurationMin    int       `json:"duration_min" gorm:"not null"`
	TotalPrice     float64   `json:"total_price" gorm:"not null"`
	DepositAmount  float64   `json:"deposit_amount" gorm:"not null"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	Business      Business  `json:"business" gorm:"foreignKey:BusinessID"`
}

type Slot struct {
	ID        uuid.UUID    `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	BusinessID uuid.UUID   `json:"business_id" gorm:"type:uuid;not null;index"`
	StartTime time.Time    `json:"start_time" gorm:"not null;index"`
	EndTime   time.Time    `json:"end_time" gorm:"not null"`
	IsBooked  bool         `json:"is_booked" gorm:"default:false"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
	Business  Business    `json:"business" gorm:"foreignKey:BusinessID"`
}

type CustomerDetails struct {
	Name  string `json:"name" gorm:"not null"`
	Email string `json:"email" gorm:"not null"`
	Phone string `json:"phone" gorm:"not null"`
}

type Booking struct {
	ID            uuid.UUID        `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	BusinessID    uuid.UUID        `json:"business_id" gorm:"type:uuid;not null;index"`
	ServiceID     uuid.UUID        `json:"service_id" gorm:"type:uuid;not null"`
	SlotID        uuid.UUID        `json:"slot_id" gorm:"type:uuid;not null;index"`
	ServiceName   string           `json:"service_name" gorm:"not null"`
	SlotTime      time.Time        `json:"slot_time" gorm:"not null"`
	Customer      CustomerDetails   `json:"customer" gorm:"embedded"`
	Status        BookingStatus     `json:"status" gorm:"not null;default:'PENDING'"`
	DepositPaid   float64          `json:"deposit_paid" gorm:"not null"`
	TotalPrice    float64          `json:"total_price" gorm:"not null"`
	CreatedAt     time.Time        `json:"created_at"`
	UpdatedAt     time.Time        `json:"updated_at"`
	Business      Business         `json:"business" gorm:"foreignKey:BusinessID"`
	Service       Service          `json:"service" gorm:"foreignKey:ServiceID"`
	Slot          Slot             `json:"slot" gorm:"foreignKey:SlotID"`
}

// User model for operators
type User struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Email     string    `json:"email" gorm:"uniqueIndex;not null"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
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
