package site

import (
	"net/http"
)

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v1/sites", h.GetAll)
	mux.HandleFunc("POST /api/v1/sites", h.Add)
	mux.HandleFunc("DELETE /api/v1/sites/{id}", h.DeleteByID)
	mux.HandleFunc("GET /api/v1/sites/{id}/status", h.GetStatusByID)
	mux.HandleFunc("GET /api/v1/sites/{id}/history", h.GetHistoryByID)
}
