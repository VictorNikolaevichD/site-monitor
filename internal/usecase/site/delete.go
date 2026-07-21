package site

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"gitlab.com/Dokuchaevvn/site-monitor/internal/db"
	domain "gitlab.com/Dokuchaevvn/site-monitor/internal/domain/site"
)

type deleteByIDRepository interface {
	DeleteByID(ctx context.Context, id uuid.UUID) error
}

type DeleteCommand struct {
	ID uuid.UUID
}

type DeleteUseCase struct {
	repo deleteByIDRepository
	pool *pgxpool.Pool
}

func NewDeleteUseCase(repository deleteByIDRepository, pool *pgxpool.Pool) *DeleteUseCase {
	return &DeleteUseCase{
		repo: repository,
		pool: pool,
	}
}

func (u *DeleteUseCase) Execute(ctx context.Context, command DeleteCommand) error {
	tx, err := u.pool.Begin(ctx)
	if err != nil {
		return domain.ErrStorage
	}
	defer tx.Rollback(ctx)

	if err := u.repo.DeleteByID(db.WithConn(ctx, tx), command.ID); err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return domain.ErrStorage
	}

	return nil
}
