package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"inventory-system/internal/models"
)

func (h *Handler) GetAllEmployees(w http.ResponseWriter, r *http.Request) {
	employees, err := h.employeeService.GetAllEmployees(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respondWithJSON(w, http.StatusOK, employees)
}

func (h *Handler) GetEmployee(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid employee ID", http.StatusBadRequest)
		return
	}

	employee, err := h.employeeService.GetEmployeeByID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	respondWithJSON(w, http.StatusOK, employee)
}

func (h *Handler) CreateEmployee(w http.ResponseWriter, r *http.Request) {
	log.Println("CreateEmployee called")

	var employee models.Employee
	if err := json.NewDecoder(r.Body).Decode(&employee); err != nil {
		log.Printf("Error decoding request: %v", err)
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	log.Printf("Creating employee: %+v", employee)

	id, err := h.employeeService.CreateEmployee(r.Context(), employee)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respondWithJSON(w, http.StatusCreated, map[string]int{"id": id})
}

func (h *Handler) UpdateEmployee(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid employee ID", http.StatusBadRequest)
		return
	}

	var employee models.Employee
	if err := json.NewDecoder(r.Body).Decode(&employee); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	employee.ID = id

	if err := h.employeeService.UpdateEmployee(r.Context(), employee); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) DeleteEmployee(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid employee ID", http.StatusBadRequest)
		return
	}

	if err := h.employeeService.DeleteEmployee(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
