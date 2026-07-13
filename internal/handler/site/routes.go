package site

import (
	"net/http"
)

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v1/sites", h.GetAll)
	mux.HandleFunc("POST /api/v1/sites", h.Add)
}
