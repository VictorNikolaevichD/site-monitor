package ping

import (
	"encoding/json"
	"net/http"

	"gitlab.com/Dokuchaevvn/site-monitor/internal/handler/ping/dto"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) Ping(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	_ = json.NewEncoder(w).Encode(dto.PingResponse{Message: "pong"})
}
