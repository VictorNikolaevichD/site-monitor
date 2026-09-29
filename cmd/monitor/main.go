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

	monitorv1 "gitlab.com/Dokuchaevvn/site-monitor/gen/go/monitor/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"gitlab.com/Dokuchaevvn/site-monitor/internal/buildinfo"
	"gitlab.com/Dokuchaevvn/site-monitor/internal/checker"
	"gitlab.com/Dokuchaevvn/site-monitor/internal/config"
	"gitlab.com/Dokuchaevvn/site-monitor/internal/db"
	grpcserver "gitlab.com/Dokuchaevvn/site-monitor/internal/grpc"
	grpcmonitor "gitlab.com/Dokuchaevvn/site-monitor/internal/grpc/monitor"
	"gitlab.com/Dokuchaevvn/site-monitor/internal/handler"
	healthhandler "gitlab.com/Dokuchaevvn/site-monitor/internal/handler/health"
	pinghandler "gitlab.com/Dokuchaevvn/site-monitor/internal/handler/ping"
	sitehandler "gitlab.com/Dokuchaevvn/site-monitor/internal/handler/site"
	swaggerhandler "gitlab.com/Dokuchaevvn/site-monitor/internal/handler/swagger"
	"gitlab.com/Dokuchaevvn/site-monitor/internal/health"
	"gitlab.com/Dokuchaevvn/site-monitor/internal/messaging"
	"gitlab.com/Dokuchaevvn/site-monitor/internal/middleware"
	checkresultrepo "gitlab.com/Dokuchaevvn/site-monitor/internal/repository/checkresult"
	siterepo "gitlab.com/Dokuchaevvn/site-monitor/internal/repository/site"
	"gitlab.com/Dokuchaevvn/site-monitor/internal/scheduler"
	"gitlab.com/Dokuchaevvn/site-monitor/internal/server"
	monitorusecase "gitlab.com/Dokuchaevvn/site-monitor/internal/usecase/monitor"
	siteusecase "gitlab.com/Dokuchaevvn/site-monitor/internal/usecase/site"
)

// @title Site Monitor API
// @version v1.0
// @description REST API для управления сайтами и просмотра результатов мониторинга.
// @description Опционально передавайте X-Request-ID в запросе; сервис вернёт его (или сгенерированный ID) в заголовке ответа X-Request-ID.
// @accept json
// @produce json
// @BasePath /api/v1
// @schemes http
func main() {
	startedAt := time.Now()

	bootstrapLogger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	bootstrapLogger.Info("site monitor started.")

	cfg, err := getConfig(bootstrapLogger)
	if err != nil {
		return
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: cfg.SlogLevel(),
	}))

	dbConn, err := db.OpenPostgres(context.Background(), cfg.Database)
	if err != nil {
		logger.Error("failed to connect to postgres", "error", err)
		return
	}
	defer dbConn.Close()
	logger.Info("postgres connected")

	producer := messaging.NewKafkaProducer(cfg.Kafka.Broker, cfg.Kafka.Topic, logger)
	defer func() {
		if err := producer.Close(); err != nil {
			logger.Error("failed to close kafka producer", "error", err)
		}
	}()
	logger.Info("kafka producer created", "broker", cfg.Kafka.Broker, "topic", cfg.Kafka.Topic)

	swaggerHandler := swaggerhandler.NewHandler()

	siteRepository := siterepo.NewPostgresSiteRepository(logger)
	checkResultRepository := checkresultrepo.NewPostgresCheckResultRepository(logger)
	siteGetByIDUseCase := siteusecase.NewGetByIDUseCase(siteRepository, dbConn)
	siteGetAllUseCase := siteusecase.NewGetAllUseCase(siteRepository, dbConn)
	siteAddUseCase := siteusecase.NewAddUseCase(siteRepository, dbConn, logger)
	siteDeleteUseCase := siteusecase.NewDeleteUseCase(siteRepository, checkResultRepository, dbConn, logger)
	siteGetStatusUseCase := siteusecase.NewGetStatusUseCase(siteRepository, checkResultRepository, dbConn)
	siteGetHistoryUseCase := siteusecase.NewGetHistoryUseCase(siteRepository, checkResultRepository, dbConn, logger)
	siteHandler := sitehandler.NewHandler(
		siteGetAllUseCase,
		siteAddUseCase,
		siteDeleteUseCase,
		siteGetStatusUseCase,
		siteGetHistoryUseCase,
	)
	checkSiteUseCase := monitorusecase.NewCheckSiteUseCase(
		siteRepository,
		checkResultRepository,
		checker.NewChecker(cfg.HTTPTimeout),
		producer,
		dbConn,
		logger,
	)

	monitorService := grpcmonitor.NewService(
		siteGetAllUseCase,
		siteGetByIDUseCase,
		siteAddUseCase,
		siteDeleteUseCase,
		siteGetStatusUseCase,
		siteGetHistoryUseCase,
	)

	pingHandler := pinghandler.NewHandler()
	healthChecker := health.NewHealthChecker(buildinfo.Version, startedAt, db.NewPostgresChecker(dbConn))
	healthHandler := healthhandler.NewHandler(healthChecker)

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(signals)

	router := handler.NewRouter(siteHandler, pingHandler, healthHandler, swaggerHandler)

	httpHandler := middleware.RequestID(middleware.Logging(middleware.Recovery(router, logger), logger))

	httpServer := server.NewServer(cfg, logger, httpHandler)
	go runServer(httpServer, cfg, logger)

	grpcServer := grpcserver.NewServer(cfg.GRPCAddr, logger, func(s *grpc.Server) {
		monitorv1.RegisterMonitorServiceServer(s, monitorService)
		reflection.Register(s)
	})
	go runGRPCServer(grpcServer, cfg, logger)

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
	grpcServer.GracefulStop()

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

	logger.Info("config parsed",
		"interval", cfg.Interval,
		"http_addr", cfg.HTTPAddr,
		"grpc_addr", cfg.GRPCAddr,
		"log_level", cfg.LogLevel,
		"http_timeout", cfg.HTTPTimeout,
		"kafka_broker", cfg.Kafka.Broker,
		"kafka_topic", cfg.Kafka.Topic,
	)
	return cfg, nil
}

func runServer(server *http.Server, cfg *config.Config, logger *slog.Logger) {
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Error("failed to start http server", "http_addr", cfg.HTTPAddr)
	}
}

type GRPCServer interface {
	Run() error
}

func runGRPCServer(server GRPCServer, cfg *config.Config, logger *slog.Logger) {
	if err := server.Run(); err != nil {
		logger.Error("failed to start grpc server", "grpc_addr", cfg.GRPCAddr)
	}
}
