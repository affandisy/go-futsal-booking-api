package model

import (
	"go-futsal-booking-api/internal/domain"
	"time"

	"gorm.io/gorm"
)

type ScheduleGorm struct {
	ID        uint      `gorm:"primaryKey"`
	FieldID   uint      `gorm:"column:field_id;not null"`
	DayOfWeek int       `gorm:"column:day_of_week;not null"`
	StartTime TimeOfDay `gorm:"column:start_time;type:time;not null"`
	EndTime   TimeOfDay `gorm:"column:end_time;type:time;not null"`
	Price     float64   `gorm:"column:price;type:numeric(10,2);not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	Field FieldGorm `gorm:"foreignKey:FieldID"`
}

func (ScheduleGorm) TableName() string {
	return "schedules"
}

func (sg *ScheduleGorm) ToDomain() domain.Schedule {
	var deletedAt *time.Time

	if sg.DeletedAt.Valid {
		deletedAt = &sg.DeletedAt.Time
	}

	return domain.Schedule{
		ID:        sg.ID,
		DayOfWeek: sg.DayOfWeek,
		StartTime: sg.StartTime.ToTime(),
		EndTime:   sg.EndTime.ToTime(),
		Price:     sg.Price,
		CreatedAt: sg.CreatedAt,
		UpdatedAt: sg.UpdatedAt,
		DeletedAt: deletedAt,
		Field:     sg.Field.ToDomain(),
	}
}

func (sg *ScheduleGorm) FromDomain(s domain.Schedule) {
	sg.ID = s.ID
	sg.FieldID = s.Field.ID
	sg.DayOfWeek = s.DayOfWeek
	sg.StartTime = NewTimeOfDay(s.StartTime)
	sg.EndTime = NewTimeOfDay(s.EndTime)
	sg.Price = s.Price
}
