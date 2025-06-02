package services

import (
	"context"
	"errors"
	"inventory-system/internal/models"
	"inventory-system/internal/repository"
	"time"
)

var (
	ErrTransferNotFound         = errors.New("transfer not found")
	ErrAssetNotFound            = errors.New("asset not found")
	ErrTransferEmployeeNotFound = errors.New("employee not found")
	ErrLocationNotFound         = errors.New("location not found")
	ErrSameLocations            = errors.New("source and destination locations are the same")
	ErrInvalidTransferDate      = errors.New("transfer date cannot be in the future")
)

type TransferService struct {
	transferRepo *repository.TransferRepository
	assetRepo    *repository.AssetRepository
	employeeRepo *repository.EmployeeRepository
	locationRepo *repository.LocationRepository
}

func NewTransferService(
	transferRepo *repository.TransferRepository,
	assetRepo *repository.AssetRepository,
	employeeRepo *repository.EmployeeRepository,
	locationRepo *repository.LocationRepository,
) *TransferService {
	return &TransferService{
		transferRepo: transferRepo,
		assetRepo:    assetRepo,
		employeeRepo: employeeRepo,
		locationRepo: locationRepo,
	}
}

func (s *TransferService) GetAllTransfers(ctx context.Context) ([]models.AssetTransfer, error) {
	return s.transferRepo.GetAll(ctx)
}

func (s *TransferService) GetTransferByID(ctx context.Context, id int) (*models.AssetTransfer, error) {
	transfer, err := s.transferRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if transfer == nil {
		return nil, ErrTransferNotFound
	}
	return transfer, nil
}

func (s *TransferService) CreateTransfer(ctx context.Context, transfer models.AssetTransfer) (int, error) {
	// Validate asset
	asset, err := s.assetRepo.GetByID(ctx, transfer.AssetID)
	if err != nil {
		return 0, err
	}
	if asset == nil {
		return 0, ErrAssetNotFound
	}

	// Validate employee
	exists, err := s.employeeRepo.Exists(ctx, transfer.EmployeeID)
	if err != nil {
		return 0, err
	}
	if !exists {
		return 0, ErrEmployeeNotFound
	}

	// Validate from location (should match asset's current location)
	if asset.CurrentLocationID != transfer.FromLocationID {
		return 0, errors.New("source location does not match asset's current location")
	}

	// Validate to location
	exists, err = s.locationRepo.Exists(ctx, transfer.ToLocationID)
	if err != nil {
		return 0, err
	}
	if !exists {
		return 0, ErrLocationNotFound
	}

	// Check if locations are different
	if transfer.FromLocationID == transfer.ToLocationID {
		return 0, ErrSameLocations
	}

	// Validate transfer date
	if transfer.TransferDate.After(time.Now()) {
		return 0, ErrInvalidTransferDate
	}

	// Create transfer
	transferID, err := s.transferRepo.Create(ctx, transfer)
	if err != nil {
		return 0, err
	}

	// Update asset's current location
	err = s.assetRepo.UpdateLocation(ctx, transfer.AssetID, transfer.ToLocationID)
	if err != nil {
		return 0, err
	}

	return transferID, nil
}

func (s *TransferService) GetTransfersByAsset(ctx context.Context, assetID int) ([]models.AssetTransfer, error) {
	// Validate asset
	exists, err := s.assetRepo.Exists(ctx, assetID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, ErrAssetNotFound
	}

	return s.transferRepo.GetByAsset(ctx, assetID)
}

func (s *TransferService) GetTransfersByEmployee(ctx context.Context, employeeID int) ([]models.AssetTransfer, error) {
	// Validate employee
	exists, err := s.employeeRepo.Exists(ctx, employeeID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, ErrEmployeeNotFound
	}

	return s.transferRepo.GetByEmployee(ctx, employeeID)
}

func (s *TransferService) GetTransfersByLocation(ctx context.Context, locationID int) ([]models.AssetTransfer, error) {
	// Validate location
	exists, err := s.locationRepo.Exists(ctx, locationID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, ErrLocationNotFound
	}

	return s.transferRepo.GetByLocation(ctx, locationID)
}

func (s *TransferService) GetTransfersByDateRange(ctx context.Context, from, to time.Time) ([]models.AssetTransfer, error) {
	if from.After(to) {
		return nil, errors.New("start date cannot be after end date")
	}

	return s.transferRepo.GetByDateRange(ctx, from, to)
}
