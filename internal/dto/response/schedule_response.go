package response

import (
	"go-futsal-booking-api/internal/domain"
	"time"
)

type ScheduleResponse struct {
	ID        uint      `json:"id"`
	FieldID   uint      `json:"field_id"`
	DayOfWeek int       `json:"day_of_week"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	Price     float64   `json:"price"`
	CreatedAt time.Time `json:"created_at"`
}

func ToScheduleResponse(schedule *domain.Schedule) ScheduleResponse {
	return ScheduleResponse{
		ID:        schedule.ID,
		FieldID:   schedule.Field.ID,
		DayOfWeek: schedule.DayOfWeek,
		StartTime: schedule.StartTime,
		EndTime:   schedule.EndTime,
		Price:     schedule.Price,
		CreatedAt: schedule.CreatedAt,
	}
}
