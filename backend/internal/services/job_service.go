package services

import (
	"blytz.cloud/backend/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type JobService struct {
	*BaseService
}

func NewJobService(db *gorm.DB) *JobService {
	return &JobService{BaseService: NewBaseService(db)}
}

func (s *JobService) GetByBusiness(businessID uuid.UUID) ([]models.Job, error) {
	var jobs []models.Job
	if err := s.DB.Where("business_id = ?", businessID).
		Preload("Customer", "business_id = ?", businessID).
		Preload("Vehicle", "business_id = ?", businessID).
		Order("scheduled_at ASC").
		Find(&jobs).Error; err != nil {
		return nil, err
	}
	return jobs, nil
}

func (s *JobService) Create(job *models.Job) error {
	return s.DB.Transaction(func(tx *gorm.DB) error {
		var customer models.Customer
		if err := tx.Where("id = ? AND business_id = ?", job.CustomerID, job.BusinessID).First(&customer).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return ErrBadRequest
			}
			return err
		}

		var vehicle models.Vehicle
		if err := tx.Where("id = ? AND business_id = ?", job.VehicleID, job.BusinessID).First(&vehicle).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return ErrBadRequest
			}
			return err
		}
		if vehicle.CustomerID != job.CustomerID {
			return ErrBadRequest
		}

		if job.BookingID != nil {
			var booking models.Booking
			if err := tx.Where("id = ? AND business_id = ?", *job.BookingID, job.BusinessID).First(&booking).Error; err != nil {
				if err == gorm.ErrRecordNotFound {
					return ErrBadRequest
				}
				return err
			}
		}

		return tx.Create(job).Error
	})
}
