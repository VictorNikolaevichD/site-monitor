package site

import (
	"github.com/google/uuid"
	domain "gitlab.com/Dokuchaevvn/site-monitor/internal/domain/site"
)

type getByIDRepository interface {
	GetByID(id uuid.UUID) (domain.Site, error)
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

func (u *GetStatusUseCase) Execute(command GetStatusCommand) (domain.Site, error) {
	site, err := u.repo.GetByID(command.ID)
	if err != nil {
		return domain.Site{}, err
	}

	return site, nil
}
