package services

import (
	"blytz.cloud/backend/internal/models"
	"blytz.cloud/backend/internal/repository"
	"context"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CustomerService struct {
	customerRepo repository.CustomerRepository
	db           *gorm.DB
}

func NewCustomerService(
	customerRepo repository.CustomerRepository,
	db *gorm.DB,
) *CustomerService {
	return &CustomerService{
		customerRepo: customerRepo,
		db:           db,
	}
}

type CreateCustomerRequest struct {
	BusinessID string `json:"business_id" validate:"required,uuid"`
	Name       string `json:"name" validate:"required,min=2"`
	Email      string `json:"email" validate:"required,email"`
	Phone      string `json:"phone" validate:"required,phone"`
	Notes      string `json:"notes"`
}

type UpdateCustomerRequest struct {
	Name  *string `json:"name" validate:"omitempty,min=2"`
	Email *string `json:"email" validate:"omitempty,email"`
	Phone *string `json:"phone" validate:"omitempty,phone"`
	Notes *string `json:"notes"`
}

func (s *CustomerService) CreateCustomer(ctx context.Context, req *CreateCustomerRequest) (*models.Customer, error) {
	businessID, err := uuid.Parse(req.BusinessID)
	if err != nil {
		return nil, err
	}

	customer := &models.Customer{
		ID:    businessID,
		Name:  req.Name,
		Email: req.Email,
		Phone: req.Phone,
		Notes: req.Notes,
	}

	if err := s.customerRepo.Create(ctx, customer); err != nil {
		return nil, err
	}

	return customer, nil
}

func (s *CustomerService) GetCustomer(ctx context.Context, id string) (*models.Customer, error) {
	return s.customerRepo.GetByID(ctx, id)
}

func (s *CustomerService) ListCustomers(ctx context.Context, businessID string, offset, limit int) ([]*models.Customer, int64, error) {
	return s.customerRepo.ListByBusiness(ctx, businessID, offset, limit)
}

func (s *CustomerService) UpdateCustomer(ctx context.Context, id string, req *UpdateCustomerRequest) (*models.Customer, error) {
	customer, err := s.customerRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		customer.Name = *req.Name
	}
	if req.Email != nil {
		customer.Email = *req.Email
	}
	if req.Phone != nil {
		customer.Phone = *req.Phone
	}
	if req.Notes != nil {
		customer.Notes = *req.Notes
	}

	if err := s.customerRepo.Update(ctx, customer); err != nil {
		return nil, err
	}

	return customer, nil
}

func (s *CustomerService) DeleteCustomer(ctx context.Context, id string) error {
	return s.customerRepo.Delete(ctx, id)
}
