package model

import (
	"go-futsal-booking-api/internal/domain"
	"time"

	"gorm.io/gorm"
)

type FieldGorm struct {
	ID        uint   `gorm:"primaryKey"`
	VenueID   uint   `gorm:"column:venue_id;not null"`
	Name      string `gorm:"column:name;not null"`
	Type      string `gorm:"column:type;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	Venue VenueGorm `gorm:"foreignKey:VenueID"`
}

func (FieldGorm) TableName() string {
	return "fields"
}

func (fg *FieldGorm) ToDomain() domain.Field {
	var deletedAt *time.Time
	if fg.DeletedAt.Valid {
		deletedAt = &fg.DeletedAt.Time
	}

	return domain.Field{
		ID:        fg.ID,
		Name:      fg.Name,
		Type:      fg.Type,
		CreatedAt: fg.CreatedAt,
		UpdatedAt: fg.UpdatedAt,
		DeletedAt: deletedAt,
		Venue:     fg.Venue.ToDomain(),
	}
}

func (fg *FieldGorm) FromDomain(field domain.Field) {
	fg.ID = field.ID
	fg.Name = field.Name
	fg.Type = field.Type
	fg.VenueID = field.Venue.ID
}
