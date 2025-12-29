package services

import (
	"context"
	"errors"
	"time"

	"blytz.cloud/backend/internal/models"
	"blytz.cloud/backend/internal/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BookingService struct {
	bookingRepo  repository.BookingRepository
	slotRepo     repository.SlotRepository
	serviceRepo  repository.ServiceRepository
	customerRepo repository.CustomerRepository
	historyRepo  repository.BookingHistoryRepository
	db           *gorm.DB
}

func NewBookingService(
	bookingRepo repository.BookingRepository,
	slotRepo repository.SlotRepository,
	serviceRepo repository.ServiceRepository,
	customerRepo repository.CustomerRepository,
	historyRepo repository.BookingHistoryRepository,
	db *gorm.DB,
) *BookingService {
	return &BookingService{
		bookingRepo:  bookingRepo,
		slotRepo:     slotRepo,
		serviceRepo:  serviceRepo,
		customerRepo: customerRepo,
		historyRepo:  historyRepo,
		db:           db,
	}
}

type CreateBookingRequest struct {
	BusinessID string                 `json:"business_id" validate:"required,uuid"`
	ServiceID  string                 `json:"service_id" validate:"required,uuid"`
	SlotID     string                 `json:"slot_id" validate:"required,uuid"`
	Customer   models.CustomerDetails `json:"customer" validate:"required"`
	Notes      string                 `json:"notes"`
}

type BookingResponse struct {
	Booking *models.Booking `json:"booking"`
}

func (s *BookingService) CreateBooking(ctx context.Context, req *CreateBookingRequest) (*models.Booking, error) {
	var booking *models.Booking

	err := s.db.Transaction(func(tx *gorm.DB) error {
		slotRepo := s.slotRepo.WithTx(tx)
		serviceRepo := s.serviceRepo.WithTx(tx)
		bookingRepo := s.bookingRepo.WithTx(tx)
		customerRepo := s.customerRepo.WithTx(tx)

		slot, err := slotRepo.GetByID(ctx, req.SlotID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("slot not found")
			}
			return err
		}

		if slot.BusinessID.String() != req.BusinessID {
			return errors.New("slot belongs to different business")
		}

		if slot.IsBooked {
			return errors.New("slot is no longer available")
		}

		if slot.BookedCount >= slot.Capacity {
			return errors.New("slot is at full capacity")
		}

		service, err := serviceRepo.GetByID(ctx, req.ServiceID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("service not found")
			}
			return err
		}

		if service.BusinessID.String() != req.BusinessID {
			return errors.New("service belongs to different business")
		}

		customer := &models.Customer{
			ID:         uuid.New(),
			BusinessID: slot.BusinessID,
			Name:       req.Customer.Name,
			Email:      req.Customer.Email,
			Phone:      req.Customer.Phone,
		}

		customer, err = customerRepo.GetOrCreate(ctx, customer)
		if err != nil {
			return err
		}

		booking = &models.Booking{
			ID:              uuid.New(),
			BusinessID:      slot.BusinessID,
			ServiceID:       service.ID,
			SlotID:          slot.ID,
			CustomerID:      &customer.ID,
			ServiceName:     service.Name,
			SlotTime:        slot.StartTime,
			CustomerDetails: req.Customer,
			Status:          models.BookingStatusPending,
			DepositPaid:     service.DepositAmount,
			TotalPrice:      service.TotalPrice,
			Notes:           req.Notes,
			PaymentStatus:   models.PaymentStatusPending,
		}

		customer, err = customerRepo.GetOrCreate(ctx, customer)
		if err != nil {
			return err
		}

		booking = &models.Booking{
			ID:            uuid.New(),
			BusinessID:    slot.BusinessID,
			ServiceID:     service.ID,
			SlotID:        slot.ID,
			CustomerID:    &customer.ID,
			ServiceName:   service.Name,
			SlotTime:      slot.StartTime,
			Status:        models.BookingStatusPending,
			DepositPaid:   service.DepositAmount,
			TotalPrice:    service.TotalPrice,
			Notes:         req.Notes,
			PaymentStatus: models.PaymentStatusPending,
		}

		if err := bookingRepo.Create(ctx, booking); err != nil {
			return err
		}

		if err := slotRepo.MarkBooked(ctx, req.SlotID); err != nil {
			return err
		}

		history := &models.BookingHistory{
			ID:          uuid.New(),
			BookingID:   booking.ID,
			Action:      "created",
			Previous:    "",
			Current:     "",
			PerformedBy: uuid.Nil,
			Timestamp:   time.Now(),
		}

		return s.historyRepo.Create(ctx, history)
	})

	if err != nil {
		return nil, err
	}

	return booking, nil
}

func (s *BookingService) GetBooking(ctx context.Context, id string) (*models.Booking, error) {
	return s.bookingRepo.GetByID(ctx, id)
}

func (s *BookingService) ListBookings(ctx context.Context, businessID string, offset, limit int) ([]*models.Booking, int64, error) {
	return s.bookingRepo.ListByBusiness(ctx, businessID, offset, limit)
}

func (s *BookingService) UpdateBookingStatus(ctx context.Context, id string, status models.BookingStatus, userID string) (*models.Booking, error) {
	var booking *models.Booking

	err := s.db.Transaction(func(tx *gorm.DB) error {
		bookingRepo := s.bookingRepo.WithTx(tx)
		historyRepo := s.historyRepo.WithTx(tx)

		var err error
		booking, err = bookingRepo.GetByID(ctx, id)
		if err != nil {
			return err
		}

		previousStatus := booking.Status
		if err := bookingRepo.UpdateStatus(ctx, id, status); err != nil {
			return err
		}

		history := &models.BookingHistory{
			ID:          uuid.New(),
			BookingID:   booking.ID,
			Action:      "status_changed",
			Previous:    string(previousStatus),
			Current:     string(status),
			PerformedBy: uuid.Nil,
			Timestamp:   time.Now(),
		}

		return historyRepo.Create(ctx, history)
	})

	if err != nil {
		return nil, err
	}

	return booking, nil
}

func (s *BookingService) CancelBooking(ctx context.Context, id, reason string, userID string) error {
	_, err := s.UpdateBookingStatus(ctx, id, models.BookingStatusCancelled, userID)
	return err
}
