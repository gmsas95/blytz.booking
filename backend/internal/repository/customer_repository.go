package repository

import (
	"blytz.cloud/backend/internal/models"
	"context"
	"gorm.io/gorm"
)

type CustomerRepository interface {
	Create(ctx context.Context, customer *models.Customer) error
	GetByID(ctx context.Context, id string) (*models.Customer, error)
	GetOrCreate(ctx context.Context, customer *models.Customer) (*models.Customer, error)
	ListByBusiness(ctx context.Context, businessID string, offset, limit int) ([]*models.Customer, int64, error)
	Update(ctx context.Context, customer *models.Customer) error
	Delete(ctx context.Context, id string) error
	WithTx(tx *gorm.DB) CustomerRepository
}

type customerRepository struct {
	db *gorm.DB
}

func NewCustomerRepository(db *gorm.DB) CustomerRepository {
	return &customerRepository{db: db}
}

func (r *customerRepository) Create(ctx context.Context, customer *models.Customer) error {
	return r.db.WithContext(ctx).Create(customer).Error
}

func (r *customerRepository) GetByID(ctx context.Context, id string) (*models.Customer, error) {
	var customer models.Customer
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&customer).Error
	if err != nil {
		return nil, err
	}
	return &customer, nil
}

func (r *customerRepository) GetOrCreate(ctx context.Context, customer *models.Customer) (*models.Customer, error) {
	var existing models.Customer
	err := r.db.WithContext(ctx).
		Where("business_id = ? AND email = ?", customer.BusinessID, customer.Email).
		First(&existing).Error

	if err == nil {
		return &existing, nil
	}

	if err != gorm.ErrRecordNotFound {
		return nil, err
	}

	if err := r.db.WithContext(ctx).Create(customer).Error; err != nil {
		return nil, err
	}

	return customer, nil
}

func (r *customerRepository) ListByBusiness(ctx context.Context, businessID string, offset, limit int) ([]*models.Customer, int64, error) {
	var customers []*models.Customer
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Customer{}).Where("business_id = ?", businessID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&customers).Error
	return customers, total, err
}

func (r *customerRepository) Update(ctx context.Context, customer *models.Customer) error {
	return r.db.WithContext(ctx).Save(customer).Error
}

func (r *customerRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&models.Customer{}, "id = ?", id).Error
}

func (r *customerRepository) WithTx(tx *gorm.DB) CustomerRepository {
	return &customerRepository{db: tx}
}
