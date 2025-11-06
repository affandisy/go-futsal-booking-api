package service

import (
	"context"
	"errors"
	"go-futsal-booking-api/internal/domain"
	"go-futsal-booking-api/internal/repository"
	"time"
)

type BookingService interface {
	CreateBooking(ctx context.Context, customer domain.User, scheduleID uint, bookingDate string) (domain.Booking, error)
	GetMyBookings(ctx context.Context, customer domain.User) (*domain.Booking, error)
	GetBookingDetails(ctx context.Context, bookingID uint) (domain.Booking, error)
}

type bookingService struct {
	bookingRepo  repository.BookingRepository
	scheduleRepo repository.ScheduleRepository
}

func NewBookingService(bookingRepo repository.BookingRepository, scheduleRepo repository.ScheduleRepository) BookingService {
	return &bookingService{
		bookingRepo:  bookingRepo,
		scheduleRepo: scheduleRepo,
	}
}

func parseDate(dateStr string) (time.Time, error) {
	return time.Parse("2006-01-02", dateStr)
}

func (s *bookingService) CreateBooking(ctx context.Context, customer domain.User, scheduleID uint, bookingDate string) (domain.Booking, error) {
	bookDate, err := parseDate(bookingDate)
	if err != nil {
		return domain.Booking{}, errors.New("invalide date format")
	}

	schedule, err := s.scheduleRepo.FindByID(ctx, scheduleID)
	if err != nil {
		return domain.Booking{}, err
	}

	newBooking := domain.Booking{
		User:        customer,
		Schedule:    schedule,
		BookingDate: bookDate,
		Status:      "PENDING",
		TotalPrice:  schedule.Price,
	}

	if err := s.bookingRepo.Create(ctx, &newBooking); err != nil {
		return domain.Booking{}, err
	}

	return newBooking, nil
}

func (s *bookingService) GetMyBookings(ctx context.Context, customer domain.User) (*domain.Booking, error) {
	bookings, err := s.bookingRepo.FindByUserID(ctx, customer.ID)
	if err != nil {
		return nil, err
	}

	return bookings, nil
}

func (s *bookingService) GetBookingDetails(ctx context.Context, bookingID uint) (domain.Booking, error) {
	booking, err := s.bookingRepo.FindByID(ctx, bookingID)
	if err != nil {
		return domain.Booking{}, err
	}

	return booking, nil
}
