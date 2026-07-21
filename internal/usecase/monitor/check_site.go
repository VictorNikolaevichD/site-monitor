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
	UpdateLastCheckByID(ctx context.Context, id uuid.UUID, checkStatus domain.CheckStatus) (domain.Site, error)
}

type siteChecker interface {
	Check(url string) checker.Result
}

type CheckSiteUseCase struct {
	repo        siteRepository
	siteChecker siteChecker
	pool        *pgxpool.Pool
	logger      *slog.Logger
}

func NewCheckSiteUseCase(
	repository siteRepository,
	checker siteChecker,
	pool *pgxpool.Pool,
	logger *slog.Logger,
) *CheckSiteUseCase {
	return &CheckSiteUseCase{
		repo:        repository,
		siteChecker: checker,
		pool:        pool,
		logger:      logger,
	}
}

func (u *CheckSiteUseCase) Execute(ctx context.Context) {
	ctx = db.WithConn(ctx, u.pool)

	sites, err := u.repo.GetAll(ctx)
	if err != nil {
		u.logger.Error("failed to get sites", "error", err)
		return
	}

	var result checker.Result
	for _, v := range sites {
		result = u.siteChecker.Check(v.URL)

		if result.Error != nil {
			u.logger.Error("site check failed", "status", "NOT ok", "url", v.URL, "error", result.Error)
			if _, err := u.repo.UpdateLastCheckByID(ctx, v.ID, domain.CheckStatus{
				Availability: domain.Unavailable,
				Code:         result.Code,
				CheckedAt:    time.Now(),
				Duration:     result.Duration,
				Error:        result.Error.Error(),
			}); err != nil {
				u.logger.Error("failed to update last check", "site_id", v.ID, "error", err)
			}
			continue
		}

		if !result.AvailabilityStatus {
			u.logger.Warn("site unavailable", "status", "NOT ok", "code", result.Code, "url", v.URL)
			if _, err := u.repo.UpdateLastCheckByID(ctx, v.ID, domain.CheckStatus{
				Availability: domain.Unavailable,
				Code:         result.Code,
				CheckedAt:    time.Now(),
				Duration:     result.Duration,
			}); err != nil {
				u.logger.Error("failed to update last check", "site_id", v.ID, "error", err)
			}
			continue
		}

		u.logger.Info("site available", "status", "ok", "code", result.Code, "url", v.URL)
		if _, err := u.repo.UpdateLastCheckByID(ctx, v.ID, domain.CheckStatus{
			Availability: domain.Available,
			Code:         result.Code,
			CheckedAt:    time.Now(),
			Duration:     result.Duration,
		}); err != nil {
			u.logger.Error("failed to update last check", "site_id", v.ID, "error", err)
		}
	}
}
