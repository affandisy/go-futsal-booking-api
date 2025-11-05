package domain

import (
	"time"
)

type Venue struct {
	ID        uint
	Name      string
	Address   string
	City      string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}
