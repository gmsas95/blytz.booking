package services

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"blytz.cloud/backend/internal/models"

	"github.com/google/uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupBookingTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	dsn := fmt.Sprintf("file:%s?mode=memory&cache=shared", uuid.NewString())
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite db: %v", err)
	}

	statements := []string{
		`CREATE TABLE businesses (
			id text PRIMARY KEY,
			name text NOT NULL,
			slug text NOT NULL,
			vertical text NOT NULL,
			description text,
			theme_color text,
			created_at datetime,
			updated_at datetime
		)`,
		`CREATE TABLE services (
			id text PRIMARY KEY,
			business_id text NOT NULL,
			name text NOT NULL,
			description text,
			duration_min integer NOT NULL,
			total_price_minor integer NOT NULL,
			deposit_amount_minor integer NOT NULL,
			currency_code text NOT NULL,
			created_at datetime,
			updated_at datetime
		)`,
		`CREATE TABLE slots (
			id text PRIMARY KEY,
			business_id text NOT NULL,
			start_time datetime NOT NULL,
			end_time datetime NOT NULL,
			is_booked numeric NOT NULL DEFAULT 0,
			created_at datetime,
			updated_at datetime
		)`,
		`CREATE TABLE bookings (
			id text PRIMARY KEY,
			business_id text NOT NULL,
			service_id text NOT NULL,
			slot_id text NOT NULL,
			service_name text NOT NULL,
			slot_time datetime NOT NULL,
			name text NOT NULL,
			email text NOT NULL,
			phone text NOT NULL,
			status text NOT NULL,
			deposit_paid_minor integer NOT NULL,
			total_price_minor integer NOT NULL,
			currency_code text NOT NULL,
			created_at datetime,
			updated_at datetime
		)`,
	}

	for _, statement := range statements {
		if err := db.Exec(statement).Error; err != nil {
			t.Fatalf("create test schema: %v", err)
		}
	}

	return db
}

func seedBookingTestRecords(t *testing.T, db *gorm.DB) (models.Business, models.Service, models.Slot) {
	t.Helper()

	business := models.Business{
		ID:          uuid.New(),
		Name:        "DetailPro Automotive",
		Slug:        "detail-pro",
		Vertical:    "Automotive",
		Description: "Test workshop",
		ThemeColor:  "blue",
	}
	if err := db.Create(&business).Error; err != nil {
		t.Fatalf("create business: %v", err)
	}

	service := models.Service{
		ID:                 uuid.New(),
		BusinessID:         business.ID,
		Name:               "Full Interior Detail",
		Description:        "Deep clean",
		DurationMin:        120,
		TotalPriceMinor:    20000,
		DepositAmountMinor: 5000,
		CurrencyCode:       "USD",
	}
	if err := db.Create(&service).Error; err != nil {
		t.Fatalf("create service: %v", err)
	}

	now := time.Now().UTC().Truncate(time.Second)
	slot := models.Slot{
		ID:         uuid.New(),
		BusinessID: business.ID,
		StartTime:  now,
		EndTime:    now.Add(2 * time.Hour),
		IsBooked:   false,
	}
	if err := db.Create(&slot).Error; err != nil {
		t.Fatalf("create slot: %v", err)
	}

	return business, service, slot
}

func TestBookingServiceCreateMarksSlotBookedAndCopiesMoneyFields(t *testing.T) {
	db := setupBookingTestDB(t)
	business, service, slot := seedBookingTestRecords(t, db)
	bookingService := NewBookingService(db)

	booking := &models.Booking{
		BusinessID: business.ID,
		ServiceID:  service.ID,
		SlotID:     slot.ID,
		Customer: models.CustomerDetails{
			Name:  "Alice Smith",
			Email: "alice@example.com",
			Phone: "555-0101",
		},
	}

	if err := bookingService.Create(booking); err != nil {
		t.Fatalf("create booking: %v", err)
	}

	if booking.DepositPaidMinor != service.DepositAmountMinor {
		t.Fatalf("expected deposit_paid_minor %d, got %d", service.DepositAmountMinor, booking.DepositPaidMinor)
	}
	if booking.TotalPriceMinor != service.TotalPriceMinor {
		t.Fatalf("expected total_price_minor %d, got %d", service.TotalPriceMinor, booking.TotalPriceMinor)
	}
	if booking.CurrencyCode != service.CurrencyCode {
		t.Fatalf("expected currency_code %q, got %q", service.CurrencyCode, booking.CurrencyCode)
	}

	var persistedSlot models.Slot
	if err := db.First(&persistedSlot, "id = ?", slot.ID).Error; err != nil {
		t.Fatalf("reload slot: %v", err)
	}
	if !persistedSlot.IsBooked {
		t.Fatal("expected slot to be marked booked")
	}
}

func TestBookingServiceCreateReturnsConflictWhenSlotAlreadyBooked(t *testing.T) {
	db := setupBookingTestDB(t)
	business, service, slot := seedBookingTestRecords(t, db)
	bookingService := NewBookingService(db)

	firstBooking := &models.Booking{
		BusinessID: business.ID,
		ServiceID:  service.ID,
		SlotID:     slot.ID,
		Customer:   models.CustomerDetails{Name: "Alice", Email: "alice@example.com", Phone: "555-0101"},
	}
	if err := bookingService.Create(firstBooking); err != nil {
		t.Fatalf("create first booking: %v", err)
	}

	secondBooking := &models.Booking{
		BusinessID: business.ID,
		ServiceID:  service.ID,
		SlotID:     slot.ID,
		Customer:   models.CustomerDetails{Name: "Bob", Email: "bob@example.com", Phone: "555-0102"},
	}

	err := bookingService.Create(secondBooking)
	if !errors.Is(err, ErrConflict) {
		t.Fatalf("expected ErrConflict, got %v", err)
	}

	var bookingCount int64
	if err := db.Model(&models.Booking{}).Count(&bookingCount).Error; err != nil {
		t.Fatalf("count bookings: %v", err)
	}
	if bookingCount != 1 {
		t.Fatalf("expected 1 persisted booking, got %d", bookingCount)
	}
}
