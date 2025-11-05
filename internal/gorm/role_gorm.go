package gorm

import (
	"time"

	"gorm.io/gorm"
)

type RoleGorm struct {
	ID        uint   `gorm:"primaryKey"`
	RoleName  string `gorm:"column:role_name;unique;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (RoleGorm) TableName() string {
	return "roles"
}
