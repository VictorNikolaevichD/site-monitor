package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gitlab.com/Dokuchaevvn/site-monitor/internal/notification/buildinfo"
	"gitlab.com/Dokuchaevvn/site-monitor/internal/notification/config"
	"gitlab.com/Dokuchaevvn/site-monitor/internal/notification/handler"
	healthhandler "gitlab.com/Dokuchaevvn/site-monitor/internal/notification/handler/health"
	"gitlab.com/Dokuchaevvn/site-monitor/internal/notification/health"
	httpserver "gitlab.com/Dokuchaevvn/site-monitor/internal/notification/server/http"
)

func main() {
	startedAt := time.Now()

	bootstrapLogger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	bootstrapLogger.Info("notification started.")

	cfg, err := getConfig(bootstrapLogger)
	if err != nil {
		return
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: cfg.SlogLevel(),
	}))

	// TODO: добавить кафка консьюмера
	logger.Info("kafka consumer created", "broker", cfg.Kafka.Broker, "topic", cfg.Kafka.Topic)

	healthChecker := health.NewHealthChecker(buildinfo.Version, startedAt)
	healthHandler := healthhandler.NewHandler(healthChecker)

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(signals)

	// TODO: возможно, добавить рекавери медлвейр
	router := handler.NewRouter(healthHandler)

	httpServer := httpserver.NewServer(cfg, logger, router)
	go runServer(httpServer, cfg, logger)

	sig := <-signals
	logger.Info("shutdown signal received", "signal", sig.String())

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		logger.Error("http server shutdown failed", "error", err)
	}

	logger.Info("notification stopped.")
}

func getConfig(logger *slog.Logger) (*config.Config, error) {
	cfg, err := config.Load()
	if err != nil {
		logger.Error("failed to load config", "error", err)
		return nil, err
	}

	logger.Info("config parsed",
		"http_addr", cfg.HTTPAddr,
		"log_level", cfg.LogLevel,
	)
	return cfg, nil
}

func runServer(server *http.Server, cfg *config.Config, logger *slog.Logger) {
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Error("failed to start http server", "http_addr", cfg.HTTPAddr)
	}
}
