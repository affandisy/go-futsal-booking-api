package service

import (
	"context"
	"fmt"
	"go-futsal-booking-api/internal/domain"
	"go-futsal-booking-api/internal/dto/request"
	"go-futsal-booking-api/internal/repository"
	"go-futsal-booking-api/pkg/logger"
)

type FieldService interface {
	GetFieldByID(ctx context.Context, id uint) (*domain.Field, error)
	GetFieldsByVenue(ctx context.Context, venueID uint) ([]*domain.Field, error)
	CreateField(ctx context.Context, req *request.CreateFieldRequest) (*domain.Field, error)
	UpdateField(ctx context.Context, id uint, name, fieldType string) (*domain.Field, error)
	DeleteField(ctx context.Context, id uint) error
}

type fieldService struct {
	fieldRepo    repository.FieldRepository
	venueRepo    repository.VenueRepository
	scheduleRepo repository.ScheduleRepository
}

// type CreateFieldRequest struct {
// 	VenueID   uint
// 	FieldType string
// 	Name      string
// }

func NewFieldService(fieldRepo repository.FieldRepository, venueRepo repository.VenueRepository, scheduleRepo repository.ScheduleRepository) FieldService {
	return &fieldService{
		fieldRepo:    fieldRepo,
		venueRepo:    venueRepo,
		scheduleRepo: scheduleRepo,
	}
}

func (s *fieldService) GetFieldByID(ctx context.Context, id uint) (*domain.Field, error) {
	if id == 0 {
		logger.Error("Invalid field id")
		return nil, domain.ErrInvalidFieldData
	}

	if err := ctx.Err(); err != nil {
		logger.Error("context error when get field by id")
		return nil, fmt.Errorf("context error: %w", err)
	}

	field, err := s.fieldRepo.FindByID(ctx, id)
	if err != nil {
		logger.Error("failed to find field by id", err.Error())
		return nil, err
	}

	return &field, nil
}

func (s *fieldService) GetFieldsByVenue(ctx context.Context, venueID uint) ([]*domain.Field, error) {
	if venueID == 0 {
		logger.Error("Invalid venue id")
		return nil, domain.ErrInvalidFieldData
	}

	if err := ctx.Err(); err != nil {
		logger.Error("context error when get field by venue")
		return nil, fmt.Errorf("context error: %w", err)
	}

	_, err := s.venueRepo.FindByID(ctx, venueID)
	if err != nil {
		logger.Error("failed to find venue by id", err)
		return nil, domain.ErrVenueNotFound
	}

	fields, err := s.fieldRepo.FindByVenueID(ctx, venueID)
	if err != nil {
		logger.Error("Failed to find field by venue id", err)
		return nil, err
	}

	field := make([]*domain.Field, len(fields))
	for i := range fields {
		field[i] = &fields[i]
	}

	return field, nil
}

func (s *fieldService) CreateField(ctx context.Context, req *request.CreateFieldRequest) (*domain.Field, error) {
	if req.VenueID == 0 || req.FieldType == "" {
		logger.Error("Invalid venue id and field type")
		return nil, domain.ErrInvalidFieldData
	}

	if err := ctx.Err(); err != nil {
		logger.Error("context error when create field")
		return nil, fmt.Errorf("context error: %w", err)
	}

	venue, err := s.venueRepo.FindByID(ctx, req.VenueID)
	if err != nil {
		logger.Error("venue not found", err)
		return nil, domain.ErrVenueNotFound
	}

	newField := &domain.Field{
		Name:  req.Name,
		Type:  req.FieldType,
		Venue: venue,
	}

	if err := s.fieldRepo.Create(ctx, newField); err != nil {
		logger.Error("failed to create new field", err)
		return nil, fmt.Errorf("failed to create field: %w", err)
	}

	logger.Info("field created successfully")

	return newField, nil
}

func (s *fieldService) UpdateField(ctx context.Context, id uint, name, fieldType string) (*domain.Field, error) {
	if id == 0 || name == "" || fieldType == "" {
		logger.Error("Invalid field data")
		return nil, domain.ErrInvalidFieldData
	}

	if err := ctx.Err(); err != nil {
		logger.Error("context error when updating field")
		return nil, fmt.Errorf("context error: %w", err)
	}

	fieldUpdate, err := s.fieldRepo.FindByID(ctx, id)
	if err != nil {
		logger.Error("field not found", err)
		return nil, domain.ErrFieldNotFound
	}

	fieldUpdate.Name = name
	fieldUpdate.Type = fieldType

	if err := s.fieldRepo.Update(ctx, &fieldUpdate); err != nil {
		logger.Error("failed to update field", err)
		return nil, fmt.Errorf("failed to update field: %w", err)
	}

	logger.Info("field updated success")

	return &fieldUpdate, nil
}

func (s *fieldService) DeleteField(ctx context.Context, id uint) error {
	if id == 0 {
		logger.Error("Invalid field id when deleting field")
		return domain.ErrInvalidFieldData
	}

	if err := ctx.Err(); err != nil {
		logger.Error("context error when deleting field")
		return fmt.Errorf("context error: %w", err)
	}

	_, err := s.fieldRepo.FindByID(ctx, id)
	if err != nil {
		logger.Error("field not found", err)
		return domain.ErrFieldNotFound
	}

	if err := s.fieldRepo.Delete(ctx, id); err != nil {
		logger.Error("failed to delete field", err)
		return fmt.Errorf("failed to delete field: %w", err)
	}

	logger.Info("field deleted success")

	return nil
}
