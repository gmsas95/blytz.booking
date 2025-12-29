package repository

import (
	"blytz.cloud/backend/internal/models"
	"context"
	"gorm.io/gorm"
)

type BusinessRepository interface {
	Create(ctx context.Context, business *models.Business) error
	GetByID(ctx context.Context, id string) (*models.Business, error)
	GetBySlug(ctx context.Context, slug string) (*models.Business, error)
	List(ctx context.Context, offset, limit int) ([]*models.Business, int64, error)
	Update(ctx context.Context, business *models.Business) error
	Delete(ctx context.Context, id string) error
	WithTx(tx *gorm.DB) BusinessRepository
}

type businessRepository struct {
	db *gorm.DB
}

func NewBusinessRepository(db *gorm.DB) BusinessRepository {
	return &businessRepository{db: db}
}

func (r *businessRepository) Create(ctx context.Context, business *models.Business) error {
	return r.db.WithContext(ctx).Create(business).Error
}

func (r *businessRepository) GetByID(ctx context.Context, id string) (*models.Business, error) {
	var business models.Business
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&business).Error
	if err != nil {
		return nil, err
	}
	return &business, nil
}

func (r *businessRepository) GetBySlug(ctx context.Context, slug string) (*models.Business, error) {
	var business models.Business
	err := r.db.WithContext(ctx).Where("slug = ?", slug).First(&business).Error
	if err != nil {
		return nil, err
	}
	return &business, nil
}

func (r *businessRepository) List(ctx context.Context, offset, limit int) ([]*models.Business, int64, error) {
	var businesses []*models.Business
	var total int64

	if err := r.db.WithContext(ctx).Model(&models.Business{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := r.db.WithContext(ctx).Offset(offset).Limit(limit).Find(&businesses).Error
	return businesses, total, err
}

func (r *businessRepository) Update(ctx context.Context, business *models.Business) error {
	return r.db.WithContext(ctx).Save(business).Error
}

func (r *businessRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&models.Business{}, "id = ?", id).Error
}

func (r *businessRepository) WithTx(tx *gorm.DB) BusinessRepository {
	return &businessRepository{db: tx}
}
