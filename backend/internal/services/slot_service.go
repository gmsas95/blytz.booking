package services

import (
	"blytz.cloud/backend/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SlotService struct {
	*BaseService
}

func NewSlotService(db *gorm.DB) *SlotService {
	return &SlotService{
		BaseService: NewBaseService(db),
	}
}

func (s *SlotService) GetAvailableByBusiness(businessID uuid.UUID) ([]models.Slot, error) {
	var slots []models.Slot
	if err := s.DB.Where("business_id = ? AND is_booked = ?", businessID, false).Order("start_time").Find(&slots).Error; err != nil {
		return nil, err
	}
	return slots, nil
}

func (s *SlotService) GetByBusiness(businessID uuid.UUID) ([]models.Slot, error) {
	var slots []models.Slot
	if err := s.DB.Where("business_id = ?", businessID).Order("start_time").Find(&slots).Error; err != nil {
		return nil, err
	}
	return slots, nil
}

func (s *SlotService) GetByID(id uuid.UUID) (*models.Slot, error) {
	var slot models.Slot
	if err := s.DB.Where("id = ?", id).First(&slot).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &slot, nil
}

func (s *SlotService) Create(slot *models.Slot) error {
	if err := s.DB.Create(slot).Error; err != nil {
		return err
	}
	return nil
}

func (s *SlotService) Delete(id uuid.UUID) error {
	if err := s.DB.Delete(&models.Slot{}, id).Error; err != nil {
		return err
	}
	return nil
}
