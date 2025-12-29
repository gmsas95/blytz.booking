package services

import (
	"blytz.cloud/backend/internal/models"
	"blytz.cloud/backend/internal/repository"
	"context"
	"gorm.io/gorm"
)

type BusinessService struct {
	businessRepo repository.BusinessRepository
	db           *gorm.DB
}

func NewBusinessService(
	businessRepo repository.BusinessRepository,
	db *gorm.DB,
) *BusinessService {
	return &BusinessService{
		businessRepo: businessRepo,
		db:           db,
	}
}

type UpdateBusinessRequest struct {
	Name        *string `json:"name" validate:"omitempty,min=3"`
	Description *string `json:"description"`
	ThemeColor  *string `json:"theme_color" validate:"omitempty,hexcolor"`
	LogoURL     *string `json:"logo_url" validate:"omitempty,url"`
}

type BusinessSettingsUpdateRequest struct {
	RequireDeposit          *bool   `json:"require_deposit"`
	CancellationPolicyHours *int    `json:"cancellation_policy_hours" validate:"omitempty,min=0"`
	ConfirmationEmail       *bool   `json:"confirmation_email"`
	ReminderEmail           *bool   `json:"reminder_email"`
	ReminderHoursBefore     *int    `json:"reminder_hours_before" validate:"omitempty,min=1"`
	Timezone                *string `json:"timezone"`
	BookingBufferMinutes    *int    `json:"booking_buffer_minutes" validate:"omitempty,min=0"`
}

func (s *BusinessService) GetBySlug(ctx context.Context, slug string) (*models.Business, error) {
	business, err := s.businessRepo.GetBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}

	return business, nil
}

func (s *BusinessService) UpdateBusiness(ctx context.Context, id string, req *UpdateBusinessRequest) (*models.Business, error) {
	business, err := s.businessRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		business.Name = *req.Name
	}
	if req.Description != nil {
		business.Description = *req.Description
	}
	if req.ThemeColor != nil {
		business.ThemeColor = *req.ThemeColor
	}
	if req.LogoURL != nil {
		business.LogoURL = req.LogoURL
	}

	if err := s.businessRepo.Update(ctx, business); err != nil {
		return nil, err
	}

	return business, nil
}

func (s *BusinessService) DeleteBusiness(ctx context.Context, id string) error {
	return s.businessRepo.Delete(ctx, id)
}

func (s *BusinessService) UpdateSettings(ctx context.Context, id string, req *BusinessSettingsUpdateRequest) (*models.BusinessSettings, error) {
	var settings models.BusinessSettings
	err := s.db.WithContext(ctx).Where("business_id = ?", id).First(&settings).Error
	if err != nil {
		return nil, err
	}

	if req.RequireDeposit != nil {
		settings.RequireDeposit = *req.RequireDeposit
	}
	if req.CancellationPolicyHours != nil {
		settings.CancellationPolicyHours = *req.CancellationPolicyHours
	}
	if req.ConfirmationEmail != nil {
		settings.ConfirmationEmail = *req.ConfirmationEmail
	}
	if req.ReminderEmail != nil {
		settings.ReminderEmail = *req.ReminderEmail
	}
	if req.ReminderHoursBefore != nil {
		settings.ReminderHoursBefore = *req.ReminderHoursBefore
	}
	if req.Timezone != nil {
		settings.Timezone = *req.Timezone
	}
	if req.BookingBufferMinutes != nil {
		settings.BookingBufferMinutes = *req.BookingBufferMinutes
	}

	if err := s.db.WithContext(ctx).Save(&settings).Error; err != nil {
		return nil, err
	}

	return &settings, nil
}

func (s *BusinessService) ListBusinesses(ctx context.Context, offset, limit int) ([]*models.Business, int64, error) {
	return s.businessRepo.List(ctx, offset, limit)
}
