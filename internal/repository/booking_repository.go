package repository

import (
	"context"
	"errors"
	"go-futsal-booking-api/internal/domain"
	gormContract "go-futsal-booking-api/internal/repository/model"

	"gorm.io/gorm"
)

type BookingRepository interface {
	Create(ctx context.Context, booking *domain.Booking) error
	FindByID(ctx context.Context, id uint) (domain.Booking, error)
	FindByUserID(ctx context.Context, userID uint) (*domain.Booking, error)
}

type gormBookingRepository struct {
	DB *gorm.DB
}

func NewBookingRepository(db *gorm.DB) BookingRepository {
	return &gormBookingRepository{DB: db}
}

func (r *gormBookingRepository) preload(ctx context.Context) *gorm.DB {
	return r.DB.WithContext(ctx).Preload("User.Role").Preload("Schedule.Field.Venue")
}

func (r *gormBookingRepository) Create(ctx context.Context, booking *domain.Booking) error {
	var gormBooking gormContract.BookingGorm
	gormBooking.FromDomain(*booking)

	err := r.DB.WithContext(ctx).Create(&gormBooking).Error
	if err != nil {
		return err
	}

	if err := r.preload(ctx).First(&gormBooking, gormBooking.ID).Error; err != nil {
		return err
	}

	*booking = gormBooking.ToDomain()

	return nil
}

func (r *gormBookingRepository) FindByID(ctx context.Context, id uint) (domain.Booking, error) {
	var gormBooking gormContract.BookingGorm

	err := r.preload(ctx).First(&gormBooking, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.Booking{}, domain.ErrBookingNotFound
		}
		return domain.Booking{}, err
	}

	booking := gormBooking.ToDomain()

	return booking, nil
}

func (r *gormBookingRepository) FindByUserID(ctx context.Context, userID uint) (*domain.Booking, error) {
	var gormBookings []gormContract.BookingGorm

	err := r.preload(ctx).Where("user_id = ?", userID).Find(&gormBookings).Error
	if err != nil {
		return nil, err
	}

	if len(gormBookings) == 0 {
		return nil, domain.ErrBookingNotFound
	}

	booking := gormBookings[0].ToDomain()
	return &booking, nil
}
