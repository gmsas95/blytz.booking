package services

import (
	"blytz.cloud/backend/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ServiceService struct {
	*BaseService
}

func NewServiceService(db *gorm.DB) *ServiceService {
	return &ServiceService{
		BaseService: NewBaseService(db),
	}
}

func (s *ServiceService) GetByBusiness(businessID uuid.UUID) ([]models.Service, error) {
	var services []models.Service
	if err := s.DB.Where("business_id = ?", businessID).Find(&services).Error; err != nil {
		return nil, err
	}
	return services, nil
}

func (s *ServiceService) GetByID(id uuid.UUID) (*models.Service, error) {
	var service models.Service
	if err := s.DB.Where("id = ?", id).First(&service).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &service, nil
}

func (s *ServiceService) Create(service *models.Service) error {
	if err := s.DB.Create(service).Error; err != nil {
		return err
	}
	return nil
}

func (s *ServiceService) Update(id uuid.UUID, service *models.Service) error {
	result := s.DB.Model(&models.Service{}).Where("id = ?", id).Updates(service)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

func (s *ServiceService) Delete(id uuid.UUID) error {
	result := s.DB.Delete(&models.Service{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}
