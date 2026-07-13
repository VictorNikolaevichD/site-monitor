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
	domain "gitlab.com/Dokuchaevvn/site-monitor/internal/domain/site"
	"gitlab.com/Dokuchaevvn/site-monitor/internal/handler"
	pingHandler "gitlab.com/Dokuchaevvn/site-monitor/internal/handler/ping"
	siteHandler "gitlab.com/Dokuchaevvn/site-monitor/internal/handler/site"
	siteRepo "gitlab.com/Dokuchaevvn/site-monitor/internal/repository/site"
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

	sites := toDomainSites(cfg.Sites)
	siteRepository := siteRepo.NewRepository(sites)
	siteHandler := siteHandler.NewHandler(siteRepository)

	pingHandler := pingHandler.NewHandler()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(signals)

	mux := http.NewServeMux()
	handler.NewRouter(mux, siteHandler)
	handler.NewRouter(mux, pingHandler)

	httpServer := server.NewServer(cfg, logger, mux)
	go runServer(httpServer, cfg, logger)

	sh := scheduler.New(cfg.Sites, cfg.Interval, logger)
	sh.Start()

	sig := <-signals
	logger.Info("shutdown signal received", "signal", sig.String())

	sh.Stop()

	if err := httpServer.Shutdown(ctx); err != nil {
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

func toDomainSites(cfgSites []config.Site) []domain.Site {
	sites := make([]domain.Site, 0, len(cfgSites))

	for i, cfgS := range cfgSites {
		sites = append(sites, domain.Site{
			ID:   int32(i + 1),
			URL:  cfgS.URL,
			Name: cfgS.Name,
		})
	}

	return sites
}
