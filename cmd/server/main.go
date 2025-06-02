package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"inventory-system/internal/app"
	"inventory-system/internal/config"
	"inventory-system/internal/database"
	"inventory-system/internal/handlers"
	"inventory-system/internal/repository"
	"inventory-system/internal/services"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize database
	db, err := database.NewDatabase(cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Run database migrations
	if err := runMigrations(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Initialize repositories
	repos := initializeRepositories(db)

	// Initialize services
	services := initializeServices(repos, cfg)

	// Initialize handlers
	handler := handlers.NewHandler(
		services.DepartmentService,
		services.EmployeeService,
		services.AssetService,
		services.LocationService,
		services.TransferService,
		services.AuthService,
		services.ReportService,
	)

	// Setup router
	router := app.NewRouter(handler)

	// Serve static files
	fs := http.FileServer(http.Dir("./frontend/static"))
	router.Handle("/static/*", http.StripPrefix("/static/", fs))

	// Serve HTML pages
	router.HandleFunc("/", serveIndex)
	router.HandleFunc("/login", serveLogin)
	router.HandleFunc("/register", serveRegister)
	router.HandleFunc("/assets", serveAssets)
	router.HandleFunc("/employees", serveEmployees)
	router.HandleFunc("/departments", serveDepartments)
	router.HandleFunc("/transfers", serveTransfers)
	router.HandleFunc("/reports", serveReports)

	// Create HTTP server
	server := &http.Server{
		Addr:    cfg.Server.Address,
		Handler: router,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Starting server on %s", cfg.Server.Address)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited properly")
}

func runMigrations(db *sql.DB) error {
	// Здесь можно добавить логику миграций (например, с помощью github.com/golang-migrate/migrate)
	// Для простоты можно выполнить начальный SQL скрипт
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS migrations (
		id SERIAL PRIMARY KEY,
		version INT NOT NULL,
		applied_at TIMESTAMP NOT NULL DEFAULT NOW()
	)`)
	return err
}

type Repositories struct {
	DepartmentRepo *repository.DepartmentRepository
	EmployeeRepo   *repository.EmployeeRepository
	AssetRepo      *repository.AssetRepository
	StatusRepo     *repository.StatusRepository
	LocationRepo   *repository.LocationRepository
	TransferRepo   *repository.TransferRepository
}

func initializeRepositories(db *sql.DB) *Repositories {
	return &Repositories{
		DepartmentRepo: repository.NewDepartmentRepository(db),
		EmployeeRepo:   repository.NewEmployeeRepository(db),
		AssetRepo:      repository.NewAssetRepository(db),
		StatusRepo:     repository.NewStatusRepository(db),
		LocationRepo:   repository.NewLocationRepository(db),
		TransferRepo:   repository.NewTransferRepository(db),
	}
}

type Services struct {
	DepartmentService *services.DepartmentService
	EmployeeService   *services.EmployeeService
	AssetService      *services.AssetService
	LocationService   *services.LocationService
	TransferService   *services.TransferService
	AuthService       *services.AuthService
	ReportService     *services.ReportService
}

func initializeServices(repos *Repositories, cfg *config.Config) *Services {
	return &Services{
		DepartmentService: services.NewDepartmentService(repos.DepartmentRepo, repos.EmployeeRepo),
		EmployeeService:   services.NewEmployeeService(repos.EmployeeRepo, repos.DepartmentRepo),
		AssetService: services.NewAssetService(
			repos.AssetRepo,
			repos.StatusRepo,
			repos.LocationRepo,
			repos.DepartmentRepo,
		),
		LocationService: services.NewLocationService(repos.LocationRepo),
		TransferService: services.NewTransferService(
			repos.TransferRepo,
			repos.AssetRepo,
			repos.EmployeeRepo,
			repos.LocationRepo,
		),
		AuthService: services.NewAuthService(
			repos.EmployeeRepo,
			cfg.Auth.SecretKey,
			cfg.Auth.TokenExpiry,
		),
		ReportService: services.NewReportService(
			services.NewAssetService(repos.AssetRepo, repos.StatusRepo, repos.LocationRepo, repos.DepartmentRepo),
			services.NewTransferService(repos.TransferRepo, repos.AssetRepo, repos.EmployeeRepo, repos.LocationRepo),
		),
	}
}

// Функции для обработки статических страниц
func serveIndex(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./frontend/static/index.html")
}

func serveLogin(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./frontend/static/login.html")
}

func serveRegister(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./frontend/static/register.html")
}

func serveAssets(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./frontend/static/assets.html")
}

func serveEmployees(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./frontend/static/employees.html")
}

func serveDepartments(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./frontend/static/departments.html")
}

func serveTransfers(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./frontend/static/transfers.html")
}

func serveReports(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./frontend/static/reports.html")
}
