package health

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"gitlab.com/Dokuchaevvn/site-monitor/internal/handler/health/dto"
	"gitlab.com/Dokuchaevvn/site-monitor/internal/health"
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

// Health godoc
// @Summary Проверить готовность к работе
// @Description Возвращает статус, uptime, текущее время и версию приложения
// @Tags system
// @Produce json
// @Success 200 {object} dto.HealthResponse
// @Failure 503 {object} dto.HealthResponse
// @Router /health [get]
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
