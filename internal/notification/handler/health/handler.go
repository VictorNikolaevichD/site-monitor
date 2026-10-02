package health

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/ViktorNikolaevichD/site-monitor/internal/notification/handler/health/dto"
	"github.com/ViktorNikolaevichD/site-monitor/internal/notification/health"
)

type healthChecker interface {
	Check(ctx context.Context) health.Report
}

type Handler struct {
	healthChecker healthChecker
}

func NewHandler(healthChecker healthChecker) *Handler {
	return &Handler{
		healthChecker: healthChecker,
	}
}

func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	report := h.healthChecker.Check(ctx)
	statusCode := http.StatusOK
	if report.Status == health.StatusUnhealthy {
		statusCode = http.StatusServiceUnavailable
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	_ = json.NewEncoder(w).Encode(dto.ToHealthResponse(report))
}
