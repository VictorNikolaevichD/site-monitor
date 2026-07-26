package site

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"gitlab.com/Dokuchaevvn/site-monitor/internal/db"
	domain "gitlab.com/Dokuchaevvn/site-monitor/internal/domain/site"
)

type deleteSiteRepository interface {
	GetByIDForUpdate(ctx context.Context, id uuid.UUID) (domain.Site, error)
	DeleteByID(ctx context.Context, id uuid.UUID) error
}

type deleteCheckResultsRepository interface {
	DeleteBySiteID(ctx context.Context, siteID uuid.UUID) error
}

type DeleteCommand struct {
	ID uuid.UUID
}

type DeleteUseCase struct {
	sites        deleteSiteRepository
	checkResults deleteCheckResultsRepository
	pool         *pgxpool.Pool
}

func NewDeleteUseCase(
	sites deleteSiteRepository,
	checkResults deleteCheckResultsRepository,
	pool *pgxpool.Pool,
) *DeleteUseCase {
	return &DeleteUseCase{
		sites:        sites,
		checkResults: checkResults,
		pool:         pool,
	}
}

func (u *DeleteUseCase) Execute(ctx context.Context, command DeleteCommand) error {
	err := db.WithinTx(ctx, u.pool, func(ctx context.Context) error {
		if _, err := u.sites.GetByIDForUpdate(ctx, command.ID); err != nil {
			return err
		}

		if err := u.checkResults.DeleteBySiteID(ctx, command.ID); err != nil {
			return err
		}

		return u.sites.DeleteByID(ctx, command.ID)
	})
	if err != nil {
		if errors.Is(err, db.ErrTx) {
			return domain.ErrStorage
		}
		return err
	}

	return nil
}
