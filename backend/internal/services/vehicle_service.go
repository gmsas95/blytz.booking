package services

import (
	"blytz.cloud/backend/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type VehicleService struct {
	*BaseService
}

func NewVehicleService(db *gorm.DB) *VehicleService {
	return &VehicleService{BaseService: NewBaseService(db)}
}

func (s *VehicleService) GetByBusiness(businessID uuid.UUID) ([]models.Vehicle, error) {
	var vehicles []models.Vehicle
	if err := s.DB.Where("business_id = ?", businessID).Preload("Customer", "business_id = ?", businessID).Order("created_at DESC").Find(&vehicles).Error; err != nil {
		return nil, err
	}
	return vehicles, nil
}

func (s *VehicleService) Create(vehicle *models.Vehicle) error {
	var customer models.Customer
	if err := s.DB.Where("id = ? AND business_id = ?", vehicle.CustomerID, vehicle.BusinessID).First(&customer).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return ErrBadRequest
		}
		return err
	}
	return s.DB.Create(vehicle).Error
}
