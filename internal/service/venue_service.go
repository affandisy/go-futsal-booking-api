package service

import (
	"context"
	"errors"
	"fmt"
	"go-futsal-booking-api/internal/domain"
	"go-futsal-booking-api/internal/repository"
	"go-futsal-booking-api/pkg/logger"
)

type VenueService interface {
	GetVenueByID(ctx context.Context, id uint) (*domain.Venue, error)
	GetAllVenues(ctx context.Context) ([]domain.Venue, error)
	CreateVenue(ctx context.Context, name, address, city string) (*domain.Venue, error)
	UpdateVenue(ctx context.Context, id uint, name, address, city string) (*domain.Venue, error)
	DeleteVenue(ctx context.Context, id uint) error
}

type venueService struct {
	venueRepo repository.VenueRepository
}

func NewVenueService(repo repository.VenueRepository) VenueService {
	return &venueService{
		venueRepo: repo,
	}
}

func (s *venueService) GetVenueByID(ctx context.Context, id uint) (*domain.Venue, error) {
	if id == 0 {
		logger.Error("Invalid venue id")
		return nil, errors.New("invalid venue id")
	}

	if err := ctx.Err(); err != nil {
		logger.Error("context error when get venue by id")
		return nil, fmt.Errorf("context error: %w", err)
	}

	venue, err := s.venueRepo.FindByID(ctx, id)
	if err != nil {
		logger.Error("failed to find venue by id", err.Error())
		return nil, err
	}

	return &venue, nil
}

func (s *venueService) GetAllVenues(ctx context.Context) ([]domain.Venue, error) {
	if err := ctx.Err(); err != nil {
		logger.Error("context error when get venue by venue")
		return nil, fmt.Errorf("context error: %w", err)
	}

	venues, err := s.venueRepo.FindAll(ctx)
	if err != nil {
		logger.Error("Failed to find venue by venue id", err)
		return nil, err
	}

	return venues, nil
}

func (s *venueService) CreateVenue(ctx context.Context, name, address, city string) (*domain.Venue, error) {
	if name == "" || address == "" || city == "" {
		logger.Error("Invalid venue data")
		return nil, errors.New("invalid venue data")
	}

	if err := ctx.Err(); err != nil {
		logger.Error("context error when create venue")
		return nil, fmt.Errorf("context error: %w", err)
	}

	newVenue := &domain.Venue{
		Name:    name,
		Address: address,
		City:    city,
	}

	if err := s.venueRepo.Create(ctx, newVenue); err != nil {
		logger.Error("failed to create new venue", err)
		return nil, fmt.Errorf("failed to create venue: %w", err)
	}

	logger.Info("venue created successfully")

	return newVenue, nil
}

func (s *venueService) UpdateVenue(ctx context.Context, id uint, name, address, city string) (*domain.Venue, error) {
	if id == 0 || name == "" || address == "" || city == "" {
		logger.Error("Invalid venue data")
		return nil, errors.New("invalid venue data")
	}

	if err := ctx.Err(); err != nil {
		logger.Error("context error when updating venue")
		return nil, fmt.Errorf("context error: %w", err)
	}

	venueUpdate, err := s.venueRepo.FindByID(ctx, id)
	if err != nil {
		logger.Error("venue not found", err)
		return nil, errors.New("venue not found")
	}

	venueUpdate.Name = name
	venueUpdate.Address = address
	venueUpdate.City = city

	if err := s.venueRepo.Update(ctx, &venueUpdate); err != nil {
		logger.Error("failed to update venue", err)
		return nil, fmt.Errorf("failed to update venue: %w", err)
	}

	logger.Info("venue updated success")

	return &venueUpdate, nil
}

func (s *venueService) DeleteVenue(ctx context.Context, id uint) error {
	if id == 0 {
		logger.Error("Invalid venue id when deleting venue")
		return errors.New("invalid venue id")
	}

	if err := ctx.Err(); err != nil {
		logger.Error("context error when deleting venue")
		return fmt.Errorf("context error: %w", err)
	}

	_, err := s.venueRepo.FindByID(ctx, id)
	if err != nil {
		logger.Error("venue not found", err)
		return errors.New("venue not found")
	}

	if err := s.venueRepo.Delete(ctx, id); err != nil {
		logger.Error("failed to delete venue", err)
		return fmt.Errorf("failed to delete venue: %w", err)
	}

	logger.Info("venue deleted success")

	return nil
}
