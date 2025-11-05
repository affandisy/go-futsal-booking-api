package service

import (
	"context"
	"go-futsal-booking-api/internal/domain"
	"go-futsal-booking-api/internal/repository"
)

type VenueService interface {
	GetVenueByID(ctx context.Context, id uint) (domain.Venue, error)
	GetAllVenues(ctx context.Context) ([]domain.Venue, error)
	CreateVenue(ctx context.Context, name, address, city string) (domain.Venue, error)
	UpdateVenue(ctx context.Context, id uint, name, address, city string) (domain.Venue, error)
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

func (s *venueService) GetVenueByID(ctx context.Context, id uint) (domain.Venue, error) {
	return s.venueRepo.FindByID(ctx, id)
}

func (s *venueService) GetAllVenues(ctx context.Context) ([]domain.Venue, error) {
	return s.venueRepo.FindAll(ctx)
}

func (s *venueService) isAdmin(actor domain.User) bool {
	return actor.Role.RoleName == "admin"
}

func (s *venueService) CreateVenue(ctx context.Context, name, address, city string) (domain.Venue, error) {
	newVenue := domain.Venue{
		Name:    name,
		Address: address,
		City:    city,
	}

	if err := s.venueRepo.Create(ctx, &newVenue); err != nil {
		return domain.Venue{}, err
	}

	return newVenue, nil
}

func (s *venueService) UpdateVenue(ctx context.Context, id uint, name, address, city string) (domain.Venue, error) {
	venueUpdate, err := s.venueRepo.FindByID(ctx, id)
	if err != nil {
		return domain.Venue{}, err
	}

	venueUpdate.Name = name
	venueUpdate.Address = address
	venueUpdate.City = city

	if err := s.venueRepo.Update(ctx, &venueUpdate); err != nil {
		return domain.Venue{}, err
	}

	return venueUpdate, nil
}

func (s *venueService) DeleteVenue(ctx context.Context, id uint) error {
	return s.venueRepo.Delete(ctx, id)
}
