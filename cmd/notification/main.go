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
	eventhandler "gitlab.com/Dokuchaevvn/site-monitor/internal/notification/handler/event"
	healthhandler "gitlab.com/Dokuchaevvn/site-monitor/internal/notification/handler/health"
	"gitlab.com/Dokuchaevvn/site-monitor/internal/notification/health"
	"gitlab.com/Dokuchaevvn/site-monitor/internal/notification/messaging"
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

	consumer := messaging.NewKafkaConsumer(cfg.Kafka.Broker, cfg.Kafka.GroupID, cfg.Kafka.Topic, logger)
	defer func() {
		if err := consumer.Close(); err != nil {
			logger.Error("failed to close kafka consumer", "error", err)
		}
	}()
	logger.Info("kafka consumer created", "broker", cfg.Kafka.Broker, "topic", cfg.Kafka.Topic)

	healthChecker := health.NewHealthChecker(buildinfo.Version, startedAt)
	healthHandler := healthhandler.NewHandler(healthChecker)

	eventHandler := eventhandler.NewHandler(logger)

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(signals)

	// TODO: возможно, добавить рекавери медлвейр
	router := handler.NewRouter(healthHandler)

	httpServer := httpserver.NewServer(cfg, logger, router)
	go runServer(httpServer, cfg, logger)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	consumeDone := make(chan struct{})
	go func() {
		defer close(consumeDone)
		runConsume(ctx, consumer, eventHandler.SiteEvent, logger)
	}()

	sig := <-signals
	logger.Info("shutdown signal received", "signal", sig.String())
	cancel()
	<-consumeDone

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
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

func runConsume(
	ctx context.Context,
	consumer *messaging.KafkaConsumer,
	handler messaging.EventHandler,
	logger *slog.Logger,
) {
	if err := consumer.Consume(ctx, handler); err != nil {
		logger.Error("kafka consumer stopped", "error", err)
	}
}
