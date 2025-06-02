package services

import (
	"context"
	"errors"
	"inventory-system/internal/models"
	"inventory-system/internal/repository"
)

var (
	ErrLocationHasAssets    = errors.New("cannot delete location with assigned assets")
	ErrLocationHasTransfers = errors.New("cannot delete location with transfer history")
)

type LocationService struct {
	locationRepo *repository.LocationRepository
}

func NewLocationService(locationRepo *repository.LocationRepository) *LocationService {
	return &LocationService{
		locationRepo: locationRepo,
	}
}

func (s *LocationService) GetAllLocations(ctx context.Context) ([]models.Location, error) {
	return s.locationRepo.GetAll(ctx)
}

func (s *LocationService) GetLocationByID(ctx context.Context, id int) (*models.Location, error) {
	location, err := s.locationRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if location == nil {
		return nil, ErrLocationNotFound
	}
	return location, nil
}

func (s *LocationService) CreateLocation(ctx context.Context, location models.Location) (int, error) {
	return s.locationRepo.Create(ctx, location)
}

func (s *LocationService) UpdateLocation(ctx context.Context, location models.Location) error {
	return s.locationRepo.Update(ctx, location)
}

func (s *LocationService) DeleteLocation(ctx context.Context, id int) error {
	// Check if location has assets
	hasAssets, err := s.locationRepo.HasAssets(ctx, id)
	if err != nil {
		return err
	}
	if hasAssets {
		return ErrLocationHasAssets
	}

	// Check if location has transfer history
	hasTransfers, err := s.locationRepo.HasTransfers(ctx, id)
	if err != nil {
		return err
	}
	if hasTransfers {
		return ErrLocationHasTransfers
	}

	return s.locationRepo.Delete(ctx, id)
}

func (s *LocationService) GetLocationsByType(ctx context.Context, locationType string) ([]models.Location, error) {
	return s.locationRepo.GetByType(ctx, locationType)
}
