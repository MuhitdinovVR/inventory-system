package models

import "time"

type AssetsByStatusReport struct {
	Status    string  `json:"status"`
	Count     int     `json:"count"`
	TotalCost float64 `json:"total_cost"`
}

type AssetsCostByDepartmentReport struct {
	Department  string  `json:"department"`
	Count       int     `json:"count"`
	TotalCost   float64 `json:"total_cost"`
	AverageCost float64 `json:"average_cost"`
}

type InventoryReport struct {
	TotalAssets     int             `json:"total_assets"`
	TotalValue      float64         `json:"total_value"`
	ByStatus        map[string]int  `json:"by_status"`
	ByLocation      map[string]int  `json:"by_location"`
	ByDepartment    map[string]int  `json:"by_department"`
	RecentTransfers []AssetTransfer `json:"recent_transfers"`
	GeneratedAt     time.Time       `json:"generated_at"`
}
