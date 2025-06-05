package handlers

import (
	"net/http"
)

func (h *Handler) GetAssetsByStatusReport(w http.ResponseWriter, r *http.Request) {
	report, err := h.reportService.GenerateAssetsByStatusReport(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respondWithJSON(w, http.StatusOK, report)
}

func (h *Handler) GetDepartmentCostsReport(w http.ResponseWriter, r *http.Request) {
	report, err := h.reportService.GenerateAssetsCostByDepartmentReport(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respondWithJSON(w, http.StatusOK, report)
}

func (h *Handler) GetInventoryReport(w http.ResponseWriter, r *http.Request) {
	report, err := h.reportService.GenerateInventoryReport(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respondWithJSON(w, http.StatusOK, report)
}
