package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"inventory-system/internal/models"
)

func (h *Handler) GetAllTransfers(w http.ResponseWriter, r *http.Request) {
	transfers, err := h.transferService.GetAllTransfers(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respondWithJSON(w, http.StatusOK, transfers)
}

func (h *Handler) GetTransfer(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid transfer ID", http.StatusBadRequest)
		return
	}

	transfer, err := h.transferService.GetTransferByID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	respondWithJSON(w, http.StatusOK, transfer)
}

func (h *Handler) CreateTransfer(w http.ResponseWriter, r *http.Request) {
	var transfer models.AssetTransfer
	if err := json.NewDecoder(r.Body).Decode(&transfer); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := h.transferService.CreateTransfer(r.Context(), transfer)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respondWithJSON(w, http.StatusCreated, map[string]int{"id": id})
}

func (h *Handler) GetTransfersReport(w http.ResponseWriter, r *http.Request) {
	fromStr := r.URL.Query().Get("from")
	toStr := r.URL.Query().Get("to")

	var from, to time.Time
	var err error

	if fromStr != "" {
		from, err = time.Parse("2006-01-02", fromStr)
		if err != nil {
			http.Error(w, "Invalid from date format", http.StatusBadRequest)
			return
		}
	} else {
		from = time.Now().AddDate(0, -1, 0) // Default: 1 month ago
	}

	if toStr != "" {
		to, err = time.Parse("2006-01-02", toStr)
		if err != nil {
			http.Error(w, "Invalid to date format", http.StatusBadRequest)
			return
		}
	} else {
		to = time.Now() // Default: now
	}

	transfers, err := h.reportService.GenerateTransfersReport(r.Context(), from, to)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respondWithJSON(w, http.StatusOK, transfers)
}
