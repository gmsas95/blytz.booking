package repository

import (
	"blytz.cloud/backend/internal/models"
	"context"
	"gorm.io/gorm"
)

type BookingHistoryRepository interface {
	Create(ctx context.Context, history *models.BookingHistory) error
	ListByBooking(ctx context.Context, bookingID string) ([]*models.BookingHistory, error)
	WithTx(tx *gorm.DB) BookingHistoryRepository
}

type bookingHistoryRepository struct {
	db *gorm.DB
}

func NewBookingHistoryRepository(db *gorm.DB) BookingHistoryRepository {
	return &bookingHistoryRepository{db: db}
}

func (r *bookingHistoryRepository) Create(ctx context.Context, history *models.BookingHistory) error {
	return r.db.WithContext(ctx).Create(history).Error
}

func (r *bookingHistoryRepository) ListByBooking(ctx context.Context, bookingID string) ([]*models.BookingHistory, error) {
	var histories []*models.BookingHistory
	err := r.db.WithContext(ctx).Where("booking_id = ?", bookingID).Order("timestamp DESC").Find(&histories).Error
	return histories, err
}

func (r *bookingHistoryRepository) WithTx(tx *gorm.DB) BookingHistoryRepository {
	return &bookingHistoryRepository{db: tx}
}
