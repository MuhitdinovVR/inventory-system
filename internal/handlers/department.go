package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"inventory-system/internal/models"
)

func (h *Handler) GetAllDepartments(w http.ResponseWriter, r *http.Request) {
	departments, err := h.departmentService.GetAllDepartments(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respondWithJSON(w, http.StatusOK, departments)
}

func (h *Handler) GetDepartment(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid department ID", http.StatusBadRequest)
		return
	}

	department, err := h.departmentService.GetDepartmentByID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	respondWithJSON(w, http.StatusOK, department)
}

func (h *Handler) CreateDepartment(w http.ResponseWriter, r *http.Request) {
	var department models.Department
	if err := json.NewDecoder(r.Body).Decode(&department); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := h.departmentService.CreateDepartment(r.Context(), department)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respondWithJSON(w, http.StatusCreated, map[string]int{"id": id})
}

func (h *Handler) UpdateDepartment(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid department ID", http.StatusBadRequest)
		return
	}

	var department models.Department
	if err := json.NewDecoder(r.Body).Decode(&department); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	department.ID = id

	if err := h.departmentService.UpdateDepartment(r.Context(), department); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) DeleteDepartment(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid department ID", http.StatusBadRequest)
		return
	}

	if err := h.departmentService.DeleteDepartment(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) GetDepartmentEmployees(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid department ID", http.StatusBadRequest)
		return
	}

	employees, err := h.departmentService.GetEmployeesByDepartment(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respondWithJSON(w, http.StatusOK, employees)
}
