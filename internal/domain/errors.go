package domain

import "errors"

var (
	ErrForbidden        = errors.New("forbidden: user does not have the required permissions")
	ErrBookingNotFound  = errors.New("booking not found")
	ErrSlotUnavailable  = errors.New("slot already booked for this date")
	ErrScheduleNotFound = errors.New("schedule not found")
)
