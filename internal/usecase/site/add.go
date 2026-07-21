package site

import (
	"context"

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
	repo addIfAbsentRepository
	pool *pgxpool.Pool
}

func NewAddUseCase(repository addIfAbsentRepository, pool *pgxpool.Pool) *AddUseCase {
	return &AddUseCase{
		repo: repository,
		pool: pool,
	}
}

func (u *AddUseCase) Execute(ctx context.Context, command AddCommand) (domain.Site, error) {
	site, err := domain.NewSite(command.URL, command.Name)
	if err != nil {
		return domain.Site{}, err
	}

	tx, err := u.pool.Begin(ctx)
	if err != nil {
		return domain.Site{}, domain.ErrStorage
	}
	defer tx.Rollback(ctx)

	created, err := u.repo.AddIfAbsent(db.WithConn(ctx, tx), site)
	if err != nil {
		return domain.Site{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return domain.Site{}, domain.ErrStorage
	}

	return created, nil
}
