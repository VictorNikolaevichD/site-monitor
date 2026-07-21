package site

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"gitlab.com/Dokuchaevvn/site-monitor/internal/db"
	domain "gitlab.com/Dokuchaevvn/site-monitor/internal/domain/site"
)

type getByIDRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (domain.Site, error)
}

type GetStatusCommand struct {
	ID uuid.UUID
}

type GetStatusUseCase struct {
	repo getByIDRepository
	pool *pgxpool.Pool
}

func NewGetStatusUseCase(repository getByIDRepository, pool *pgxpool.Pool) *GetStatusUseCase {
	return &GetStatusUseCase{
		repo: repository,
		pool: pool,
	}
}

func (u *GetStatusUseCase) Execute(ctx context.Context, command GetStatusCommand) (domain.Site, error) {
	ctx = db.WithConn(ctx, u.pool)
	return u.repo.GetByID(ctx, command.ID)
}
