package handlers

import (
	"encoding/json"
	"golang.org/x/crypto/bcrypt"
	"inventory-system/internal/models"
	"net/http"
)

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var creds struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	employee, err := h.authService.Authenticate(r.Context(), creds.Email, creds.Password)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	token, expiresAt, err := h.authService.GenerateToken(employee)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"token":      token,
		"expires_at": expiresAt,
		"employee":   employee,
	})
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var request struct {
		FullName string `json:"full_name"`
		Position string `json:"position"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Хеширование пароля с использованием bcrypt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to hash password")
		return
	}

	employee := models.Employee{
		FullName:     request.FullName,
		Position:     request.Position,
		Email:        request.Email,
		PasswordHash: string(hashedPassword),
		Role:         "employee",
	}

	id, err := h.employeeService.CreateEmployee(r.Context(), employee)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, map[string]int{"id": id})
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}
