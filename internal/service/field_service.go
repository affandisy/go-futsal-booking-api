package service

import (
	"context"
	"go-futsal-booking-api/internal/domain"
	"go-futsal-booking-api/internal/repository"
)

type FieldService interface {
	GetFieldByID(ctx context.Context, id uint) (domain.Field, error)
	GetFieldsByVenue(ctx context.Context, venueID uint) ([]domain.Field, error)
	CreateField(ctx context.Context, venueID uint, name, fieldType string) (domain.Field, error)
	UpdateField(ctx context.Context, id uint, name, fieldType string) (domain.Field, error)
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

func (s *fieldService) GetFieldByID(ctx context.Context, id uint) (domain.Field, error) {
	return s.fieldRepo.FindByID(ctx, id)
}

func (s *fieldService) GetFieldsByVenue(ctx context.Context, venueID uint) ([]domain.Field, error) {
	_, err := s.venueRepo.FindByID(ctx, venueID)
	if err != nil {
		return nil, err
	}
	return s.fieldRepo.FindByVenueID(ctx, venueID)
}

func (s *fieldService) CreateField(ctx context.Context, venueID uint, name, fieldType string) (domain.Field, error) {
	venue, err := s.venueRepo.FindByID(ctx, venueID)
	if err != nil {
		return domain.Field{}, err
	}

	newField := domain.Field{
		Name:  name,
		Type:  fieldType,
		Venue: venue,
	}

	if err := s.fieldRepo.Create(ctx, &newField); err != nil {
		return domain.Field{}, err
	}

	return newField, nil
}

func (s *fieldService) UpdateField(ctx context.Context, id uint, name, fieldType string) (domain.Field, error) {
	fieldUpdate, err := s.fieldRepo.FindByID(ctx, id)
	if err != nil {
		return domain.Field{}, err
	}

	fieldUpdate.Name = name
	fieldUpdate.Type = fieldType

	if err := s.fieldRepo.Update(ctx, &fieldUpdate); err != nil {
		return domain.Field{}, err
	}

	return fieldUpdate, nil
}

func (s *fieldService) DeleteField(ctx context.Context, id uint) error {
	return s.fieldRepo.Delete(ctx, id)
}
