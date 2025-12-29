package services

import (
	"blytz.cloud/backend/internal/models"
	"blytz.cloud/backend/internal/repository"
	"context"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type SlotService struct {
	slotRepo     repository.SlotRepository
	businessRepo repository.BusinessRepository
	db           *gorm.DB
}

func NewSlotService(
	slotRepo repository.SlotRepository,
	businessRepo repository.BusinessRepository,
	db *gorm.DB,
) *SlotService {
	return &SlotService{
		slotRepo:     slotRepo,
		businessRepo: businessRepo,
		db:           db,
	}
}

type CreateSlotsRequest struct {
	BusinessID string `json:"business_id" validate:"required,uuid"`
	Slots      []Slot `json:"slots" validate:"required,dive"`
}

type Slot struct {
	StartTime time.Time `json:"start_time" validate:"required"`
	EndTime   time.Time `json:"end_time" validate:"required,gtfield=StartTime"`
	ServiceID *string   `json:"service_id" validate:"omitempty,uuid"`
}

type RecurringScheduleRequest struct {
	BusinessID   string      `json:"business_id" validate:"required,uuid"`
	Name         string      `json:"name" validate:"required"`
	DaysOfWeek   []int       `json:"days_of_week" validate:"required,min=1,max=7,dive,min=0,max=6"`
	StartTime    string      `json:"start_time" validate:"required,^([0-1]?[0-9]|2[0-3]):[0-5][0-9]$"`
	EndTime      string      `json:"end_time" validate:"required,^([0-1]?[0-9]|2[0-3]):[0-5][0-9]$"`
	StartDate    time.Time   `json:"start_date" validate:"required"`
	EndDate      time.Time   `json:"end_date" validate:"required,gtfield=StartDate"`
	ExcludeDates []time.Time `json:"exclude_dates"`
}

func (s *SlotService) CreateSlots(ctx context.Context, req *CreateSlotsRequest) ([]*models.Slot, error) {
	businessID, err := uuid.Parse(req.BusinessID)
	if err != nil {
		return nil, err
	}

	var slots []*models.Slot
	for _, slotData := range req.Slots {
		var serviceID *uuid.UUID
		if slotData.ServiceID != nil {
			sid, err := uuid.Parse(*slotData.ServiceID)
			if err != nil {
				return nil, err
			}
			serviceID = &sid
		}

		slot := &models.Slot{
			ID:         uuid.New(),
			BusinessID: businessID,
			ServiceID:  serviceID,
			StartTime:  slotData.StartTime,
			EndTime:    slotData.EndTime,
			IsBooked:   false,
			Capacity:   1,
		}
		slots = append(slots, slot)
	}

	if err := s.db.WithContext(ctx).Create(&slots).Error; err != nil {
		return nil, err
	}

	return slots, nil
}

func (s *SlotService) ListAvailableSlots(ctx context.Context, businessID, serviceID string, startDate, endDate time.Time) ([]*models.Slot, error) {
	return s.slotRepo.ListAvailable(ctx, businessID, serviceID, startDate, endDate)
}

func (s *SlotService) ListAvailableSlotsBySlug(ctx context.Context, slug, serviceID string, startDate, endDate *time.Time) ([]*models.Slot, error) {
	business, err := s.businessRepo.GetBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}

	businessID := business.ID.String()

	var start, end time.Time
	if startDate != nil {
		start = *startDate
	}
	if endDate != nil {
		end = *endDate
	}

	return s.slotRepo.ListAvailable(ctx, businessID, serviceID, start, end)
}

func (s *SlotService) GetSlot(ctx context.Context, id string) (*models.Slot, error) {
	return s.slotRepo.GetByID(ctx, id)
}

func (s *SlotService) DeleteSlot(ctx context.Context, id string) error {
	return s.slotRepo.Delete(ctx, id)
}

func (s *SlotService) CreateRecurringSchedule(ctx context.Context, req *RecurringScheduleRequest) (*models.RecurringSchedule, error) {
	businessID, err := uuid.Parse(req.BusinessID)
	if err != nil {
		return nil, err
	}

	schedule := &models.RecurringSchedule{
		ID:           uuid.New(),
		BusinessID:   businessID,
		Name:         req.Name,
		DaysOfWeek:   req.DaysOfWeek,
		StartTime:    req.StartTime,
		EndTime:      req.EndTime,
		StartDate:    req.StartDate,
		EndDate:      req.EndDate,
		ExcludeDates: req.ExcludeDates,
		IsActive:     true,
	}

	if err := s.db.WithContext(ctx).Create(schedule).Error; err != nil {
		return nil, err
	}

	return schedule, nil
}

func (s *SlotService) DeleteRecurringSchedule(ctx context.Context, id string) error {
	return s.db.WithContext(ctx).Delete(&models.RecurringSchedule{}, "id = ?", id).Error
}
