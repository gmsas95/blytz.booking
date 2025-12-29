package services

import (
	"blytz.cloud/backend/internal/models"
	"blytz.cloud/backend/internal/repository"
	"context"
	"errors"
	"gorm.io/gorm"
)

type ServiceService struct {
	serviceRepo  repository.ServiceRepository
	businessRepo repository.BusinessRepository
	db           *gorm.DB
}

func NewServiceService(
	serviceRepo repository.ServiceRepository,
	businessRepo repository.BusinessRepository,
	db *gorm.DB,
) *ServiceService {
	return &ServiceService{
		serviceRepo:  serviceRepo,
		businessRepo: businessRepo,
		db:           db,
	}
}

type CreateServiceRequest struct {
	BusinessID    string  `json:"business_id" validate:"required,uuid"`
	Name          string  `json:"name" validate:"required,min=3"`
	Description   string  `json:"description"`
	DurationMin   int     `json:"duration_min" validate:"required,min=5,max=480"`
	TotalPrice    float64 `json:"total_price" validate:"required,min=0"`
	DepositAmount float64 `json:"deposit_amount" validate:"omitempty,min=0"`
	IsActive      bool    `json:"is_active"`
	MaxCapacity   int     `json:"max_capacity" validate:"omitempty,min=1"`
}

type UpdateServiceRequest struct {
	Name          *string  `json:"name" validate:"omitempty,min=3"`
	Description   *string  `json:"description"`
	DurationMin   *int     `json:"duration_min" validate:"omitempty,min=5,max=480"`
	TotalPrice    *float64 `json:"total_price" validate:"omitempty,min=0"`
	DepositAmount *float64 `json:"deposit_amount" validate:"omitempty,min=0"`
	IsActive      *bool    `json:"is_active"`
	MaxCapacity   *int     `json:"max_capacity" validate:"omitempty,min=1"`
}

func (s *ServiceService) CreateService(ctx context.Context, req *CreateServiceRequest) (*models.Service, error) {
	business, err := s.businessRepo.GetByID(ctx, req.BusinessID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("business not found")
		}
		return nil, err
	}

	service := &models.Service{
		ID:            business.ID,
		BusinessID:    business.ID,
		Name:          req.Name,
		Description:   req.Description,
		DurationMin:   req.DurationMin,
		TotalPrice:    req.TotalPrice,
		DepositAmount: req.DepositAmount,
		IsActive:      req.IsActive,
		MaxCapacity:   req.MaxCapacity,
	}

	if err := s.serviceRepo.Create(ctx, service); err != nil {
		return nil, err
	}

	return service, nil
}

func (s *ServiceService) GetService(ctx context.Context, id string) (*models.Service, error) {
	return s.serviceRepo.GetByID(ctx, id)
}

func (s *ServiceService) ListServices(ctx context.Context, businessID string, offset, limit int) ([]*models.Service, int64, error) {
	return s.serviceRepo.ListByBusiness(ctx, businessID, offset, limit)
}

func (s *ServiceService) ListServicesBySlug(ctx context.Context, slug string, offset, limit int) ([]*models.Service, int64, error) {
	business, err := s.businessRepo.GetBySlug(ctx, slug)
	if err != nil {
		return nil, 0, err
	}

	businessID := business.ID.String()
	return s.serviceRepo.ListByBusiness(ctx, businessID, offset, limit)
}

func (s *ServiceService) UpdateService(ctx context.Context, id string, req *UpdateServiceRequest) (*models.Service, error) {
	service, err := s.serviceRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		service.Name = *req.Name
	}
	if req.Description != nil {
		service.Description = *req.Description
	}
	if req.DurationMin != nil {
		service.DurationMin = *req.DurationMin
	}
	if req.TotalPrice != nil {
		service.TotalPrice = *req.TotalPrice
	}
	if req.DepositAmount != nil {
		service.DepositAmount = *req.DepositAmount
	}
	if req.IsActive != nil {
		service.IsActive = *req.IsActive
	}
	if req.MaxCapacity != nil {
		service.MaxCapacity = *req.MaxCapacity
	}

	if err := s.serviceRepo.Update(ctx, service); err != nil {
		return nil, err
	}

	return service, nil
}

func (s *ServiceService) DeleteService(ctx context.Context, id string) error {
	return s.serviceRepo.Delete(ctx, id)
}
