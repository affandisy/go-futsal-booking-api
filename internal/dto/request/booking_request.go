package request

type CreateBookingRequest struct {
	// UserID      uint    `json:"user_id" validate:"required"`
	ScheduleID  uint   `json:"schedule_id" validate:"required"`
	BookingDate string `json:"booking_date" validate:"required"`
	// Status      string  `json:"status"`
	// TotalPrice  float64 `json:"total_price" validate:"required"`
}

type UpdateBookingRequest struct {
	// UserID      uint    `json:"user_id" validate:"required"`
	ScheduleID  uint   `json:"schedule_id" validate:"required"`
	BookingDate string `json:"booking_date" validate:"required"`
	// Status      string  `json:"status"`
	// TotalPrice  float64 `json:"total_price" validate:"required"`
}
