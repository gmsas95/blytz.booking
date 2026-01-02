package services

import (
	"time"

	"blytz.cloud/backend/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AvailabilityService struct {
	*BaseService
}

func NewAvailabilityService(db *gorm.DB) *AvailabilityService {
	return &AvailabilityService{
		BaseService: NewBaseService(db),
	}
}

func (s *AvailabilityService) GetByBusiness(businessID uuid.UUID) ([]models.BusinessAvailability, error) {
	var availabilities []models.BusinessAvailability
	if err := s.DB.Where("business_id = ?", businessID).Order("day_of_week").Find(&availabilities).Error; err != nil {
		return nil, err
	}
	return availabilities, nil
}

func (s *AvailabilityService) Upsert(businessID uuid.UUID, availability *models.BusinessAvailability) error {
	// Check if availability exists for this day
	var existing models.BusinessAvailability
	err := s.DB.Where("business_id = ? AND day_of_week = ?", businessID, availability.DayOfWeek).First(&existing).Error

	if err == gorm.ErrRecordNotFound {
		// Create new
		availability.ID = uuid.New()
		availability.BusinessID = businessID
		return s.DB.Create(availability).Error
	}

	if err != nil {
		return err
	}

	// Update existing
	availability.ID = existing.ID
	availability.BusinessID = businessID
	return s.DB.Save(availability).Error
}

func (s *AvailabilityService) GenerateSlotsFromAvailability(businessID uuid.UUID, startDate, endDate string, durationMin int) ([]models.Slot, error) {
	// Get business
	var business models.Business
	if err := s.DB.Where("id = ?", businessID).First(&business).Error; err != nil {
		return nil, err
	}

	// Get availability schedule
	availabilities, err := s.GetByBusiness(businessID)
	if err != nil {
		return nil, err
	}

	// Parse dates
	start, _ := time.Parse("2006-01-02", startDate)
	end, _ := time.Parse("2006-01-02", endDate)

	var slots []models.Slot
	slotDuration := time.Duration(durationMin) * time.Minute

	// Generate slots for each day in range
	for d := start; !d.After(end); d = d.Add(24 * time.Hour) {
		dayOfWeek := int(d.Weekday())
		if dayOfWeek == 0 {
			dayOfWeek = 6 // Sunday should be last
		} else {
			dayOfWeek-- // Monday = 0
		}

		// Find availability for this day
		var avail *models.BusinessAvailability
		for _, a := range availabilities {
			if a.DayOfWeek == dayOfWeek {
				avail = &a
				break
			}
		}

		// Skip if closed or no availability
		if avail == nil || avail.IsClosed || avail.StartTime == "" || avail.EndTime == "" {
			continue
		}

		// Parse availability times
		availStart, _ := time.Parse("15:04", avail.StartTime)
		availEnd, _ := time.Parse("15:04", avail.EndTime)

		// Set date component to availability times
		dayStart := time.Date(d.Year(), d.Month(), d.Day(), availStart.Hour(), availStart.Minute(), 0, 0, d.Location())
		dayEnd := time.Date(d.Year(), d.Month(), d.Day(), availEnd.Hour(), availEnd.Minute(), 0, 0, d.Location())

		// Generate time slots
		for slotStart := dayStart; slotStart.Add(slotDuration).Before(dayEnd); slotStart = slotStart.Add(slotDuration) {
			slot := models.Slot{
				ID:         uuid.New(),
				BusinessID: businessID,
				StartTime:  slotStart,
				EndTime:    slotStart.Add(slotDuration),
				IsBooked:   false,
			}
			slots = append(slots, slot)
		}
	}

	return slots, nil
}
