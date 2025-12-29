package repository

import (
	"blytz.cloud/backend/internal/models"
	"context"
	"gorm.io/gorm"
	"time"
)

type SlotRepository interface {
	Create(ctx context.Context, slot *models.Slot) error
	GetByID(ctx context.Context, id string) (*models.Slot, error)
	ListAvailable(ctx context.Context, businessID, serviceID string, startDate, endDate time.Time) ([]*models.Slot, error)
	Update(ctx context.Context, slot *models.Slot) error
	MarkBooked(ctx context.Context, id string) error
	Delete(ctx context.Context, id string) error
	WithTx(tx *gorm.DB) SlotRepository
}

type slotRepository struct {
	db *gorm.DB
}

func NewSlotRepository(db *gorm.DB) SlotRepository {
	return &slotRepository{db: db}
}

func (r *slotRepository) Create(ctx context.Context, slot *models.Slot) error {
	return r.db.WithContext(ctx).Create(slot).Error
}

func (r *slotRepository) GetByID(ctx context.Context, id string) (*models.Slot, error) {
	var slot models.Slot
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&slot).Error
	if err != nil {
		return nil, err
	}
	return &slot, nil
}

func (r *slotRepository) ListAvailable(ctx context.Context, businessID, serviceID string, startDate, endDate time.Time) ([]*models.Slot, error) {
	var slots []*models.Slot

	query := r.db.WithContext(ctx).
		Where("business_id = ?", businessID).
		Where("is_booked = ?", false).
		Where("start_time >= ?", startDate).
		Where("end_time <= ?", endDate).
		Preload("Service")

	if serviceID != "" {
		query = query.Where("service_id = ? OR service_id IS NULL", serviceID)
	}

	err := query.Order("start_time ASC").Find(&slots).Error
	return slots, err
}

func (r *slotRepository) Update(ctx context.Context, slot *models.Slot) error {
	return r.db.WithContext(ctx).Save(slot).Error
}

func (r *slotRepository) MarkBooked(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Model(&models.Slot{}).Where("id = ?", id).Update("is_booked", true).Error
}

func (r *slotRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&models.Slot{}, "id = ?", id).Error
}

func (r *slotRepository) WithTx(tx *gorm.DB) SlotRepository {
	return &slotRepository{db: tx}
}
