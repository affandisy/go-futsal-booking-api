package model

import (
	"go-futsal-booking-api/internal/domain"
	"time"

	"gorm.io/gorm"
)

type VenueGorm struct {
	ID        uint   `gorm:"primaryKey"`
	Name      string `gorm:"column:name;unique;not null"`
	Address   string `gorm:"column:address;not null"`
	City      string `gorm:"column:city;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (VenueGorm) TableName() string {
	return "venues"
}

func (vg *VenueGorm) ToDomain() domain.Venue {
	var deletedAt *time.Time
	if vg.DeletedAt.Valid {
		deletedAt = &vg.DeletedAt.Time
	}

	return domain.Venue{
		ID:        vg.ID,
		Name:      vg.Name,
		Address:   vg.Address,
		City:      vg.City,
		CreatedAt: vg.CreatedAt,
		UpdatedAt: vg.UpdatedAt,
		DeletedAt: deletedAt,
	}
}

func (vg *VenueGorm) FromDomain(venue domain.Venue) {
	vg.ID = venue.ID
	vg.Name = venue.Name
	vg.Address = venue.Address
	vg.City = venue.City
}
