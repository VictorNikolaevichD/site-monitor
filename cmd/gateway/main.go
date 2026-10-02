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

	"gitlab.com/Dokuchaevvn/site-monitor/internal/gateway/client"
	"gitlab.com/Dokuchaevvn/site-monitor/internal/gateway/config"
	sitehandler "gitlab.com/Dokuchaevvn/site-monitor/internal/gateway/handler/site"
	"gitlab.com/Dokuchaevvn/site-monitor/internal/gateway/server"
	"gitlab.com/Dokuchaevvn/site-monitor/internal/middleware"
)

func main() {
	bootstrapLogger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	bootstrapLogger.Info("gateway started.")

	cfg, err := getConfig(bootstrapLogger)
	if err != nil {
		return
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: cfg.SlogLevel(),
	}))

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(signals)

	monitorClient, conn, err := client.NewMonitorClient(cfg.GRPCAddr, cfg.GRPCTimeout, logger)
	if err != nil {
		logger.Error("failed to create gRPC client", "error", err)
		return
	}
	defer func() {
		_ = conn.Close()
	}()

	logger.Info("monitor gRPC client created", "grpc_addr", cfg.GRPCAddr)

	mux := http.NewServeMux()

	siteHandler := sitehandler.NewHandler(monitorClient)
	siteHandler.RegisterRoutes(mux)

	httpHandler := middleware.RequestID(mux)

	httpServer := server.NewServer(cfg, logger, httpHandler)
	go runServer(httpServer, cfg, logger)

	sig := <-signals
	logger.Info("shutdown signal received", "signal", sig.String())

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		logger.Error("http server shutdown failed", "error", err)
	}

	logger.Info("gateway stopped.")
}

func getConfig(logger *slog.Logger) (*config.Config, error) {
	configPath := flag.String("config", "gateway.yaml", "path to YAML config file")
	flag.Parse()

	cfg, err := config.Load(*configPath)
	if err != nil {
		logger.Error("failed to load config", "error", err)
		return nil, err
	}

	logger.Info("config parsed",
		"http_addr", cfg.HTTPAddr,
		"grpc_addr", cfg.GRPCAddr,
		"log_level", cfg.LogLevel,
	)
	return cfg, nil
}

func runServer(server *http.Server, cfg *config.Config, logger *slog.Logger) {
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Error("failed to start http server", "http_addr", cfg.HTTPAddr)
	}
}
