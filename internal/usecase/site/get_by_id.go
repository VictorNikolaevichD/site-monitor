package site

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"gitlab.com/Dokuchaevvn/site-monitor/internal/db"
	domain "gitlab.com/Dokuchaevvn/site-monitor/internal/domain/site"
)

type GetByIDCommand struct {
	ID uuid.UUID
}

type GetByIDUseCase struct {
	sites getByIDRepository
	pool  *pgxpool.Pool
}

func NewGetByIDUseCase(sites getByIDRepository, pool *pgxpool.Pool) *GetByIDUseCase {
	return &GetByIDUseCase{
		sites: sites,
		pool:  pool,
	}
}

func (u *GetByIDUseCase) Execute(ctx context.Context, command GetByIDCommand) (domain.Site, error) {
	ctx = db.WithConn(ctx, u.pool)
	return u.sites.GetByID(ctx, command.ID)
}
