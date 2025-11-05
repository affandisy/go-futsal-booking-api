package domain

import (
	"time"
)

type Field struct {
	ID        uint
	Name      string
	Type      string
	Venue     Venue
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}
