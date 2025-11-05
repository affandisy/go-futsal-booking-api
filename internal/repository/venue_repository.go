package repository

import (
	"context"
	"errors"
	"go-futsal-booking-api/internal/domain"
	gormContract "go-futsal-booking-api/internal/repository/model"

	"gorm.io/gorm"
)

type VenueRepository interface {
	Create(ctx context.Context, venue *domain.Venue) error
	FindByID(ctx context.Context, id uint) (domain.Venue, error)
	FindAll(ctx context.Context) ([]domain.Venue, error)
	Update(ctx context.Context, venue *domain.Venue) error
	Delete(ctx context.Context, id uint) error
}

type gormVenueRepository struct {
	DB *gorm.DB
}

func NewVenueRepository(db *gorm.DB) VenueRepository {
	return &gormVenueRepository{
		DB: db,
	}
}

func (r *gormVenueRepository) Create(ctx context.Context, venue *domain.Venue) error {
	var gormVenue gormContract.VenueGorm
	gormVenue.FromDomain(*venue)

	if err := r.DB.WithContext(ctx).Create(&gormVenue).Error; err != nil {
		return err
	}

	*venue = gormVenue.ToDomain()

	return nil
}

func (r *gormVenueRepository) FindByID(ctx context.Context, id uint) (domain.Venue, error) {
	var gormVenue gormContract.VenueGorm
	err := r.DB.WithContext(ctx).First(&gormVenue, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.Venue{}, err
		}
		return domain.Venue{}, err
	}
	return gormVenue.ToDomain(), nil
}

func (r *gormVenueRepository) FindAll(ctx context.Context) ([]domain.Venue, error) {
	var gormVenues []gormContract.VenueGorm
	if err := r.DB.WithContext(ctx).Find(&gormVenues).Error; err != nil {
		return nil, err
	}

	var domainVenues []domain.Venue
	for _, v := range gormVenues {
		domainVenues = append(domainVenues, v.ToDomain())
	}

	return domainVenues, nil
}

func (r *gormVenueRepository) Update(ctx context.Context, venue *domain.Venue) error {
	var gormVenue gormContract.VenueGorm
	gormVenue.FromDomain(*venue)

	result := r.DB.WithContext(ctx).Model(&gormVenue).Updates(gormVenue)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("venue not found")
	}

	return r.DB.WithContext(ctx).First(&gormVenue, venue.ID).Error
}

func (r *gormVenueRepository) Delete(ctx context.Context, id uint) error {
	result := r.DB.WithContext(ctx).Delete(&gormContract.VenueGorm{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("venue not found")
	}

	return nil
}
