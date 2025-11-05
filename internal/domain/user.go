package domain

import (
	"time"
)

type User struct {
	ID         uint
	FullName   string
	Email      string
	IsVerified bool
	Password   string
	Age        int
	Address    string
	Role       Role
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  *time.Time
}
