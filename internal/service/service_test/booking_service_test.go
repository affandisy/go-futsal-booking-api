package service_test

import (
	"context"
	"errors"
	"go-futsal-booking-api/internal/domain"
	"go-futsal-booking-api/internal/dto/request"
	"go-futsal-booking-api/internal/repository/mock"
	"go-futsal-booking-api/internal/service"
	"go-futsal-booking-api/pkg/logger"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func init() {
	// Initialize logger for testing
	logger.Init("test")
}

func TestBookingService_CreateBooking(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBookingRepo := mock.NewMockBookingRepository(ctrl)
	mockScheduleRepo := mock.NewMockScheduleRepository(ctrl)
	mockUserRepo := mock.NewMockUserRepository(ctrl)

	bookingService := service.NewBookingService(mockBookingRepo, mockScheduleRepo, mockUserRepo)

	t.Run("Success - Create booking", func(t *testing.T) {
		ctx := context.Background()

		// Tomorrow's date
		tomorrow := time.Now().Add(24 * time.Hour)
		bookingDate := tomorrow.Format("2006-01-02")
		dayOfWeek := int(tomorrow.Weekday())
		if dayOfWeek == 0 {
			dayOfWeek = 7
		}

		req := &request.CreateBookingRequest{
			UserID:      1,
			ScheduleID:  1,
			BookingDate: bookingDate,
		}

		schedule := domain.Schedule{
			ID:        1,
			DayOfWeek: dayOfWeek,
			StartTime: time.Date(0, 1, 1, 10, 0, 0, 0, time.UTC),
			EndTime:   time.Date(0, 1, 1, 12, 0, 0, 0, time.UTC),
			Price:     100000,
			Field: domain.Field{
				ID:   1,
				Name: "Field A",
				Type: "Futsal",
				Venue: domain.Venue{
					ID:      1,
					Name:    "Venue 1",
					Address: "123 Main St",
					City:    "Jakarta",
				},
			},
		}

		user := domain.User{
			ID:       1,
			FullName: "John Doe",
			Email:    "john@example.com",
			Role: domain.Role{
				ID:       2,
				RoleName: "customer",
			},
		}

		mockScheduleRepo.EXPECT().
			FindByID(ctx, req.ScheduleID).
			Return(schedule, nil)

		mockUserRepo.EXPECT().
			FindByID(ctx, req.UserID).
			Return(user, nil)

		mockBookingRepo.EXPECT().
			Create(ctx, gomock.Any()).
			DoAndReturn(func(ctx context.Context, booking *domain.Booking) error {
				booking.ID = 1
				booking.CreatedAt = time.Now()
				return nil
			})

		result, err := bookingService.CreateBooking(ctx, req)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, uint(1), result.ID)
		assert.Equal(t, "PENDING", result.Status)
		assert.Equal(t, schedule.Price, result.TotalPrice)
	})

	t.Run("Fail - Invalid booking request (nil)", func(t *testing.T) {
		ctx := context.Background()

		result, err := bookingService.CreateBooking(ctx, nil)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, "invalid booking request", err.Error())
	})

	t.Run("Fail - Invalid booking request (missing fields)", func(t *testing.T) {
		ctx := context.Background()

		req := &request.CreateBookingRequest{
			UserID:      0,
			ScheduleID:  1,
			BookingDate: "2024-12-01",
		}

		result, err := bookingService.CreateBooking(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, "invalid booking request", err.Error())
	})

	t.Run("Fail - Invalid date format", func(t *testing.T) {
		ctx := context.Background()

		req := &request.CreateBookingRequest{
			UserID:      1,
			ScheduleID:  1,
			BookingDate: "invalid-date",
		}

		result, err := bookingService.CreateBooking(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, domain.ErrInvalidBookingDate, err)
	})

	t.Run("Fail - Cannot book past date", func(t *testing.T) {
		ctx := context.Background()

		yesterday := time.Now().Add(-24 * time.Hour)
		pastDate := yesterday.Format("2006-01-02")

		req := &request.CreateBookingRequest{
			UserID:      1,
			ScheduleID:  1,
			BookingDate: pastDate,
		}

		result, err := bookingService.CreateBooking(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, domain.ErrPastDateBooking, err)
	})

	t.Run("Fail - Schedule not found", func(t *testing.T) {
		ctx := context.Background()

		tomorrow := time.Now().Add(24 * time.Hour)
		bookingDate := tomorrow.Format("2006-01-02")

		req := &request.CreateBookingRequest{
			UserID:      1,
			ScheduleID:  999,
			BookingDate: bookingDate,
		}

		mockScheduleRepo.EXPECT().
			FindByID(ctx, req.ScheduleID).
			Return(domain.Schedule{}, domain.ErrScheduleNotFound)

		result, err := bookingService.CreateBooking(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, domain.ErrScheduleNotFound, err)
	})

	t.Run("Fail - Day mismatch", func(t *testing.T) {
		ctx := context.Background()

		tomorrow := time.Now().Add(24 * time.Hour)
		bookingDate := tomorrow.Format("2006-01-02")
		dayOfWeek := int(tomorrow.Weekday())
		if dayOfWeek == 0 {
			dayOfWeek = 7
		}

		// Schedule with different day
		wrongDay := dayOfWeek + 1
		if wrongDay > 7 {
			wrongDay = 1
		}

		req := &request.CreateBookingRequest{
			UserID:      1,
			ScheduleID:  1,
			BookingDate: bookingDate,
		}

		schedule := domain.Schedule{
			ID:        1,
			DayOfWeek: wrongDay, // Different day
			Price:     100000,
		}

		mockScheduleRepo.EXPECT().
			FindByID(ctx, req.ScheduleID).
			Return(schedule, nil)

		result, err := bookingService.CreateBooking(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, domain.ErrDayMistmatch, err)
	})

	t.Run("Fail - User not found", func(t *testing.T) {
		ctx := context.Background()

		tomorrow := time.Now().Add(24 * time.Hour)
		bookingDate := tomorrow.Format("2006-01-02")
		dayOfWeek := int(tomorrow.Weekday())
		if dayOfWeek == 0 {
			dayOfWeek = 7
		}

		req := &request.CreateBookingRequest{
			UserID:      999,
			ScheduleID:  1,
			BookingDate: bookingDate,
		}

		schedule := domain.Schedule{
			ID:        1,
			DayOfWeek: dayOfWeek,
			Price:     100000,
		}

		mockScheduleRepo.EXPECT().
			FindByID(ctx, req.ScheduleID).
			Return(schedule, nil)

		mockUserRepo.EXPECT().
			FindByID(ctx, req.UserID).
			Return(domain.User{}, errors.New("user not found"))

		result, err := bookingService.CreateBooking(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "user not found")
	})

	t.Run("Fail - Database error on create", func(t *testing.T) {
		ctx := context.Background()

		tomorrow := time.Now().Add(24 * time.Hour)
		bookingDate := tomorrow.Format("2006-01-02")
		dayOfWeek := int(tomorrow.Weekday())
		if dayOfWeek == 0 {
			dayOfWeek = 7
		}

		req := &request.CreateBookingRequest{
			UserID:      1,
			ScheduleID:  1,
			BookingDate: bookingDate,
		}

		schedule := domain.Schedule{
			ID:        1,
			DayOfWeek: dayOfWeek,
			Price:     100000,
		}

		user := domain.User{
			ID:       1,
			FullName: "John Doe",
		}

		mockScheduleRepo.EXPECT().
			FindByID(ctx, req.ScheduleID).
			Return(schedule, nil)

		mockUserRepo.EXPECT().
			FindByID(ctx, req.UserID).
			Return(user, nil)

		mockBookingRepo.EXPECT().
			Create(ctx, gomock.Any()).
			Return(errors.New("database error"))

		result, err := bookingService.CreateBooking(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "failed to create booking")
	})
}

