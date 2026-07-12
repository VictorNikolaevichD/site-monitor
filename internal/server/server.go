package server

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"gitlab.com/Dokuchaevvn/site-monitor/internal/config"
)

type PingResponse struct {
	Message string `json:"message"`
}

func New(cfg *config.Config, logger *slog.Logger) *http.Server {
	logger.Info("start of server creation", "http_addr", cfg.HTTPAddr)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v1/ping", pingHandler)

	return &http.Server{
		Addr:    cfg.HTTPAddr,
		Handler: mux,
		ReadHeaderTimeout: 5 * time.Second,
	}
}

func pingHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(PingResponse{Message: "pong"})
}
