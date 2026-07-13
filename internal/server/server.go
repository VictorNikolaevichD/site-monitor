package server

import (
	"log/slog"
	"net/http"
	"time"

	"gitlab.com/Dokuchaevvn/site-monitor/internal/config"
)

func NewServer(cfg *config.Config, logger *slog.Logger, handler http.Handler) *http.Server {
	logger.Info("start of server creation", "http_addr", cfg.HTTPAddr)

	return &http.Server{
		Addr:              cfg.HTTPAddr,
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
	}
}
