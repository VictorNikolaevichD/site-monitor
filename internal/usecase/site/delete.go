package site

import "github.com/google/uuid"

type DeleteByIDRepository interface {
	DeleteByID(id uuid.UUID) error
}

type DeleteCommand struct {
	ID uuid.UUID
}

type DeleteUseCase struct {
	repo DeleteByIDRepository
}

func NewDeleteUseCase(repository DeleteByIDRepository) *DeleteUseCase {
	return &DeleteUseCase{
		repo: repository,
	}
}

func (u *DeleteUseCase) Execute(command DeleteCommand) error {
	err := u.repo.DeleteByID(command.ID)
	if err != nil {
		return err
	}

	return nil
}
