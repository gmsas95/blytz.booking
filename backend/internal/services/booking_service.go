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
	return s.DB.Transaction(func(tx *gorm.DB) error {
		var service models.Service
		if err := tx.Where("id = ? AND business_id = ?", booking.ServiceID, booking.BusinessID).First(&service).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return ErrBadRequest
			}
			return err
		}

		var slot models.Slot
		if err := tx.Where("id = ? AND business_id = ?", booking.SlotID, booking.BusinessID).First(&slot).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return ErrBadRequest
			}
			return err
		}

		reserveResult := tx.Model(&models.Slot{}).
			Where("id = ? AND business_id = ? AND is_booked = ?", booking.SlotID, booking.BusinessID, false).
			Update("is_booked", true)
		if reserveResult.Error != nil {
			return reserveResult.Error
		}
		if reserveResult.RowsAffected == 0 {
			return ErrConflict
		}

		booking.ServiceName = service.Name
		booking.SlotTime = slot.StartTime
		booking.DepositPaidMinor = service.DepositAmountMinor
		booking.TotalPriceMinor = service.TotalPriceMinor
		booking.CurrencyCode = service.CurrencyCode

		if err := tx.Create(booking).Error; err != nil {
			return err
		}

		return nil
	})
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

	// Mark slot as available again
	if err := tx.Model(&models.Slot{}).Where("id = ?", booking.SlotID).Update("is_booked", false).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
