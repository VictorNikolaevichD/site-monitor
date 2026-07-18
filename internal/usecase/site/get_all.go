package site

import (
	domain "gitlab.com/Dokuchaevvn/site-monitor/internal/domain/site"
)

type getAllRepository interface {
	GetAll() []domain.Site
}

type GetAllUseCase struct {
	repo getAllRepository
}

func NewGetAllUseCase(repository getAllRepository) *GetAllUseCase {
	return &GetAllUseCase{
		repo: repository,
	}
}

func (u *GetAllUseCase) Execute() []domain.Site {
	return u.repo.GetAll()
}
