package site

import (
	domain "gitlab.com/Dokuchaevvn/site-monitor/internal/domain/site"
)

type GetAllRepository interface {
	GetAll() []domain.Site
}

type GetAllUseCase struct {
	repo GetAllRepository
}

func NewGetAllUseCase(repository GetAllRepository) *GetAllUseCase {
	return &GetAllUseCase{
		repo: repository,
	}
}

func (u *GetAllUseCase) Execute() []domain.Site {
	return u.repo.GetAll()
}
