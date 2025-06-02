package models

type Asset struct {
	ID                int     `json:"id"`
	Name              string  `json:"name"`
	Category          string  `json:"category"`
	AcquisitionDate   string  `json:"acquisition_date"`
	Cost              float64 `json:"cost"`
	StatusID          int     `json:"status_id"`
	Status            string  `json:"status"`
	CurrentLocationID int     `json:"current_location_id"`
	Location          string  `json:"location"`
	DepartmentID      *int    `json:"department_id,omitempty"`
	Department        *string `json:"department,omitempty"`
}
