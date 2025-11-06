package repository

import (
	"context"
	"errors"
	"go-futsal-booking-api/internal/domain"
	gormContract "go-futsal-booking-api/internal/repository/model"

	"gorm.io/gorm"
)

type ScheduleRepository interface {
	Create(ctx context.Context, schedule *domain.Schedule) error
	FindByID(ctx context.Context, id uint) (domain.Schedule, error)
	FindByFieldID(ctx context.Context, fieldID uint) ([]domain.Schedule, error) // Metode kustom
	Update(ctx context.Context, schedule *domain.Schedule) error
	Delete(ctx context.Context, id uint) error
}

type gormScheduleRepository struct {
	DB *gorm.DB
}

func NewScheduleRepository(db *gorm.DB) ScheduleRepository {
	return &gormScheduleRepository{DB: db}
}

func (r *gormScheduleRepository) preload(ctx context.Context) *gorm.DB {
	return r.DB.WithContext(ctx).Preload("Field.Venue")
}

func (r *gormScheduleRepository) Create(ctx context.Context, schedule *domain.Schedule) error {
	var gormSchedule gormContract.ScheduleGorm
	gormSchedule.FromDomain(*schedule)

	if err := r.DB.WithContext(ctx).Create(&gormSchedule).Error; err != nil {
		return err
	}

	if err := r.preload(ctx).First(&gormSchedule, gormSchedule.ID).Error; err != nil {
		return err
	}

	*schedule = gormSchedule.ToDomain()

	return nil
}

func (r *gormScheduleRepository) FindByID(ctx context.Context, id uint) (domain.Schedule, error) {
	var gormSchedule gormContract.ScheduleGorm

	err := r.preload(ctx).First(&gormSchedule, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.Schedule{}, domain.ErrScheduleNotFound
		}
		return domain.Schedule{}, err
	}

	return gormSchedule.ToDomain(), nil
}

func (r *gormScheduleRepository) FindByFieldID(ctx context.Context, fieldID uint) ([]domain.Schedule, error) {
	var gormSchedules []gormContract.ScheduleGorm

	err := r.preload(ctx).Where("field_id = ?", fieldID).Find(&gormSchedules).Error
	if err != nil {
		return nil, err
	}

	var domainSchedules []domain.Schedule
	for _, s := range gormSchedules {
		domainSchedules = append(domainSchedules, s.ToDomain())
	}

	return domainSchedules, nil
}

func (r *gormScheduleRepository) Update(ctx context.Context, schedule *domain.Schedule) error {
	var gormSchedule gormContract.ScheduleGorm
	gormSchedule.FromDomain(*schedule)

	result := r.DB.WithContext(ctx).Model(&gormSchedule).Updates(gormSchedule)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return domain.ErrScheduleNotFound
	}

	if err := r.preload(ctx).First(&gormSchedule, schedule.ID).Error; err != nil {
		return err
	}

	*schedule = gormSchedule.ToDomain()

	return nil
}

func (r *gormScheduleRepository) Delete(ctx context.Context, id uint) error {
	result := r.DB.WithContext(ctx).Delete(&gormContract.ScheduleGorm{}, id)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return domain.ErrScheduleNotFound
	}

	return nil
}
