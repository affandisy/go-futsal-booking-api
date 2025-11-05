package gorm

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

func (userGorm *UserGorm) ToDomain() domain.User {
	var deletedAt *time.Time
	if userGorm.DeletedAt.Valid {
		deletedAt = &userGorm.DeletedAt.Time
	}

	return domain.User{
		ID:         userGorm.ID,
		FullName:   userGorm.FullName,
		Email:      userGorm.Email,
		IsVerified: userGorm.IsVerified,
		Password:   userGorm.Password,
		Age:        userGorm.Age,
		Address:    userGorm.Address,
		CreatedAt:  userGorm.CreatedAt,
		UpdatedAt:  userGorm.UpdatedAt,
		DeletedAt:  deletedAt,
		Role: domain.Role{
			ID:       userGorm.Role.ID,
			RoleName: userGorm.Role.RoleName,
		},
	}
}

func (userGorm *UserGorm) FromDomain(user domain.User) {
	userGorm.ID = user.ID
	userGorm.FullName = user.FullName
	userGorm.Email = user.Email
	userGorm.Password = user.Password
	userGorm.Age = user.Age
	userGorm.Address = user.Address
	userGorm.CreatedAt = user.CreatedAt
	userGorm.UpdatedAt = user.UpdatedAt
	userGorm.RoleID = user.Role.ID
}
