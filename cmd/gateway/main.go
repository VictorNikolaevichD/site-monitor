package main

import (
	"flag"
	"log/slog"
	"os"

	"gitlab.com/Dokuchaevvn/site-monitor/internal/gateway/config"
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

	// TODO

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
