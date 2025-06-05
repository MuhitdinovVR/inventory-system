package services

import (
	"context"
	"inventory-system/internal/models"
	"time"
)

type ReportService struct {
	assetService    *AssetService
	transferService *TransferService
}

func NewReportService(assetService *AssetService, transferService *TransferService) *ReportService {
	return &ReportService{
		assetService:    assetService,
		transferService: transferService,
	}
}

func (s *ReportService) GenerateAssetsByStatusReport(ctx context.Context) ([]models.AssetsByStatusReport, error) {
	statuses, err := s.assetService.GetAllStatuses(ctx)
	if err != nil {
		return nil, err
	}

	var report []models.AssetsByStatusReport
	for _, status := range statuses {
		assets, err := s.assetService.GetAssetsByStatus(ctx, status.ID)
		if err != nil {
			return nil, err
		}

		totalCost := 0.0
		for _, asset := range assets {
			totalCost += asset.Cost
		}

		report = append(report, models.AssetsByStatusReport{
			Status:    status.Name,
			Count:     len(assets),
			TotalCost: totalCost,
		})
	}

	return report, nil
}

func (s *ReportService) GenerateTransfersReport(ctx context.Context, from, to time.Time) ([]models.AssetTransfer, error) {
	return s.transferService.GetTransfersByDateRange(ctx, from, to)
}

func (s *ReportService) GenerateAssetsCostByDepartmentReport(ctx context.Context) ([]models.AssetsCostByDepartmentReport, error) {
	assets, err := s.assetService.GetAllAssets(ctx)
	if err != nil {
		return nil, err
	}

	departmentMap := make(map[int]*models.AssetsCostByDepartmentReport)

	for _, asset := range assets {
		var departmentID int
		var departmentName string

		if asset.DepartmentID != nil {
			departmentID = *asset.DepartmentID
			departmentName = *asset.Department
		} else {
			departmentID = 0
			departmentName = "Не назначен"
		}

		if _, ok := departmentMap[departmentID]; !ok {
			departmentMap[departmentID] = &models.AssetsCostByDepartmentReport{
				Department: departmentName,
				Count:      0,
				TotalCost:  0,
			}
		}

		departmentMap[departmentID].Count++
		departmentMap[departmentID].TotalCost += asset.Cost
	}

	var report []models.AssetsCostByDepartmentReport
	for _, item := range departmentMap {
		item.AverageCost = item.TotalCost / float64(item.Count)
		report = append(report, *item)
	}

	return report, nil
}

func (s *ReportService) GenerateInventoryReport(ctx context.Context) (*models.InventoryReport, error) {
	assets, err := s.assetService.GetAllAssets(ctx)
	if err != nil {
		return nil, err
	}
	
	report := &models.InventoryReport{
		TotalAssets:     len(assets),
		TotalValue:      0,
		ByStatus:        make(map[string]int),
		ByLocation:      make(map[string]int),
		ByDepartment:    make(map[string]int),
		RecentTransfers: make([]models.AssetTransfer, 0),
	}

	for _, asset := range assets {
		report.TotalValue += asset.Cost
		report.ByStatus[asset.Status]++
		report.ByLocation[asset.Location]++

		if asset.Department != nil {
			report.ByDepartment[*asset.Department]++
		} else {
			report.ByDepartment["Не назначен"]++
		}
	}

	// Get recent transfers (last 10)
	transfers, err := s.transferService.GetTransfersByDateRange(ctx, time.Now().AddDate(0, 0, -30), time.Now())
	if err != nil {
		return nil, err
	}

	if len(transfers) > 10 {
		report.RecentTransfers = transfers[:10]
	} else {
		report.RecentTransfers = transfers
	}

	return report, nil
}