func TestBookingService_GetMyBookings(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBookingRepo := mock.NewMockBookingRepository(ctrl)
	mockScheduleRepo := mock.NewMockScheduleRepository(ctrl)
	mockUserRepo := mock.NewMockUserRepository(ctrl)

	bookingService := service.NewBookingService(mockBookingRepo, mockScheduleRepo, mockUserRepo)

	t.Run("Success - Get user bookings", func(t *testing.T) {
		ctx := context.Background()
		userID := uint(1)

		booking := &domain.Booking{
			ID:          1,
			BookingDate: time.Now(),
			Status:      "PENDING",
			TotalPrice:  100000,
			User: domain.User{
				ID:       userID,
				FullName: "John Doe",
			},
			Schedule: domain.Schedule{
				ID:    1,
				Price: 100000,
			},
		}

		mockBookingRepo.EXPECT().
			FindByUserID(ctx, userID).
			Return(booking, nil)

		result, err := bookingService.GetMyBookings(ctx, userID)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result, 1)
		assert.Equal(t, booking.ID, result[0].ID)
	})

	t.Run("Fail - Invalid user ID", func(t *testing.T) {
		ctx := context.Background()

		result, err := bookingService.GetMyBookings(ctx, 0)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, "invalid user id", err.Error())
	})

	t.Run("Fail - Database error", func(t *testing.T) {
		ctx := context.Background()
		userID := uint(1)

		mockBookingRepo.EXPECT().
			FindByUserID(ctx, userID).
			Return(nil, errors.New("database error"))

		result, err := bookingService.GetMyBookings(ctx, userID)

		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestBookingService_GetBookingByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBookingRepo := mock.NewMockBookingRepository(ctrl)
	mockScheduleRepo := mock.NewMockScheduleRepository(ctrl)
	mockUserRepo := mock.NewMockUserRepository(ctrl)

	bookingService := service.NewBookingService(mockBookingRepo, mockScheduleRepo, mockUserRepo)

	t.Run("Success - Get booking by ID", func(t *testing.T) {
		ctx := context.Background()
		bookingID := uint(1)

		booking := domain.Booking{
			ID:          bookingID,
			BookingDate: time.Now(),
			Status:      "PENDING",
			TotalPrice:  100000,
			User: domain.User{
				ID:       1,
				FullName: "John Doe",
			},
			Schedule: domain.Schedule{
				ID:    1,
				Price: 100000,
			},
		}

		mockBookingRepo.EXPECT().
			FindByID(ctx, bookingID).
			Return(booking, nil)

		result, err := bookingService.GetBookingByID(ctx, bookingID)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, bookingID, result.ID)
	})

	t.Run("Fail - Invalid booking ID", func(t *testing.T) {
		ctx := context.Background()

		result, err := bookingService.GetBookingByID(ctx, 0)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, "invalid booking id", err.Error())
	})

	t.Run("Fail - Booking not found", func(t *testing.T) {
		ctx := context.Background()
		bookingID := uint(999)

		mockBookingRepo.EXPECT().
			FindByID(ctx, bookingID).
			Return(domain.Booking{}, domain.ErrBookingNotFound)

		result, err := bookingService.GetBookingByID(ctx, bookingID)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, domain.ErrBookingNotFound, err)
	})
}

