package services

import (
	"blytz.cloud/backend/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BookingService struct {
	*BaseService
}

func NewBookingService(db *gorm.DB) *BookingService {
	return &BookingService{
		BaseService: NewBaseService(db),
	}
}

func (s *BookingService) Create(booking *models.Booking) error {
	tx := s.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var slot models.Slot
	if err := tx.Set("gorm:query_option", "FOR UPDATE").
		Where("id = ?", booking.SlotID).First(&slot).Error; err != nil {
		tx.Rollback()
		if err == gorm.ErrRecordNotFound {
			return ErrBadRequest
		}
		return err
	}

	var business models.Business
	if err := tx.Where("id = ?", booking.BusinessID).First(&business).Error; err != nil {
		tx.Rollback()
		return err
	}

	if slot.BookingCount >= business.MaxBookings {
		tx.Rollback()
		return ErrSlotFull
	}

	var service models.Service
	if err := tx.Where("id = ?", booking.ServiceID).First(&service).Error; err != nil {
		tx.Rollback()
		if err == gorm.ErrRecordNotFound {
			return ErrBadRequest
		}
		return err
	}

	booking.ServiceName = service.Name
	booking.SlotTime = slot.StartTime
	booking.DepositPaid = service.DepositAmount
	booking.TotalPrice = service.TotalPrice

	if err := tx.Create(booking).Error; err != nil {
		tx.Rollback()
		return err
	}

	result := tx.Model(&models.Slot{}).
		Where("id = ? AND booking_count < ?", booking.SlotID, business.MaxBookings).
		Update("booking_count", gorm.Expr("booking_count + 1"))

	if result.RowsAffected == 0 {
		tx.Rollback()
		return ErrSlotFull
	}

	if slot.BookingCount+1 >= business.MaxBookings {
		tx.Model(&models.Slot{}).
			Where("id = ?", booking.SlotID).
			Update("is_booked", true)
	}

	return tx.Commit().Error
}

func (s *BookingService) GetByBusiness(businessID uuid.UUID) ([]models.Booking, error) {
	var bookings []models.Booking
	if err := s.DB.Where("business_id = ?", businessID).Order("created_at DESC").Find(&bookings).Error; err != nil {
		return nil, err
	}
	return bookings, nil
}

func (s *BookingService) GetByID(id uuid.UUID) (*models.Booking, error) {
	var booking models.Booking
	if err := s.DB.Where("id = ?", id).Preload("Service").Preload("Slot").First(&booking).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &booking, nil
}

func (s *BookingService) UpdateStatus(id uuid.UUID, status models.BookingStatus) error {
	result := s.DB.Model(&models.Booking{}).Where("id = ?", id).Update("status", status)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

func (s *BookingService) Cancel(id uuid.UUID) error {
	// Get booking
	var booking models.Booking
	if err := s.DB.Where("id = ?", id).First(&booking).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return ErrNotFound
		}
		return err
	}

	// Get slot
	var slot models.Slot
	if err := s.DB.Where("id = ?", booking.SlotID).First(&slot).Error; err != nil {
		return err
	}

	// Start transaction
	tx := s.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Update booking status
	if err := tx.Model(&booking).Update("status", models.BookingStatusCancelled).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Decrement slot booking count
	if slot.BookingCount > 0 {
		if err := tx.Model(&models.Slot{}).Where("id = ?", booking.SlotID).Update("booking_count", slot.BookingCount-1).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	// Mark slot as available
	if err := tx.Model(&models.Slot{}).Where("id = ?", booking.SlotID).Update("is_booked", false).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
