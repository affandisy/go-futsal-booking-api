package model

import (
	"go-futsal-booking-api/internal/domain"
	"time"

	"gorm.io/gorm"
)

type BookingGorm struct {
	ID          uint      `gorm:"primaryKey"`
	UserID      uint      `gorm:"column:user_id;not null"`
	ScheduleID  uint      `gorm:"column:schedule_id;not null"`
	BookingDate time.Time `gorm:"column:booking_date;type:date;not null"`
	Status      string    `gorm:"column:status; not null"`
	TotalPrice  float64   `gorm:"column:total_price;type:numeric(10,2);not null"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`

	User     UserGorm     `gorm:"foreignKey:UserID"`
	Schedule ScheduleGorm `gorm:"foreignKey:ScheduleID"`
}

func (BookingGorm) TableName() string {
	return "bookings"
}

func (bg *BookingGorm) ToDomain() domain.Booking {
	var deletedAt *time.Time
	if bg.DeletedAt.Valid {
		deletedAt = &bg.DeletedAt.Time
	}

	return domain.Booking{
		ID:          bg.ID,
		BookingDate: bg.BookingDate,
		Status:      bg.Status,
		TotalPrice:  bg.TotalPrice,
		CreatedAt:   bg.CreatedAt,
		UpdatedAt:   bg.UpdatedAt,
		DeletedAt:   deletedAt,
		User:        bg.User.ToDomain(),
		Schedule:    bg.Schedule.ToDomain(),
	}
}

func (bg *BookingGorm) FromDomain(b domain.Booking) {
	bg.ID = b.ID
	bg.UserID = b.User.ID
	bg.ScheduleID = b.Schedule.ID
	bg.BookingDate = b.BookingDate
	bg.Status = b.Status
	bg.TotalPrice = b.TotalPrice
}
