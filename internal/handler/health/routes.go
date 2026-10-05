package health

import (
	"net/http"
)

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /health/live", h.Live)
	mux.HandleFunc("GET /api/v1/health/live", h.Live)
	mux.HandleFunc("GET /health/ready", h.Health)
	mux.HandleFunc("GET /api/v1/health/ready", h.Health)
	mux.HandleFunc("GET /api/v1/health", h.Health)
}
