package site

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"gitlab.com/Dokuchaevvn/site-monitor/internal/db"
	domain "gitlab.com/Dokuchaevvn/site-monitor/internal/domain/site"
)

const (
	defaultHistoryLimit = 20
	maxHistoryLimit     = 100
)

type getHistorySiteRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (domain.Site, error)
}

type getHistoryCheckResultsRepository interface {
	GetBySiteID(ctx context.Context, siteID uuid.UUID, limit, offset int) ([]domain.CheckResult, error)
	CountBySiteID(ctx context.Context, siteID uuid.UUID) (int, error)
}

type GetHistoryCommand struct {
	SiteID uuid.UUID
	Limit  int
	Offset int
}

type GetHistoryResult struct {
	Items  []domain.CheckResult
	Total  int
	Limit  int
	Offset int
}

type GetHistoryUseCase struct {
	sites        getHistorySiteRepository
	checkResults getHistoryCheckResultsRepository
	pool         *pgxpool.Pool
}

func NewGetHistoryUseCase(
	sites getHistorySiteRepository,
	checkResults getHistoryCheckResultsRepository,
	pool *pgxpool.Pool,
) *GetHistoryUseCase {
	return &GetHistoryUseCase{
		sites:        sites,
		checkResults: checkResults,
		pool:         pool,
	}
}

func (u *GetHistoryUseCase) Execute(ctx context.Context, command GetHistoryCommand) (GetHistoryResult, error) {
	limit := command.Limit
	if limit == 0 {
		limit = defaultHistoryLimit
	}
	if limit < 0 {
		return GetHistoryResult{}, ErrInvalidLimit
	}
	if limit > maxHistoryLimit {
		limit = maxHistoryLimit
	}
	if command.Offset < 0 {
		return GetHistoryResult{}, ErrInvalidOffset
	}

	ctx = db.WithConn(ctx, u.pool)

	if _, err := u.sites.GetByID(ctx, command.SiteID); err != nil {
		return GetHistoryResult{}, err
	}

	total, err := u.checkResults.CountBySiteID(ctx, command.SiteID)
	if err != nil {
		return GetHistoryResult{}, err
	}

	items, err := u.checkResults.GetBySiteID(ctx, command.SiteID, limit, command.Offset)
	if err != nil {
		return GetHistoryResult{}, err
	}

	return GetHistoryResult{
		Items:  items,
		Total:  total,
		Limit:  limit,
		Offset: command.Offset,
	}, nil
}
