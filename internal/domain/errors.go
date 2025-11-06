package domain

import "errors"

var (
	ErrForbidden            = errors.New("forbidden: user does not have the required permissions")
	ErrBookingNotFound      = errors.New("booking not found")
	ErrSlotUnavailable      = errors.New("slot already booked for this date")
	ErrScheduleNotFound     = errors.New("schedule not found")
	ErrSlotAlreadyBooked    = errors.New("slot already booked for this date")
	ErrInvalidBookingDate   = errors.New("invalid booking date")
	ErrDayMistmatch         = errors.New("booking date does not match schedule day")
	ErrPastDateBooking      = errors.New("cannot book past date")
	ErrScheduleNotAvailable = errors.New("schedule is not available")
	ErrInvalidDayOfWeek     = errors.New("day of week must be between 1-7")
	ErrInvalidPrice         = errors.New("price must be postive")
	ErrScheduleHasBookings  = errors.New("cannot modifty schedule with existing bookings")
	ErrInvalidDuration      = errors.New("duration must be at least 1 hour")
	ErrInvalidTimeRange     = errors.New("invalid time range")
	ErrDuplicateFieldName   = errors.New("field with this name already exists in venue")
	ErrFieldNotFound        = errors.New("field not found")
	ErrVenueNotFound        = errors.New("venue not found")
	ErrInvalidFieldData     = errors.New("invalid field data")
	ErrFieldTypeNotFound    = errors.New("field type not found")
	ErrUserNotFound         = errors.New("user not found")
)