func TestBookingService_CancelBooking(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBookingRepo := mock.NewMockBookingRepository(ctrl)
	mockScheduleRepo := mock.NewMockScheduleRepository(ctrl)
	mockUserRepo := mock.NewMockUserRepository(ctrl)

	bookingService := service.NewBookingService(mockBookingRepo, mockScheduleRepo, mockUserRepo)

	t.Run("Success - Cancel booking", func(t *testing.T) {
		ctx := context.Background()
		bookingID := uint(1)
		userID := uint(1)

		booking := domain.Booking{
			ID:          bookingID,
			BookingDate: time.Now(),
			Status:      "PENDING",
			TotalPrice:  100000,
			User: domain.User{
				ID:       userID,
				FullName: "John Doe",
			},
		}

		mockBookingRepo.EXPECT().
			FindByID(ctx, bookingID).
			Return(booking, nil)

		mockBookingRepo.EXPECT().
			CancelBooking(ctx, bookingID).
			Return(nil)

		err := bookingService.CancelBooking(ctx, bookingID, userID)

		assert.NoError(t, err)
	})

	t.Run("Fail - Invalid booking ID", func(t *testing.T) {
		ctx := context.Background()

		err := bookingService.CancelBooking(ctx, 0, 1)

		assert.Error(t, err)
		assert.Equal(t, "invalid booking or user id", err.Error())
	})

	t.Run("Fail - Invalid user ID", func(t *testing.T) {
		ctx := context.Background()

		err := bookingService.CancelBooking(ctx, 1, 0)

		assert.Error(t, err)
		assert.Equal(t, "invalid booking or user id", err.Error())
	})

	t.Run("Fail - Booking not found", func(t *testing.T) {
		ctx := context.Background()
		bookingID := uint(999)
		userID := uint(1)

		mockBookingRepo.EXPECT().
			FindByID(ctx, bookingID).
			Return(domain.Booking{}, domain.ErrBookingNotFound)

		err := bookingService.CancelBooking(ctx, bookingID, userID)

		assert.Error(t, err)
		assert.Equal(t, domain.ErrBookingNotFound, err)
	})

	t.Run("Fail - Forbidden (wrong user)", func(t *testing.T) {
		ctx := context.Background()
		bookingID := uint(1)
		userID := uint(2)

		booking := domain.Booking{
			ID:     bookingID,
			Status: "PENDING",
			User: domain.User{
				ID: 1, // Different user
			},
		}

		mockBookingRepo.EXPECT().
			FindByID(ctx, bookingID).
			Return(booking, nil)

		err := bookingService.CancelBooking(ctx, bookingID, userID)

		assert.Error(t, err)
		assert.Equal(t, domain.ErrForbidden, err)
	})

	t.Run("Fail - Cannot cancel cancelled booking", func(t *testing.T) {
		ctx := context.Background()
		bookingID := uint(1)
		userID := uint(1)

		booking := domain.Booking{
			ID:     bookingID,
			Status: "CANCELLED",
			User: domain.User{
				ID: userID,
			},
		}

		mockBookingRepo.EXPECT().
			FindByID(ctx, bookingID).
			Return(booking, nil)

		err := bookingService.CancelBooking(ctx, bookingID, userID)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cannot cancel booking with status")
	})

	t.Run("Fail - Cannot cancel completed booking", func(t *testing.T) {
		ctx := context.Background()
		bookingID := uint(1)
		userID := uint(1)

		booking := domain.Booking{
			ID:     bookingID,
			Status: "COMPLETED",
			User: domain.User{
				ID: userID,
			},
		}

		mockBookingRepo.EXPECT().
			FindByID(ctx, bookingID).
			Return(booking, nil)

		err := bookingService.CancelBooking(ctx, bookingID, userID)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cannot cancel booking with status")
	})

	t.Run("Fail - Database error on cancel", func(t *testing.T) {
		ctx := context.Background()
		bookingID := uint(1)
		userID := uint(1)

		booking := domain.Booking{
			ID:     bookingID,
			Status: "PENDING",
			User: domain.User{
				ID: userID,
			},
		}

		mockBookingRepo.EXPECT().
			FindByID(ctx, bookingID).
			Return(booking, nil)

		mockBookingRepo.EXPECT().
			CancelBooking(ctx, bookingID).
			Return(errors.New("database error"))

		err := bookingService.CancelBooking(ctx, bookingID, userID)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to cancel booking")
	})
}
