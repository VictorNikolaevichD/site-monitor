package health

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/ViktorNikolaevichD/site-monitor/internal/handler/health/dto"
	"github.com/ViktorNikolaevichD/site-monitor/internal/health"
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
// @Description Пингует Postgres. 200, если база доступна, и 503, если нет. Тот же ответ на /health/ready. Не использовать как liveness: сбой базы уберёт под из Service, но не должен его перезапускать.
// @Tags system
// @Produce json
// @Param X-Request-ID header string false "Идентификатор запроса для трассировки"
// @Success 200 {object} dto.HealthResponse
// @Failure 503 {object} dto.HealthResponse
// @Header 200 {string} X-Request-ID "Идентификатор запроса"
// @Header 503 {string} X-Request-ID "Идентификатор запроса"
// @Router /health [get]
// @Router /health/ready [get]
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

// Live godoc
// @Summary Проверить, что процесс жив
// @Description Не проверяет Postgres и Kafka. 200 означает, что HTTP-сервер обрабатывает запросы. Этот путь использует liveness probe.
// @Tags system
// @Produce json
// @Param X-Request-ID header string false "Идентификатор запроса для трассировки"
// @Success 200 {object} dto.LiveResponse
// @Header 200 {string} X-Request-ID "Идентификатор запроса"
// @Router /health/live [get]
func (h *Handler) Live(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(dto.LiveResponse{Status: string(health.StatusHealthy)})
}
