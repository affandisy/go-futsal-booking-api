package domain

import "time"

type Booking struct {
	ID          uint
	User        User
	Schedule    Schedule
	BookingDate time.Time
	Status      string
	TotalPrice  float64
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
}
