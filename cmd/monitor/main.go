package main

import (
	"flag"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"gitlab.com/Dokuchaevvn/site-monitor/internal/config"
	"gitlab.com/Dokuchaevvn/site-monitor/internal/scheduler"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	logger.Info("site monitor started.")

	signals := make(chan os.Signal, 1)

	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)

	configPath := flag.String("config", "config.yaml", "path to YAML config file")
	flag.Parse()

	cfg, err := config.Load(*configPath)
	if err != nil {
		logger.Error("failed to load config", "error", err, "config_path", *configPath)
		return
	}

	logger.Info("config parsed", "sites_count", len(cfg.Sites), "interval", cfg.Interval)

	sh := scheduler.New(cfg.Sites, cfg.Interval, logger)
	sh.Start()

	sig := <-signals

	logger.Info(
		"shutdown signal received",
		"signal", sig.String(),
	)

	sh.Stop()

	logger.Info("site monitor stopped.")
}
