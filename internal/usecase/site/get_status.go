package site

import (
	"context"
	"errors"

	"github.com/ViktorNikolaevichD/site-monitor/internal/db"
	domain "github.com/ViktorNikolaevichD/site-monitor/internal/domain/site"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type getByIDRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (domain.Site, error)
}

type latestCheckRepository interface {
	GetLatestBySiteID(ctx context.Context, siteID uuid.UUID) (domain.CheckStatus, error)
}

type GetStatusCommand struct {
	ID uuid.UUID
}

type GetStatusUseCase struct {
	sites        getByIDRepository
	checkResults latestCheckRepository
	pool         *pgxpool.Pool
}

func NewGetStatusUseCase(
	sites getByIDRepository,
	checkResults latestCheckRepository,
	pool *pgxpool.Pool,
) *GetStatusUseCase {
	return &GetStatusUseCase{
		sites:        sites,
		checkResults: checkResults,
		pool:         pool,
	}
}

func (u *GetStatusUseCase) Execute(ctx context.Context, command GetStatusCommand) (domain.Site, error) {
	ctx = db.WithConn(ctx, u.pool)

	site, err := u.sites.GetByID(ctx, command.ID)
	if err != nil {
		return domain.Site{}, err
	}

	latest, err := u.checkResults.GetLatestBySiteID(ctx, command.ID)
	if err != nil {
		if errors.Is(err, domain.ErrCheckResultNotFound) {
			return site, nil
		}
		return domain.Site{}, err
	}

	site.LastCheck = &latest
	return site, nil
}
