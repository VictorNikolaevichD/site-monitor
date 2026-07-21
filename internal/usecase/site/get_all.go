package site

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"gitlab.com/Dokuchaevvn/site-monitor/internal/db"
	domain "gitlab.com/Dokuchaevvn/site-monitor/internal/domain/site"
)

type getAllRepository interface {
	GetAll(ctx context.Context) ([]domain.Site, error)
}

type GetAllUseCase struct {
	repo getAllRepository
	pool *pgxpool.Pool
}

func NewGetAllUseCase(repository getAllRepository, pool *pgxpool.Pool) *GetAllUseCase {
	return &GetAllUseCase{
		repo: repository,
		pool: pool,
	}
}

func (u *GetAllUseCase) Execute(ctx context.Context) ([]domain.Site, error) {
	ctx = db.WithConn(ctx, u.pool)
	return u.repo.GetAll(ctx)
}
