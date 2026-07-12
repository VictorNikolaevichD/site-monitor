package main

import (
	"context"
	"errors"
	"flag"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gitlab.com/Dokuchaevvn/site-monitor/internal/config"
	"gitlab.com/Dokuchaevvn/site-monitor/internal/scheduler"
	"gitlab.com/Dokuchaevvn/site-monitor/internal/server"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	logger.Info("site monitor started.")

	cfg, err := getConfig(logger)
	if err != nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(signals)

	server := server.New(cfg, logger)
	go runServer(server, cfg, logger)

	sh := scheduler.New(cfg.Sites, cfg.Interval, logger)
	sh.Start()

	sig := <-signals
	logger.Info("shutdown signal received", "signal", sig.String())

	sh.Stop()

	if err := server.Shutdown(ctx); err != nil {
		logger.Error("http server shutdown failed", "error", err)
	}

	logger.Info("site monitor stopped.")
}

func getConfig(logger *slog.Logger) (*config.Config, error) {
	configPath := flag.String("config", "config.yaml", "path to YAML config file")
	flag.Parse()

	cfg, err := config.Load(*configPath)
	if err != nil {
		logger.Error("failed to load config", "error", err, "config_path", *configPath)
		return nil, err
	}

	logger.Info("config parsed", "sites_count", len(cfg.Sites), "interval", cfg.Interval)
	return cfg, nil
}

func runServer(server *http.Server, cfg *config.Config, logger *slog.Logger) {
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Error("failed to start http server", "http_addr", cfg.HTTPAddr)
	}
}
