package repository

import (
	"blytz.cloud/backend/internal/models"
	"context"
	"gorm.io/gorm"
)

type ServiceRepository interface {
	Create(ctx context.Context, service *models.Service) error
	GetByID(ctx context.Context, id string) (*models.Service, error)
	ListByBusiness(ctx context.Context, businessID string, offset, limit int) ([]*models.Service, int64, error)
	Update(ctx context.Context, service *models.Service) error
	Delete(ctx context.Context, id string) error
	WithTx(tx *gorm.DB) ServiceRepository
}

type serviceRepository struct {
	db *gorm.DB
}

func NewServiceRepository(db *gorm.DB) ServiceRepository {
	return &serviceRepository{db: db}
}

func (r *serviceRepository) Create(ctx context.Context, service *models.Service) error {
	return r.db.WithContext(ctx).Create(service).Error
}

func (r *serviceRepository) GetByID(ctx context.Context, id string) (*models.Service, error) {
	var service models.Service
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&service).Error
	if err != nil {
		return nil, err
	}
	return &service, nil
}

func (r *serviceRepository) ListByBusiness(ctx context.Context, businessID string, offset, limit int) ([]*models.Service, int64, error) {
	var services []*models.Service
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Service{}).Where("business_id = ?", businessID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Offset(offset).Limit(limit).Find(&services).Error
	return services, total, err
}

func (r *serviceRepository) Update(ctx context.Context, service *models.Service) error {
	return r.db.WithContext(ctx).Save(service).Error
}

func (r *serviceRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&models.Service{}, "id = ?", id).Error
}

func (r *serviceRepository) WithTx(tx *gorm.DB) ServiceRepository {
	return &serviceRepository{db: tx}
}
