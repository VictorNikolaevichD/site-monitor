package monitor

import (
	"context"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	checker "gitlab.com/Dokuchaevvn/site-monitor/internal/checker"
	"gitlab.com/Dokuchaevvn/site-monitor/internal/db"
	domain "gitlab.com/Dokuchaevvn/site-monitor/internal/domain/site"
)

type siteRepository interface {
	GetAll(ctx context.Context) ([]domain.Site, error)
}

type checkResultRepository interface {
	Create(ctx context.Context, siteID uuid.UUID, result domain.CheckStatus) error
}

type siteChecker interface {
	Check(url string) checker.Result
}

type CheckSiteUseCase struct {
	repo         siteRepository
	checkResults checkResultRepository
	siteChecker  siteChecker
	pool         *pgxpool.Pool
	logger       *slog.Logger
}

func NewCheckSiteUseCase(
	repository siteRepository,
	checkResults checkResultRepository,
	checker siteChecker,
	pool *pgxpool.Pool,
	logger *slog.Logger,
) *CheckSiteUseCase {
	return &CheckSiteUseCase{
		repo:         repository,
		checkResults: checkResults,
		siteChecker:  checker,
		pool:         pool,
		logger:       logger,
	}
}

func (u *CheckSiteUseCase) Execute(ctx context.Context) {
	ctx = db.WithConn(ctx, u.pool)

	sites, err := u.repo.GetAll(ctx)
	if err != nil {
		u.logger.Error("failed to get sites", "error", err)
		return
	}

	for _, v := range sites {
		result := u.siteChecker.Check(v.URL)
		status := domain.CheckStatus{
			Code:      result.Code,
			CheckedAt: time.Now().UTC(),
			Duration:  result.Duration,
		}

		switch {
		case result.Error != nil:
			u.logger.Error("site check failed", "status", "NOT ok", "url", v.URL, "error", result.Error)
			status.Availability = domain.Unavailable
			status.Error = result.Error.Error()
		case !result.AvailabilityStatus:
			u.logger.Warn("site unavailable", "status", "NOT ok", "code", result.Code, "url", v.URL)
			status.Availability = domain.Unavailable
		default:
			u.logger.Info("site available", "status", "ok", "code", result.Code, "url", v.URL)
			status.Availability = domain.Available
		}

		if err := u.checkResults.Create(ctx, v.ID, status); err != nil {
			u.logger.Error("failed to save check result", "site_id", v.ID, "error", err)
		}
	}
}
