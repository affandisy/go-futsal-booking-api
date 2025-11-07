package service

import (
	"context"
	"errors"
	"fmt"
	"go-futsal-booking-api/internal/domain"
	"go-futsal-booking-api/internal/dto/request"
	"go-futsal-booking-api/internal/repository"
	"go-futsal-booking-api/pkg/logger"
	"time"
)

type ScheduleService interface {
	GetScheduleByID(ctx context.Context, id uint) (*domain.Schedule, error)
	GetScheduleByField(ctx context.Context, fieldID uint) ([]*domain.Schedule, error)
	CreateSchedule(ctx context.Context, req *request.CreateScheduleRequest) (*domain.Schedule, error)
	UpdateSchedule(ctx context.Context, id uint, dayOfWeek int, startTime, endTime string, price float64) (*domain.Schedule, error)
	DeleteSchedule(ctx context.Context, id uint) error
}

type scheduleService struct {
	scheduleRepo repository.ScheduleRepository
	fieldRepo    repository.FieldRepository
	bookingRepo  repository.BookingRepository
}

// type CreateScheduleRequest struct {
// 	FieldID   uint
// 	DayOfWeek int
// 	StartTime string
// 	EndTime   string
// 	Price     float64
// }

func NewScheduleService(scheduleRepo repository.ScheduleRepository, fieldRepo repository.FieldRepository, bookingRepo repository.BookingRepository) ScheduleService {
	return &scheduleService{scheduleRepo: scheduleRepo, fieldRepo: fieldRepo, bookingRepo: bookingRepo}
}

func parseTime(timeStr string) (time.Time, error) {
	return time.Parse("15:04", timeStr)
}

func (s *scheduleService) GetScheduleByID(ctx context.Context, id uint) (*domain.Schedule, error) {
	if id == 0 {
		return nil, errors.New("invalid schedule id")
	}

	if err := ctx.Err(); err != nil {
		return nil, fmt.Errorf("context error: %w", err)
	}

	schedule, err := s.scheduleRepo.FindByID(ctx, id)
	if err != nil {
		logger.Error("schedule not found", err.Error())
		return nil, domain.ErrScheduleNotFound
	}

	return &schedule, nil
}

func (s *scheduleService) GetScheduleByField(ctx context.Context, fieldID uint) ([]*domain.Schedule, error) {
	if fieldID == 0 {
		logger.Error("field id not found when get schedule by field")
		return nil, errors.New("invalid field id")
	}

	if err := ctx.Err(); err != nil {
		return nil, fmt.Errorf("context error: %w", err)
	}

	if _, err := s.fieldRepo.FindByID(ctx, fieldID); err != nil {
		logger.Error("field not found when get schedule by field id", err.Error())
		return nil, errors.New("field not found")
	}

	schedules, err := s.scheduleRepo.FindByFieldID(ctx, fieldID)
	if err != nil {
		logger.Error("failed to get schedule by field id", err.Error())
		return nil, err
	}

	result := make([]*domain.Schedule, len(schedules))
	for i := range schedules {
		result[i] = &schedules[i]
	}

	return result, nil
}

func (s *scheduleService) CreateSchedule(ctx context.Context, req *request.CreateScheduleRequest) (*domain.Schedule, error) {
	if req == nil || req.FieldID == 0 {
		logger.Error("missing request value to create schedule")
		return nil, errors.New("invalid schedule request")
	}

	if req.DayOfWeek < 1 || req.DayOfWeek > 7 {
		logger.Error("invalid day of week request")
		return nil, domain.ErrInvalidDayOfWeek
	}

	if req.Price <= 0 {
		logger.Error("invalid price when creating schedule")
		return nil, domain.ErrInvalidPrice
	}

	if err := ctx.Err(); err != nil {
		return nil, fmt.Errorf("context error: %w", err)
	}

	startT, err := parseTime(req.StartTime)
	if err != nil {
		return nil, errors.New("invalid start time format, use HH:MM")
	}
	endT, err := parseTime(req.EndTime)
	if err != nil {
		return nil, errors.New("invalid end time format, use HH:MM")
	}

	if !endT.After(startT) {
		return nil, errors.New("end time must be after start time")
	}

	field, err := s.fieldRepo.FindByID(ctx, req.FieldID)
	if err != nil {
		logger.Error("field not found when creating schedule", err.Error())
		return nil, errors.New("field not found")
	}

	newSchedule := domain.Schedule{
		Field:     field,
		DayOfWeek: req.DayOfWeek,
		StartTime: startT,
		EndTime:   endT,
		Price:     req.Price,
	}

	if err := s.scheduleRepo.Create(ctx, &newSchedule); err != nil {
		logger.Error("failed to create schedule", map[string]any{
			"field_id": req.FieldID,
			"error":    err.Error(),
		})
		return nil, fmt.Errorf("failed to create schedule: %w", err)
	}

	logger.Info("schedule created successfully")

	return &newSchedule, nil
}

func (s *scheduleService) UpdateSchedule(ctx context.Context, id uint, dayOfWeek int, startTime, endTime string, price float64) (*domain.Schedule, error) {
	if id == 0 {
		return nil, errors.New("invalid schedule id")
	}

	if dayOfWeek == 0 {
		return nil, errors.New("invalid schedule dayOfWeek")
	}

	if startTime == "" {
		return nil, errors.New("invalid schedule startTime")
	}

	if endTime == "" {
		return nil, errors.New("invalid schedule endTime")
	}

	if price == 0 {
		return nil, errors.New("invalid schedule price")
	}

	if err := ctx.Err(); err != nil {
		return nil, fmt.Errorf("context error: %w", err)
	}

	scheduleUpdate, err := s.scheduleRepo.FindByID(ctx, id)
	if err != nil {
		logger.Error("schedule not found when updating schedule", err.Error())
		return nil, domain.ErrScheduleNotFound
	}

	if dayOfWeek != 0 {
		if dayOfWeek < 1 || dayOfWeek > 7 {
			return nil, domain.ErrInvalidDayOfWeek
		}
		scheduleUpdate.DayOfWeek = dayOfWeek
	}

	if startTime != "" {
		startT, err := parseTime(startTime)
		if err != nil {
			return nil, errors.New("invalid start time format")
		}
		scheduleUpdate.StartTime = startT
	}

	if endTime != "" {
		endT, err := parseTime(endTime)
		if err != nil {
			return nil, errors.New("invalid end time format")
		}
		scheduleUpdate.EndTime = endT
	}

	if price != 0 {
		if price <= 0 {
			return nil, domain.ErrInvalidPrice
		}
		scheduleUpdate.Price = price
	}

	if err := s.scheduleRepo.Update(ctx, &scheduleUpdate); err != nil {
		logger.Error("failed to update schedule", map[string]any{
			"schedule_id": id,
			"error":       err.Error(),
		})
		return nil, fmt.Errorf("failed to update schedule: %w", err)
	}

	logger.Info("schedule update successfully")

	return &scheduleUpdate, nil
}

func (s *scheduleService) DeleteSchedule(ctx context.Context, id uint) error {
	if id == 0 {
		return errors.New("invalid schedule id")
	}

	if err := ctx.Err(); err != nil {
		return fmt.Errorf("context error: %w", err)
	}

	if _, err := s.scheduleRepo.FindByID(ctx, id); err != nil {
		return domain.ErrScheduleNotFound
	}

	err := s.scheduleRepo.Delete(ctx, id)
	if err != nil {
		logger.Error("failed to delete schedule", map[string]any{
			"schedule_id": id,
			"error":       err.Error(),
		})
		return fmt.Errorf("failed to delete schedule: %w", err)
	}

	return nil
}
