package ping

import (
	"encoding/json"
	"net/http"

	"github.com/ViktorNikolaevichD/site-monitor/internal/handler/ping/dto"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

// Ping godoc
// @Summary Проверить доступность сервиса
// @Description Возвращает pong, если HTTP-сервер работает
// @Tags system
// @Produce json
// @Param X-Request-ID header string false "Идентификатор запроса для трассировки"
// @Success 200 {object} dto.PingResponse
// @Header 200 {string} X-Request-ID "Идентификатор запроса"
// @Router /ping [get]
func (h *Handler) Ping(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	_ = json.NewEncoder(w).Encode(dto.PingResponse{Message: "pong"})
}
