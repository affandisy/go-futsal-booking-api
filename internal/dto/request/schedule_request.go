package request

type CreateScheduleRequest struct {
	FieldID   uint    `json:"field_id" validate:"required"`
	DayOfWeek int     `json:"day_of_week" validate:"required"`
	StartTime string  `json:"start_time" validate:"required"`
	EndTime   string  `json:"end_time" validate:"required"`
	Price     float64 `json:"price" validate:"required"`
}

type UpdateScheduleRequest struct {
	FieldID   uint    `json:"field_id" validate:"required"`
	DayOfWeek int     `json:"day_of_week" validate:"required"`
	StartTime string  `json:"start_time" validate:"required"`
	EndTime   string  `json:"end_time" validate:"required"`
	Price     float64 `json:"price" validate:"required"`
}
