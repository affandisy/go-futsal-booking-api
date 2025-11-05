package model

import (
	"go-futsal-booking-api/internal/domain"
	"time"

	"gorm.io/gorm"
)

type UserGorm struct {
	ID         uint   `gorm:"primaryKey"`
	RoleID     uint   `gorm:"column:role_id;not null"`
	FullName   string `gorm:"column:full_name;not null"`
	Email      string `gorm:"column:email;unique;not null"`
	IsVerified bool   `gorm:"column:is_verified;default:false"`
	Password   string `gorm:"column:password;not null"`
	Age        int    `gorm:"column:age;not null"`
	Address    string `gorm:"column:address;not null"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index"`

	Role RoleGorm `gorm:"foreignKey:RoleID"`
}

func (UserGorm) TableName() string {
	return "users"
}

func (ug *UserGorm) ToDomain() domain.User {
	var deletedAt *time.Time
	if ug.DeletedAt.Valid {
		deletedAt = &ug.DeletedAt.Time
	}

	return domain.User{
		ID:         ug.ID,
		FullName:   ug.FullName,
		Email:      ug.Email,
		IsVerified: ug.IsVerified,
		Password:   ug.Password,
		Age:        ug.Age,
		Address:    ug.Address,
		CreatedAt:  ug.CreatedAt,
		UpdatedAt:  ug.UpdatedAt,
		DeletedAt:  deletedAt,
		Role: domain.Role{
			ID:       ug.Role.ID,
			RoleName: ug.Role.RoleName,
		},
	}
}

func (ug *UserGorm) FromDomain(user domain.User) {
	ug.ID = user.ID
	ug.FullName = user.FullName
	ug.Email = user.Email
	ug.Password = user.Password
	ug.Age = user.Age
	ug.Address = user.Address
	ug.CreatedAt = user.CreatedAt
	ug.UpdatedAt = user.UpdatedAt
	ug.RoleID = user.Role.ID
}
