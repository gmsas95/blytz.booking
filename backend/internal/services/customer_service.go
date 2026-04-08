package services

import (
	"blytz.cloud/backend/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CustomerService struct {
	*BaseService
}

func NewCustomerService(db *gorm.DB) *CustomerService {
	return &CustomerService{BaseService: NewBaseService(db)}
}

func (s *CustomerService) GetByBusiness(businessID uuid.UUID) ([]models.Customer, error) {
	var customers []models.Customer
	if err := s.DB.Where("business_id = ?", businessID).Order("created_at DESC").Find(&customers).Error; err != nil {
		return nil, err
	}
	return customers, nil
}

func (s *CustomerService) Create(customer *models.Customer) error {
	return s.DB.Create(customer).Error
}
