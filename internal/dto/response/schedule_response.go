package response

import (
	"go-futsal-booking-api/internal/domain"
	"time"
)

type ScheduleResponse struct {
	ID        uint      `json:"id"`
	FieldID   uint      `json:"field_id"`
	DayOfWeek int       `json:"day_of_week"`
	StartTime string    `json:"start_time"`
	EndTime   string    `json:"end_time"`
	Price     float64   `json:"price"`
	CreatedAt time.Time `json:"created_at"`
}

func ToScheduleResponse(schedule *domain.Schedule) ScheduleResponse {
	return ScheduleResponse{
		ID:        schedule.ID,
		FieldID:   schedule.Field.ID,
		DayOfWeek: schedule.DayOfWeek,
		StartTime: schedule.StartTime.Format("15:04:05"),
		EndTime:   schedule.EndTime.Format("15:04:05"),
		Price:     schedule.Price,
		CreatedAt: schedule.CreatedAt,
	}
}
