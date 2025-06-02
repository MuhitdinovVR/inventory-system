package handlers

import (
	"inventory-system/internal/services"
)

type Handler struct {
	departmentService *services.DepartmentService
	employeeService   *services.EmployeeService
	assetService      *services.AssetService
	locationService   *services.LocationService
	transferService   *services.TransferService
	authService       *services.AuthService
	reportService     *services.ReportService
}

func NewHandler(
	departmentService *services.DepartmentService,
	employeeService *services.EmployeeService,
	assetService *services.AssetService,
	locationService *services.LocationService,
	transferService *services.TransferService,
	authService *services.AuthService,
	reportService *services.ReportService,
) *Handler {
	return &Handler{
		departmentService: departmentService,
		employeeService:   employeeService,
		assetService:      assetService,
		locationService:   locationService,
		transferService:   transferService,
		authService:       authService,
		reportService:     reportService,
	}
}
