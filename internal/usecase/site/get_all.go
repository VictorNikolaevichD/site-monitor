package site

import (
	"context"

	domain "gitlab.com/Dokuchaevvn/site-monitor/internal/domain/site"
)

type getAllRepository interface {
	GetAll(ctx context.Context) ([]domain.Site, error)
}

type GetAllUseCase struct {
	repo getAllRepository
}

func NewGetAllUseCase(repository getAllRepository) *GetAllUseCase {
	return &GetAllUseCase{
		repo: repository,
	}
}

func (u *GetAllUseCase) Execute(ctx context.Context) ([]domain.Site, error) {
	return u.repo.GetAll(ctx)
}
