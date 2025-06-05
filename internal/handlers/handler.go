package handlers

import (
	"inventory-system/internal/services"
	"log"
	"net/http"
	"os"
	"strings"
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

// AuthMiddleware проверяет JWT токен в заголовке Authorization
func (h *Handler) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Получаем токен из заголовка
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			respondWithError(w, http.StatusUnauthorized, "Authorization header required")
			return
		}

		// Проверяем формат заголовка
		if !strings.HasPrefix(authHeader, "Bearer ") {
			respondWithError(w, http.StatusUnauthorized, "Invalid authorization format")
			return
		}

		// Извлекаем токен
		token := strings.TrimPrefix(authHeader, "Bearer ")

		// Проверяем токен
		_, err := h.authService.ValidateToken(r.Context(), token)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "Invalid token")
			return
		}

		next.ServeHTTP(w, r)
	})
}

// ServeIndex обрабатывает запрос к главной странице
func (h *Handler) ServeIndex(w http.ResponseWriter, r *http.Request) {
	if _, err := os.Stat("./frontend/static/index.html"); os.IsNotExist(err) {
		log.Printf("index.html not found: %v", err)
		http.NotFound(w, r)
		return
	}
	http.ServeFile(w, r, "./frontend/static/index.html")
}

// ServeLogin обрабатывает запрос к странице входа
func (h *Handler) ServeLogin(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./frontend/static/login.html")
}

// ServeRegister обрабатывает запрос к странице регистрации
func (h *Handler) ServeRegister(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./frontend/static/register.html")
}

// ServeAssets обрабатывает запрос к странице активов
func (h *Handler) ServeAssets(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./frontend/static/assets.html")
}

// ServeEmployees обрабатывает запрос к странице сотрудников
func (h *Handler) ServeEmployees(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./frontend/static/employees.html")
}

// ServeDepartments обрабатывает запрос к странице отделов
func (h *Handler) ServeDepartments(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./frontend/static/departments.html")
}

// ServeTransfers обрабатывает запрос к странице перемещений
func (h *Handler) ServeTransfers(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./frontend/static/transfers.html")
}

// ServeReports обрабатывает запрос к странице отчетов
func (h *Handler) ServeReports(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./frontend/static/reports.html")
}
