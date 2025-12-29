package repository

import (
	"blytz.cloud/backend/internal/models"
	"context"
	"gorm.io/gorm"
)

type BookingRepository interface {
	Create(ctx context.Context, booking *models.Booking) error
	GetByID(ctx context.Context, id string) (*models.Booking, error)
	ListByBusiness(ctx context.Context, businessID string, offset, limit int) ([]*models.Booking, int64, error)
	Update(ctx context.Context, booking *models.Booking) error
	UpdateStatus(ctx context.Context, id string, status models.BookingStatus) error
	Delete(ctx context.Context, id string) error
	WithTx(tx *gorm.DB) BookingRepository
}

type bookingRepository struct {
	db *gorm.DB
}

func NewBookingRepository(db *gorm.DB) BookingRepository {
	return &bookingRepository{db: db}
}

func (r *bookingRepository) Create(ctx context.Context, booking *models.Booking) error {
	return r.db.WithContext(ctx).Create(booking).Error
}

func (r *bookingRepository) GetByID(ctx context.Context, id string) (*models.Booking, error) {
	var booking models.Booking
	err := r.db.WithContext(ctx).
		Preload("Business").
		Preload("Service").
		Preload("Slot").
		Preload("Customer").
		Where("id = ?", id).
		First(&booking).Error
	if err != nil {
		return nil, err
	}
	return &booking, nil
}

func (r *bookingRepository) ListByBusiness(ctx context.Context, businessID string, offset, limit int) ([]*models.Booking, int64, error) {
	var bookings []*models.Booking
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Booking{}).Where("business_id = ?", businessID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Offset(offset).Limit(limit).
		Preload("Business").
		Preload("Service").
		Preload("Slot").
		Preload("Customer").
		Order("created_at DESC").
		Find(&bookings).Error

	return bookings, total, err
}

func (r *bookingRepository) Update(ctx context.Context, booking *models.Booking) error {
	return r.db.WithContext(ctx).Save(booking).Error
}

func (r *bookingRepository) UpdateStatus(ctx context.Context, id string, status models.BookingStatus) error {
	return r.db.WithContext(ctx).Model(&models.Booking{}).Where("id = ?", id).Update("status", status).Error
}

func (r *bookingRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&models.Booking{}, "id = ?", id).Error
}

func (r *bookingRepository) WithTx(tx *gorm.DB) BookingRepository {
	return &bookingRepository{db: tx}
}
