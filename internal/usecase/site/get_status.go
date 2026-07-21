package site

import (
	"context"

	"github.com/google/uuid"
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
}

func NewGetStatusUseCase(repository getByIDRepository) *GetStatusUseCase {
	return &GetStatusUseCase{
		repo: repository,
	}
}

func (u *GetStatusUseCase) Execute(ctx context.Context, command GetStatusCommand) (domain.Site, error) {
	return u.repo.GetByID(ctx, command.ID)
}
