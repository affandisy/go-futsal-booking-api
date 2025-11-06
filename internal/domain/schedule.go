package domain

import "time"

type Schedule struct {
	ID        uint
	Field     Field
	DayOfWeek int
	StartTime time.Time
	EndTime   time.Time
	Price     float64
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}
