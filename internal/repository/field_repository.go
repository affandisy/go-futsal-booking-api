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

type FieldRepository interface {
	Create(ctx context.Context, field *domain.Field) error
	FindByID(ctx context.Context, id uint) (domain.Field, error)
	FindByVenueID(ctx context.Context, venueID uint) ([]domain.Field, error)
	Update(ctx context.Context, field *domain.Field) error
	Delete(ctx context.Context, id uint) error
}

type gormFieldRepository struct {
	DB *gorm.DB
}

func NewFieldRepository(db *gorm.DB) FieldRepository {
	return &gormFieldRepository{
		DB: db,
	}
}

func (r *gormFieldRepository) preloadVenue(ctx context.Context) *gorm.DB {
	return r.DB.WithContext(ctx).Preload("Venue")
}

func (r *gormFieldRepository) Create(ctx context.Context, field *domain.Field) error {
	if err := ctx.Err(); err != nil {
		return fmt.Errorf("context error: %w", err)
	}

	var gormField gormContract.FieldGorm
	gormField.FromDomain(*field)

	now := time.Now()
	gormField.CreatedAt = now
	gormField.UpdatedAt = now

	if err := r.DB.WithContext(ctx).Create(&gormField).Error; err != nil {
		return fmt.Errorf("failed to create field: %w", err)
	}

	if err := r.preloadVenue(ctx).First(&gormField, gormField.ID).Error; err != nil {
		return fmt.Errorf("failed to preload venue: %w", err)
	}

	*field = gormField.ToDomain()
	return nil
}

func (r *gormFieldRepository) FindByID(ctx context.Context, id uint) (domain.Field, error) {
	if err := ctx.Err(); err != nil {
		return domain.Field{}, fmt.Errorf("context error: %w", err)
	}

	var gormField gormContract.FieldGorm
	err := r.preloadVenue(ctx).Where("deleted_at IS NULL").First(&gormField, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.Field{}, errors.New("field not found")
		}
		return domain.Field{}, fmt.Errorf("failed to find field: %w", err)
	}

	field := gormField.ToDomain()

	return field, nil
}

func (r *gormFieldRepository) FindByVenueID(ctx context.Context, venueID uint) ([]domain.Field, error) {
	if err := ctx.Err(); err != nil {
		return nil, fmt.Errorf("context error: %w", err)
	}

	var gormFields []gormContract.FieldGorm
	err := r.preloadVenue(ctx).Where("venue_id = ? AND deleted_at IS NULL", venueID).Find(&gormFields).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find fields: %w", err)
	}

	var domainFields []domain.Field
	for _, f := range gormFields {
		domainFields = append(domainFields, f.ToDomain())
	}
	return domainFields, nil
}

func (r *gormFieldRepository) Update(ctx context.Context, field *domain.Field) error {
	if err := ctx.Err(); err != nil {
		return fmt.Errorf("context error: %w", err)
	}

	var gormField gormContract.FieldGorm
	gormField.FromDomain(*field)

	updates := map[string]interface{}{
		"name":       gormField.Name,
		"type":       gormField.Type,
		"updated_at": time.Now(),
	}

	result := r.DB.WithContext(ctx).Model(&gormField).Where("id = ? AND deleted_at IS NULL", field.ID).Updates(updates)
	if result.Error != nil {
		return fmt.Errorf("failed to update field: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return errors.New("field not found or already deleted")
	}

	if err := r.preloadVenue(ctx).First(&gormField, field.ID).Error; err != nil {
		return fmt.Errorf("failed to update field: %w", err)
	}

	*field = gormField.ToDomain()

	return nil
}

func (r *gormFieldRepository) Delete(ctx context.Context, id uint) error {
	if err := ctx.Err(); err != nil {
		return fmt.Errorf("context error: %w", err)
	}

	result := r.DB.WithContext(ctx).Model(&gormContract.FieldGorm{}).Where("id = ? AND deleted_at IS NULL", id).Update("deleted_at", time.Now())
	if result.Error != nil {
		return fmt.Errorf("failed to delete field: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return errors.New("field not found or already deleted")
	}

	return nil
}
