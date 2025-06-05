package services

import (
	"context"
	"errors"
	"inventory-system/internal/models"
	"inventory-system/internal/repository"
)

var (
	ErrStatusNotFound = errors.New("status not found")
)

type AssetService struct {
	assetRepo      *repository.AssetRepository
	statusRepo     *repository.StatusRepository
	locationRepo   *repository.LocationRepository
	departmentRepo *repository.DepartmentRepository
}

func NewAssetService(
	assetRepo *repository.AssetRepository,
	statusRepo *repository.StatusRepository,
	locationRepo *repository.LocationRepository,
	departmentRepo *repository.DepartmentRepository,
) *AssetService {
	return &AssetService{
		assetRepo:      assetRepo,
		statusRepo:     statusRepo,
		locationRepo:   locationRepo,
		departmentRepo: departmentRepo,
	}
}

func (s *AssetService) GetAllAssets(ctx context.Context) ([]models.Asset, error) {
	return s.assetRepo.GetAll(ctx)
}

func (s *AssetService) GetAssetByID(ctx context.Context, id int) (*models.Asset, error) {
	asset, err := s.assetRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if asset == nil {
		return nil, ErrAssetNotFound
	}
	return asset, nil
}

func (s *AssetService) CreateAsset(ctx context.Context, asset models.Asset) (int, error) {
	// Validate status
	exists, err := s.statusRepo.Exists(ctx, asset.StatusID)
	if err != nil {
		return 0, err
	}
	if !exists {
		return 0, ErrStatusNotFound
	}

	// Validate location
	exists, err = s.locationRepo.Exists(ctx, asset.CurrentLocationID)
	if err != nil {
		return 0, err
	}
	if !exists {
		return 0, ErrLocationNotFound
	}

	// Validate department if specified
	if asset.DepartmentID != nil {
		exists, err = s.departmentRepo.Exists(ctx, *asset.DepartmentID)
		if err != nil {
			return 0, err
		}
		if !exists {
			return 0, ErrDepartmentNotFound
		}
	}

	return s.assetRepo.Create(ctx, asset)
}

func (s *AssetService) UpdateAsset(ctx context.Context, asset models.Asset) error {
	// Validate status
	exists, err := s.statusRepo.Exists(ctx, asset.StatusID)
	if err != nil {
		return err
	}
	if !exists {
		return ErrStatusNotFound
	}

	// Validate location
	exists, err = s.locationRepo.Exists(ctx, asset.CurrentLocationID)
	if err != nil {
		return err
	}
	if !exists {
		return ErrLocationNotFound
	}

	// Validate department if specified
	if asset.DepartmentID != nil {
		exists, err = s.departmentRepo.Exists(ctx, *asset.DepartmentID)
		if err != nil {
			return err
		}
		if !exists {
			return ErrDepartmentNotFound
		}
	}

	return s.assetRepo.Update(ctx, asset)
}

func (s *AssetService) DeleteAsset(ctx context.Context, id int) error {
	// Check if asset exists
	exists, err := s.assetRepo.Exists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return ErrAssetNotFound
	}

	return s.assetRepo.Delete(ctx, id)
}

func (s *AssetService) UpdateStatus(ctx context.Context, assetID, statusID int) error {
	// Validate asset
	exists, err := s.assetRepo.Exists(ctx, assetID)
	if err != nil {
		return err
	}
	if !exists {
		return ErrAssetNotFound
	}

	// Validate status
	exists, err = s.statusRepo.Exists(ctx, statusID)
	if err != nil {
		return err
	}
	if !exists {
		return ErrStatusNotFound
	}

	return s.assetRepo.UpdateStatus(ctx, assetID, statusID)
}

func (s *AssetService) UpdateLocation(ctx context.Context, assetID, locationID int) error {
	// Validate asset
	exists, err := s.assetRepo.Exists(ctx, assetID)
	if err != nil {
		return err
	}
	if !exists {
		return ErrAssetNotFound
	}

	// Validate location
	exists, err = s.locationRepo.Exists(ctx, locationID)
	if err != nil {
		return err
	}
	if !exists {
		return ErrLocationNotFound
	}

	return s.assetRepo.UpdateLocation(ctx, assetID, locationID)
}

func (s *AssetService) GetAssetsByStatus(ctx context.Context, statusID int) ([]models.Asset, error) {
	// Validate status
	exists, err := s.statusRepo.Exists(ctx, statusID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, ErrStatusNotFound
	}

	return s.assetRepo.GetByStatus(ctx, statusID)
}

func (s *AssetService) GetAssetsByLocation(ctx context.Context, locationID int) ([]models.Asset, error) {
	// Validate location
	exists, err := s.locationRepo.Exists(ctx, locationID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, ErrLocationNotFound
	}

	return s.assetRepo.GetByLocation(ctx, locationID)
}

func (s *AssetService) GetAssetsByDepartment(ctx context.Context, departmentID int) ([]models.Asset, error) {
	// Validate department
	exists, err := s.departmentRepo.Exists(ctx, departmentID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, ErrDepartmentNotFound
	}

	return s.assetRepo.GetByDepartment(ctx, departmentID)
}

func (s *AssetService) GetAssetTransferHistory(ctx context.Context, assetID int) ([]models.AssetTransfer, error) {
	// Validate asset
	exists, err := s.assetRepo.Exists(ctx, assetID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, ErrAssetNotFound
	}

	return s.assetRepo.GetTransferHistory(ctx, assetID)
}

func (s *AssetService) GetAllStatuses(ctx context.Context) ([]models.AssetStatus, error) {
	return s.statusRepo.GetAll(ctx)
}
