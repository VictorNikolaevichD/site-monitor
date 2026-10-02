package server

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/ViktorNikolaevichD/site-monitor/internal/gateway/config"
)

func NewServer(cfg *config.Config, logger *slog.Logger, handler http.Handler) *http.Server {
	logger.Info("start of gateway server creation", "http_addr", cfg.HTTPAddr)

	return &http.Server{
		Addr:              cfg.HTTPAddr,
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
}
