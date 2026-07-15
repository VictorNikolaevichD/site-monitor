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

	"gitlab.com/Dokuchaevvn/site-monitor/internal/checker"
	"gitlab.com/Dokuchaevvn/site-monitor/internal/config"
	domain "gitlab.com/Dokuchaevvn/site-monitor/internal/domain/site"
	"gitlab.com/Dokuchaevvn/site-monitor/internal/handler"
	healthhandler "gitlab.com/Dokuchaevvn/site-monitor/internal/handler/health"
	pinghandler "gitlab.com/Dokuchaevvn/site-monitor/internal/handler/ping"
	sitehandler "gitlab.com/Dokuchaevvn/site-monitor/internal/handler/site"
	"gitlab.com/Dokuchaevvn/site-monitor/internal/health"
	siterepo "gitlab.com/Dokuchaevvn/site-monitor/internal/repository/site"
	"gitlab.com/Dokuchaevvn/site-monitor/internal/scheduler"
	"gitlab.com/Dokuchaevvn/site-monitor/internal/server"
	monitorusecase "gitlab.com/Dokuchaevvn/site-monitor/internal/usecase/monitor"
	siteusecase "gitlab.com/Dokuchaevvn/site-monitor/internal/usecase/site"
)

func main() {
	startedAt := time.Now()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	logger.Info("site monitor started.")

	cfg, err := getConfig(logger)
	if err != nil {
		return
	}

	sites := toDomainSites(cfg.Sites, logger)
	siteRepository := siterepo.NewRepository(sites)
	siteGetAllUseCase := siteusecase.NewGetAllUseCase(siteRepository)
	siteAddUseCase := siteusecase.NewAddUseCase(siteRepository)
	siteDeleteUseCase := siteusecase.NewDeleteUseCase(siteRepository)
	siteGetStatusUseCase := siteusecase.NewGetStatusUseCase(siteRepository)
	siteHandler := sitehandler.NewHandler(siteGetAllUseCase, siteAddUseCase, siteDeleteUseCase, siteGetStatusUseCase)
	checkSiteUseCase := monitorusecase.NewCheckSiteUseCase(siteRepository, checker.NewChecker(), logger)

	pingHandler := pinghandler.NewHandler()
	healthChecker := health.NewHealthChecker("v1.1.0-beta", startedAt)
	healthHandler := healthhandler.NewHandler(healthChecker)

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(signals)

	router := handler.NewRouter(siteHandler, pingHandler, healthHandler)

	httpServer := server.NewServer(cfg, logger, router)
	go runServer(httpServer, cfg, logger)

	sh := scheduler.New(checkSiteUseCase, cfg.Interval, logger)
	sh.Start()

	sig := <-signals
	logger.Info("shutdown signal received", "signal", sig.String())

	sh.Stop()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

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

func toDomainSites(cfgSites []config.Site, logger *slog.Logger) []domain.Site {
	sites := make([]domain.Site, 0, len(cfgSites))

	for _, cfgS := range cfgSites {
		site, err := domain.NewSite(cfgS.URL, cfgS.Name)
		if err != nil {
			logger.Error("failed to create site", "error", err, "url", cfgS.URL, "name", cfgS.Name)
			continue
		}
		sites = append(sites, site)
	}

	return sites
}
