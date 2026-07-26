package site

import (
	"context"
	"errors"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
	"gitlab.com/Dokuchaevvn/site-monitor/internal/db"
	domain "gitlab.com/Dokuchaevvn/site-monitor/internal/domain/site"
)

type addIfAbsentRepository interface {
	AddIfAbsent(ctx context.Context, site domain.Site) (domain.Site, error)
}

type AddCommand struct {
	URL  string
	Name string
}

type AddUseCase struct {
	repo   addIfAbsentRepository
	pool   *pgxpool.Pool
	logger *slog.Logger
}

func NewAddUseCase(repository addIfAbsentRepository, pool *pgxpool.Pool, logger *slog.Logger) *AddUseCase {
	return &AddUseCase{
		repo:   repository,
		pool:   pool,
		logger: logger,
	}
}

func (u *AddUseCase) Execute(ctx context.Context, command AddCommand) (domain.Site, error) {
	site, err := domain.NewSite(command.URL, command.Name)
	if err != nil {
		return domain.Site{}, err
	}

	created, err := db.WithinTxResult(ctx, u.pool, u.logger, func(ctx context.Context) (domain.Site, error) {
		return u.repo.AddIfAbsent(ctx, site)
	})
	if err != nil {
		if errors.Is(err, db.ErrTx) {
			return domain.Site{}, domain.ErrStorage
		}
		return domain.Site{}, err
	}

	return created, nil
}
