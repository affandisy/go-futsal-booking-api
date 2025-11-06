package repository

import (
	"context"
	"errors"
	"fmt"
	"go-futsal-booking-api/internal/domain"
	gormContract "go-futsal-booking-api/internal/repository/model"
	"time"

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
	if err := ctx.Err(); err != nil {
		return fmt.Errorf("context error: %w", err)
	}

	var gormVenue gormContract.VenueGorm
	gormVenue.FromDomain(*venue)

	now := time.Now()
	gormVenue.CreatedAt = now
	gormVenue.UpdatedAt = now

	if err := r.DB.WithContext(ctx).Create(&gormVenue).Error; err != nil {
		return fmt.Errorf("failed to create venue: %w", err)
	}

	*venue = gormVenue.ToDomain()

	return nil
}

func (r *gormVenueRepository) FindByID(ctx context.Context, id uint) (domain.Venue, error) {
	if err := ctx.Err(); err != nil {
		return domain.Venue{}, fmt.Errorf("context error: %w", err)
	}

	var gormVenue gormContract.VenueGorm
	err := r.DB.WithContext(ctx).Where("deleted_at IS NULL").First(&gormVenue, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.Venue{}, errors.New("venue not found")
		}
		return domain.Venue{}, fmt.Errorf("failed to find field: %w", err)
	}

	venue := gormVenue.ToDomain()

	return venue, nil
}

func (r *gormVenueRepository) FindAll(ctx context.Context) ([]domain.Venue, error) {
	if err := ctx.Err(); err != nil {
		return nil, fmt.Errorf("context error: %w", err)
	}

	var gormVenues []gormContract.VenueGorm
	err := r.DB.WithContext(ctx).Where("deleted_at IS NULL").Find(&gormVenues).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find venues: %w", err)
	}

	var domainVenues []domain.Venue
	for _, v := range gormVenues {
		domainVenues = append(domainVenues, v.ToDomain())
	}

	return domainVenues, nil
}

func (r *gormVenueRepository) Update(ctx context.Context, venue *domain.Venue) error {
	if err := ctx.Err(); err != nil {
		return fmt.Errorf("context error: %w", err)
	}

	var gormVenue gormContract.VenueGorm
	gormVenue.FromDomain(*venue)

	updateVenue := map[string]interface{}{
		"name":       gormVenue.Name,
		"address":    gormVenue.Address,
		"city":       gormVenue.City,
		"updated_at": time.Now(),
	}

	result := r.DB.WithContext(ctx).Model(&gormVenue).Where("id = ? AND deleted_at IS NULL", venue.ID).Updates(updateVenue)
	if result.Error != nil {
		return fmt.Errorf("failed to update venue: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return errors.New("venue not found or already deleted")
	}

	*venue = gormVenue.ToDomain()

	return nil
}

func (r *gormVenueRepository) Delete(ctx context.Context, id uint) error {
	if err := ctx.Err(); err != nil {
		return fmt.Errorf("context error: %w", err)
	}

	result := r.DB.WithContext(ctx).Model(&gormContract.VenueGorm{}).Where("id = ? AND deleted_at IS NULL", id).Update("deleted_at", time.Now())
	if result.Error != nil {
		return fmt.Errorf("failed to delete field: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return errors.New("venue not found or already deleted")
	}

	return nil
}
