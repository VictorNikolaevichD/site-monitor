package site

import (
	"context"
	"errors"
	"log/slog"

	"github.com/ViktorNikolaevichD/site-monitor/internal/db"
	domain "github.com/ViktorNikolaevichD/site-monitor/internal/domain/site"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
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
	logger       *slog.Logger
}

func NewDeleteUseCase(
	sites deleteSiteRepository,
	checkResults deleteCheckResultsRepository,
	pool *pgxpool.Pool,
	logger *slog.Logger,
) *DeleteUseCase {
	return &DeleteUseCase{
		sites:        sites,
		checkResults: checkResults,
		pool:         pool,
		logger:       logger,
	}
}

func (u *DeleteUseCase) Execute(ctx context.Context, command DeleteCommand) error {
	err := db.WithinTx(ctx, u.pool, u.logger, func(ctx context.Context) error {
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
