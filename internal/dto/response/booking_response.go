package response

import (
	"go-futsal-booking-api/internal/domain"
	"time"
)

type BookingResponse struct {
	ID          uint      `json:"id"`
	UserID      uint      `json:"user_id"`
	ScheduleID  uint      `json:"schedule_id"`
	BookingDate time.Time `json:"booking_date"`
	Status      string    `json:"status"`
	TotalPrice  float64   `json:"total_price"`
	CreatedAt   time.Time `json:"created_at"`
}

func ToBookingResponse(booking *domain.Booking) BookingResponse {
	return BookingResponse{
		ID:          booking.ID,
		UserID:      booking.User.ID,
		ScheduleID:  booking.Schedule.ID,
		BookingDate: booking.BookingDate,
		Status:      booking.Status,
		TotalPrice:  booking.TotalPrice,
		CreatedAt:   booking.CreatedAt,
	}
}
