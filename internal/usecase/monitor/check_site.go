package monitor

import (
	"log/slog"

	checker "gitlab.com/Dokuchaevvn/site-monitor/internal/checker"
	domain "gitlab.com/Dokuchaevvn/site-monitor/internal/domain/site"
)

type getAllRepository interface {
	GetAll() []domain.Site
}

type siteChecker interface {
	Check(url string) checker.Result
}

type CheckSiteUseCase struct {
	repo        getAllRepository
	siteChecker siteChecker
	logger      *slog.Logger
}

func NewCheckSiteUseCase(repository getAllRepository, checker siteChecker, logger *slog.Logger) *CheckSiteUseCase {
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
			continue
		}

		if !result.AvailabilityStatus {
			u.logger.Warn("site unavailable", "status", "NOT ok", "code", result.Code, "url", v.URL)
			continue
		}

		u.logger.Info("site available", "status", "ok", "code", result.Code, "url", v.URL)
	}
}
