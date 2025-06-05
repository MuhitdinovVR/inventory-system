package models

import "time"

type AssetTransfer struct {
	ID             int       `json:"id"`
	AssetID        int       `json:"asset_id"`
	AssetName      string    `json:"asset_name"`
	EmployeeID     int       `json:"employee_id"`
	EmployeeName   string    `json:"employee_name"`
	FromLocationID int       `json:"from_location_id"`
	FromLocation   string    `json:"from_location"`
	ToLocationID   int       `json:"to_location_id"`
	ToLocation     string    `json:"to_location"`
	TransferDate   time.Time `json:"transfer_date"`
	Notes          string    `json:"notes,omitempty"`
}
