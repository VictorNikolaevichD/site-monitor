package monitor

import (
	"log/slog"
	"time"

	"github.com/google/uuid"
	checker "gitlab.com/Dokuchaevvn/site-monitor/internal/checker"
	domain "gitlab.com/Dokuchaevvn/site-monitor/internal/domain/site"
)

type siteRepository interface {
	GetAll() []domain.Site
	UpdateLastCheckByID(id uuid.UUID, checkStatus domain.CheckStatus) (domain.Site, error)
}

type siteChecker interface {
	Check(url string) checker.Result
}

type CheckSiteUseCase struct {
	repo        siteRepository
	siteChecker siteChecker
	logger      *slog.Logger
}

func NewCheckSiteUseCase(repository siteRepository, checker siteChecker, logger *slog.Logger) *CheckSiteUseCase {
	return &CheckSiteUseCase{
		repo:        repository,
		siteChecker: checker,
		logger:      logger,
	}
}

func (u *CheckSiteUseCase) Execute() {
	sites := u.repo.GetAll()

	var result checker.Result
	for _, v := range sites {
		result = u.siteChecker.Check(v.URL)

		if result.Error != nil {
			u.logger.Error("site check failed", "status", "NOT ok", "url", v.URL, "error", result.Error)
			u.repo.UpdateLastCheckByID(v.ID, domain.CheckStatus{
				Availability: domain.Unavailable,
				Code:         result.Code,
				CheckedAt:    time.Now(),
				Duration:     result.Duration,
				Error:        result.Error.Error(),
			})
			continue
		}

		if !result.AvailabilityStatus {
			u.logger.Warn("site unavailable", "status", "NOT ok", "code", result.Code, "url", v.URL)
			u.repo.UpdateLastCheckByID(v.ID, domain.CheckStatus{
				Availability: domain.Unavailable,
				Code:         result.Code,
				CheckedAt:    time.Now(),
				Duration:     result.Duration,
			})
			continue
		}

		u.logger.Info("site available", "status", "ok", "code", result.Code, "url", v.URL)
		u.repo.UpdateLastCheckByID(v.ID, domain.CheckStatus{
			Availability: domain.Available,
			Code:         result.Code,
			CheckedAt:    time.Now(),
			Duration:     result.Duration,
		})
	}
}
