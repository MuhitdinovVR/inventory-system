package app

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"inventory-system/internal/handlers"
	"log"
	"net/http"
)

func NewRouter(
	h *handlers.Handler,
) *chi.Mux {
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Auth routes
	r.Group(func(r chi.Router) {
		r.Post("/login", h.Login)
		r.Post("/register", h.Register)
	})

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Inventory System API is running"))
	})

	// Protected routes
	r.Group(func(r chi.Router) {
		r.Use(h.AuthMiddleware)

		// Assets
		r.Route("/assets", func(r chi.Router) {
			r.Get("/", h.GetAllAssets)
			r.Post("/", h.CreateAsset)
			r.Get("/{id}", h.GetAsset)
			r.Put("/{id}", h.UpdateAsset)
			r.Delete("/{id}", h.DeleteAsset)
			r.Get("/{id}/transfers", h.GetAssetTransfers)
		})

		// Employees
		r.Route("/employees", func(r chi.Router) {
			r.Get("/", h.GetAllEmployees)
			r.Post("/", h.CreateEmployee)
			r.Get("/{id}", h.GetEmployee)
			r.Put("/{id}", h.UpdateEmployee)
			r.Delete("/{id}", h.DeleteEmployee)
		})

		// Departments
		r.Route("/departments", func(r chi.Router) {
			r.Get("/", h.GetAllDepartments)
			r.Post("/", h.CreateDepartment)
			r.Get("/{id}", h.GetDepartment)
			r.Put("/{id}", h.UpdateDepartment)
			r.Delete("/{id}", h.DeleteDepartment)
			r.Get("/{id}/employees", h.GetDepartmentEmployees)
		})

		// Locations
		r.Route("/locations", func(r chi.Router) {
			r.Get("/", h.GetAllLocations)
			r.Post("/", h.CreateLocation)
			r.Get("/{id}", h.GetLocation)
			r.Put("/{id}", h.UpdateLocation)
			r.Delete("/{id}", h.DeleteLocation)
		})

		// Transfers
		r.Route("/transfers", func(r chi.Router) {
			r.Get("/", h.GetAllTransfers)
			r.Post("/", h.CreateTransfer)
			r.Get("/{id}", h.GetTransfer)
		})

		// Reports
		r.Route("/reports", func(r chi.Router) {
			r.Get("/assets-by-status", h.GetAssetsByStatusReport)
			r.Get("/transfers", h.GetTransfersReport)
			r.Get("/department-costs", h.GetDepartmentCostsReport)
		})
	})
	// Логирование роутов
	walkFunc := func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		log.Printf("%s %s\n", method, route)
		return nil
	}

	if err := chi.Walk(r, walkFunc); err != nil {
		log.Printf("Logging err: %s\n", err.Error())
	}
	return r
}
