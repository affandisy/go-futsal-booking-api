package repository

import (
	"context"
	"errors"
	"go-futsal-booking-api/internal/domain"
	gormContract "go-futsal-booking-api/internal/repository/model"

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
	var gormField gormContract.FieldGorm
	gormField.FromDomain(*field)

	if err := r.DB.WithContext(ctx).Create(&gormField).Error; err != nil {
		return err
	}

	*field = gormField.ToDomain()
	return r.preloadVenue(ctx).First(&gormField, gormField.ID).Error
}

func (r *gormFieldRepository) FindByID(ctx context.Context, id uint) (domain.Field, error) {
	var gormField gormContract.FieldGorm
	err := r.preloadVenue(ctx).First(&gormField, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.Field{}, errors.New("field not found")
		}
		return domain.Field{}, err
	}
	return gormField.ToDomain(), nil
}

func (r *gormFieldRepository) FindByVenueID(ctx context.Context, venueID uint) ([]domain.Field, error) {
	var gormFields []gormContract.FieldGorm
	err := r.preloadVenue(ctx).Where("venue_id = ?", venueID).Find(&gormFields).Error
	if err != nil {
		return nil, err
	}

	var domainFields []domain.Field
	for _, f := range gormFields {
		domainFields = append(domainFields, f.ToDomain())
	}
	return domainFields, nil
}

func (r *gormFieldRepository) Update(ctx context.Context, field *domain.Field) error {
	var gormField gormContract.FieldGorm
	gormField.FromDomain(*field)

	result := r.DB.WithContext(ctx).Model(&gormField).Updates(gormField)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("field not found")
	}

	return r.preloadVenue(ctx).First(&gormField, field.ID).Error
}

func (r *gormFieldRepository) Delete(ctx context.Context, id uint) error {
	result := r.DB.WithContext(ctx).Delete(&gormContract.FieldGorm{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("field not found")
	}

	return nil
}
