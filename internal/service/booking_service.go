package service

import (
	"context"
	"errors"
	"fmt"
	"go-futsal-booking-api/internal/domain"
	"go-futsal-booking-api/internal/dto/request"
	"go-futsal-booking-api/internal/repository"
	"go-futsal-booking-api/pkg/logger"
	"time"
)

type BookingService interface {
	CreateBooking(ctx context.Context, req *request.CreateBookingRequest) (*domain.Booking, error)
	GetMyBookings(ctx context.Context, userID uint) ([]*domain.Booking, error)
	GetBookingByID(ctx context.Context, bookingID uint) (*domain.Booking, error)
	CancelBooking(ctx context.Context, bookingID uint, userID uint) error
}

type bookingService struct {
	bookingRepo  repository.BookingRepository
	scheduleRepo repository.ScheduleRepository
	userRepo     repository.UserRepository
}

// type CreateBookingRequest struct {
// 	UserID      uint
// 	ScheduleID  uint
// 	BookingDate string
// }

func NewBookingService(bookingRepo repository.BookingRepository, scheduleRepo repository.ScheduleRepository, userRepo repository.UserRepository) BookingService {
	return &bookingService{
		bookingRepo:  bookingRepo,
		scheduleRepo: scheduleRepo,
		userRepo:     userRepo,
	}
}

func parseDate(dateStr string) (time.Time, error) {
	return time.Parse("2006-01-02", dateStr)
}

func (s *bookingService) CreateBooking(ctx context.Context, req *request.CreateBookingRequest) (*domain.Booking, error) {
	if req == nil || req.UserID == 0 || req.ScheduleID == 0 || req.BookingDate == "" {
		return nil, errors.New("invalid booking request")
	}

	if err := ctx.Err(); err != nil {
		return nil, fmt.Errorf("context error: %w", err)
	}

	bookDate, err := parseDate(req.BookingDate)
	if err != nil {
		logger.Error("Invalid date format", err.Error())
		return nil, domain.ErrInvalidBookingDate
	}

	// cannot booking in the past
	now := time.Now()
	if bookDate.Before(now.Truncate(24 * time.Hour)) {
		logger.Error("Attempt to booking past date")
		return nil, domain.ErrPastDateBooking
	}

	schedule, err := s.scheduleRepo.FindByID(ctx, req.ScheduleID)
	if err != nil {
		logger.Error("schedule not found", err.Error())
		return nil, domain.ErrScheduleNotFound
	}

	bookingDayOfWeek := int(bookDate.Weekday())
	if bookingDayOfWeek == 0 {
		bookingDayOfWeek = 7
	}

	if schedule.DayOfWeek != bookingDayOfWeek {
		logger.Warn("day mistmatch", map[string]any{
			"booking_date":        bookDate,
			"booking_day_of_week": bookingDayOfWeek,
			"schedule_day":        schedule.DayOfWeek,
		})
		return nil, domain.ErrDayMistmatch
	}

	user, err := s.userRepo.FindByID(ctx, req.UserID)
	if err != nil {
		logger.Error("user not found", err.Error())
		return nil, errors.New("user not found")
	}

	newBooking := &domain.Booking{
		User:        user,
		Schedule:    schedule,
		BookingDate: bookDate,
		Status:      "PENDING",
		TotalPrice:  schedule.Price,
	}

	if err := s.bookingRepo.Create(ctx, newBooking); err != nil {
		logger.Error("failed to create booking", err.Error())
		return nil, fmt.Errorf("failed to create booking: %w", err)
	}

	logger.Info("booking created successfully", map[string]any{
		"booking_id": newBooking.ID,
	})

	return newBooking, nil
}

func (s *bookingService) GetMyBookings(ctx context.Context, userID uint) ([]*domain.Booking, error) {
	if userID == 0 {
		return nil, errors.New("invalid user id")
	}

	if err := ctx.Err(); err != nil {
		return nil, fmt.Errorf("context error: %w", err)
	}

	booking, err := s.bookingRepo.FindByUserID(ctx, userID)
	if err != nil {
		logger.Error("failed to get user bookings", err.Error())
		return nil, err
	}

	return []*domain.Booking{booking}, nil
}

func (s *bookingService) GetBookingByID(ctx context.Context, bookingID uint) (*domain.Booking, error) {
	if bookingID == 0 {
		return nil, errors.New("invalid booking id")
	}

	if err := ctx.Err(); err != nil {
		return nil, fmt.Errorf("context error: %w", err)
	}

	booking, err := s.bookingRepo.FindByID(ctx, bookingID)
	if err != nil {
		logger.Error("booking not found", err.Error())
		return nil, domain.ErrBookingNotFound
	}

	return &booking, nil
}

func (s *bookingService) CancelBooking(ctx context.Context, bookingID uint, userID uint) error {
	if bookingID == 0 || userID == 0 {
		return errors.New("invalid booking or user id")
	}

	if err := ctx.Err(); err != nil {
		return fmt.Errorf("context error: %w", err)
	}

	booking, err := s.bookingRepo.FindByID(ctx, bookingID)
	if err != nil {
		return domain.ErrBookingNotFound
	}

	if booking.User.ID != userID {
		return domain.ErrForbidden
	}

	if booking.Status == "CANCELLED" || booking.Status == "COMPLETED" {
		return errors.New("cannot cancel booking with status: " + booking.Status)
	}

	if err := s.bookingRepo.CancelBooking(ctx, bookingID); err != nil {
		logger.Error("failed to cancel booking", err.Error())
		return fmt.Errorf("failed to cancel booking: %w", err)
	}

	logger.Info("booking cancelled", map[string]any{
		"booking_id": bookingID,
	})

	return nil
}
