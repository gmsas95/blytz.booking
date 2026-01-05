package services

import (
	"blytz.cloud/backend/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BusinessService struct {
	*BaseService
}

func NewBusinessService(db *gorm.DB) *BusinessService {
	return &BusinessService{
		BaseService: NewBaseService(db),
	}
}

func (s *BusinessService) GetAll() ([]models.Business, error) {
	var businesses []models.Business
	if err := s.DB.Find(&businesses).Error; err != nil {
		return nil, err
	}
	return businesses, nil
}

func (s *BusinessService) GetByID(id uuid.UUID) (*models.Business, error) {
	var business models.Business
	if err := s.DB.Where("id = ?", id).First(&business).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &business, nil
}

func (s *BusinessService) GetBySlug(slug string) (*models.Business, error) {
	var business models.Business
	if err := s.DB.Where("slug = ?", slug).First(&business).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &business, nil
}

func (s *BusinessService) GetByUser(userID uuid.UUID) ([]models.Business, error) {
	var businesses []models.Business
	if err := s.DB.Where("owner_id = ?", userID).Find(&businesses).Error; err != nil {
		return nil, err
	}
	return businesses, nil
}

func (s *BusinessService) Create(business *models.Business) error {
	if err := s.DB.Create(business).Error; err != nil {
		return err
	}
	return nil
}

func (s *BusinessService) Update(id uuid.UUID, business *models.Business) error {
	result := s.DB.Model(&models.Business{}).Where("id = ?", id).Updates(business)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

func (s *BusinessService) Delete(id uuid.UUID) error {
	result := s.DB.Delete(&models.Business{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}
