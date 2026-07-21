package site

import (
	"context"

	"github.com/google/uuid"
)

type deleteByIDRepository interface {
	DeleteByID(ctx context.Context, id uuid.UUID) error
}

type DeleteCommand struct {
	ID uuid.UUID
}

type DeleteUseCase struct {
	repo deleteByIDRepository
}

func NewDeleteUseCase(repository deleteByIDRepository) *DeleteUseCase {
	return &DeleteUseCase{
		repo: repository,
	}
}

func (u *DeleteUseCase) Execute(ctx context.Context, command DeleteCommand) error {
	return u.repo.DeleteByID(ctx, command.ID)
}
