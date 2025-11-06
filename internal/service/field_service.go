package service

import (
	"context"
	"errors"
	"fmt"
	"go-futsal-booking-api/internal/domain"
	"go-futsal-booking-api/internal/repository"
	"go-futsal-booking-api/pkg/logger"
)

var (
	ErrFieldNotFound     = errors.New("field not found")
	ErrVenueNotFound     = errors.New("venue not found")
	ErrInvalidFieldData  = errors.New("invalid field data")
	ErrFieldTypeNotFound = errors.New("field type not found")
)

type FieldService interface {
	GetFieldByID(ctx context.Context, id uint) (*domain.Field, error)
	GetFieldsByVenue(ctx context.Context, venueID uint) ([]domain.Field, error)
	CreateField(ctx context.Context, venueID uint, name, fieldType string) (*domain.Field, error)
	UpdateField(ctx context.Context, id uint, name, fieldType string) (*domain.Field, error)
	DeleteField(ctx context.Context, id uint) error
}

type fieldService struct {
	fieldRepo repository.FieldRepository
	venueRepo repository.VenueRepository
}

func NewFieldService(fieldRepo repository.FieldRepository, venueRepo repository.VenueRepository) FieldService {
	return &fieldService{
		fieldRepo: fieldRepo,
		venueRepo: venueRepo,
	}
}

func (s *fieldService) GetFieldByID(ctx context.Context, id uint) (*domain.Field, error) {
	if id == 0 {
		logger.Error("Invalid field id")
		return nil, ErrInvalidFieldData
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

func (s *fieldService) GetFieldsByVenue(ctx context.Context, venueID uint) ([]domain.Field, error) {
	if venueID == 0 {
		logger.Error("Invalid venue id")
		return nil, ErrInvalidFieldData
	}

	if err := ctx.Err(); err != nil {
		logger.Error("context error when get field by venue")
		return nil, fmt.Errorf("context error: %w", err)
	}

	_, err := s.venueRepo.FindByID(ctx, venueID)
	if err != nil {
		logger.Error("failed to find venue by id", err)
		return nil, ErrVenueNotFound
	}

	fields, err := s.fieldRepo.FindByVenueID(ctx, venueID)
	if err != nil {
		logger.Error("Failed to find field by venue id", err)
		return nil, err
	}

	return fields, nil
}

func (s *fieldService) CreateField(ctx context.Context, venueID uint, name, fieldType string) (*domain.Field, error) {
	if venueID == 0 || fieldType == "" {
		logger.Error("Invalid venue id and field type")
		return nil, ErrInvalidFieldData
	}

	if err := ctx.Err(); err != nil {
		logger.Error("context error when create field")
		return nil, fmt.Errorf("context error: %w", err)
	}

	venue, err := s.venueRepo.FindByID(ctx, venueID)
	if err != nil {
		logger.Error("venue not found", err)
		return nil, ErrVenueNotFound
	}

	newField := &domain.Field{
		Name:  name,
		Type:  fieldType,
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
		return nil, ErrInvalidFieldData
	}

	if err := ctx.Err(); err != nil {
		logger.Error("context error when updating field")
		return nil, fmt.Errorf("context error: %w", err)
	}

	fieldUpdate, err := s.fieldRepo.FindByID(ctx, id)
	if err != nil {
		logger.Error("field not found", err)
		return nil, ErrFieldNotFound
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
		return ErrInvalidFieldData
	}

	if err := ctx.Err(); err != nil {
		logger.Error("context error when deleting field")
		return fmt.Errorf("context error: %w", err)
	}

	_, err := s.fieldRepo.FindByID(ctx, id)
	if err != nil {
		logger.Error("field not found", err)
		return ErrFieldNotFound
	}

	if err := s.fieldRepo.Delete(ctx, id); err != nil {
		logger.Error("failed to delete field", err)
		return fmt.Errorf("failed to delete field: %w", err)
	}

	logger.Info("field deleted success")

	return nil
}
