package service

import (
	"context"
	"errors"
	"go-futsal-booking-api/internal/domain"
	"go-futsal-booking-api/internal/repository"
	"time"
)

type ScheduleService interface {
	GetScheduleByID(ctx context.Context, id uint) (domain.Schedule, error)
	GetScheduleByField(ctx context.Context, fieldID uint) ([]domain.Schedule, error)
	CreateSchedule(ctx context.Context, fieldId uint, dayOfWeek int, startTime, endTime string, price float64) (domain.Schedule, error)
	UpdateSchedule(ctx context.Context, id uint, dayOfWeek int, startTime, endTime string, price float64) (domain.Schedule, error)
	DeleteSchedule(ctx context.Context, id uint) error
}

type scheduleService struct {
	scheduleRepo repository.ScheduleRepository
	fieldRepo    repository.FieldRepository
}

func NewScheduleService(scheduleRepo repository.ScheduleRepository, fieldRepo repository.FieldRepository) ScheduleService {
	return &scheduleService{scheduleRepo: scheduleRepo, fieldRepo: fieldRepo}
}

func parseTime(timeStr string) (time.Time, error) {
	return time.Parse("15:04", timeStr)
}

func (s *scheduleService) GetScheduleByID(ctx context.Context, id uint) (domain.Schedule, error) {
	return s.scheduleRepo.FindByID(ctx, id)
}

func (s *scheduleService) GetScheduleByField(ctx context.Context, fieldID uint) ([]domain.Schedule, error) {
	if _, err := s.fieldRepo.FindByID(ctx, fieldID); err != nil {
		return nil, err
	}

	return s.scheduleRepo.FindByFieldID(ctx, fieldID)
}

func (s *scheduleService) CreateSchedule(ctx context.Context, fieldId uint, dayOfWeek int, startTime, endTime string, price float64) (domain.Schedule, error) {
	startT, err := parseTime(startTime)
	if err != nil {
		return domain.Schedule{}, errors.New("invalid start time format")
	}
	endT, err := parseTime(endTime)
	if err != nil {
		return domain.Schedule{}, errors.New("invalid end time format")
	}

	if !endT.After(startT) {
		return domain.Schedule{}, errors.New("end time must be after start time")
	}

	field, err := s.fieldRepo.FindByID(ctx, fieldId)
	if err != nil {
		return domain.Schedule{}, err
	}

	newSchedule := domain.Schedule{
		Field:     field,
		DayOfWeek: dayOfWeek,
		StartTime: startT,
		EndTime:   endT,
		Price:     price,
	}

	if err := s.scheduleRepo.Create(ctx, &newSchedule); err != nil {
		return domain.Schedule{}, err
	}

	return newSchedule, nil
}

func (s *scheduleService) UpdateSchedule(ctx context.Context, id uint, dayOfWeek int, startTime, endTime string, price float64) (domain.Schedule, error) {
	scheduleUpdate, err := s.scheduleRepo.FindByID(ctx, id)
	if err != nil {
		return domain.Schedule{}, err
	}

	startT, err := parseTime(startTime)
	if err != nil {
		return domain.Schedule{}, errors.New("invalid start time format")
	}

	endT, err := parseTime(endTime)
	if err != nil {
		return domain.Schedule{}, errors.New("invalid end time format")
	}

	if !endT.After(startT) {
		return domain.Schedule{}, errors.New("end time must be after start time")
	}

	scheduleUpdate.DayOfWeek = dayOfWeek
	scheduleUpdate.StartTime = startT
	scheduleUpdate.EndTime = endT
	scheduleUpdate.Price = price

	if err := s.scheduleRepo.Update(ctx, &scheduleUpdate); err != nil {
		return domain.Schedule{}, err
	}

	return scheduleUpdate, nil
}

func (s *scheduleService) DeleteSchedule(ctx context.Context, id uint) error {
	if _, err := s.scheduleRepo.FindByID(ctx, id); err != nil {
		return err
	}

	err := s.scheduleRepo.Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
