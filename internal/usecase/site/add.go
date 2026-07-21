package site

import (
	"context"

	domain "gitlab.com/Dokuchaevvn/site-monitor/internal/domain/site"
)

type addIfAbsentRepository interface {
	AddIfAbsent(ctx context.Context, site domain.Site) (domain.Site, error)
}

type AddCommand struct {
	URL  string
	Name string
}

type AddUseCase struct {
	repo addIfAbsentRepository
}

func NewAddUseCase(repository addIfAbsentRepository) *AddUseCase {
	return &AddUseCase{
		repo: repository,
	}
}

func (u *AddUseCase) Execute(ctx context.Context, command AddCommand) (domain.Site, error) {
	site, err := domain.NewSite(command.URL, command.Name)
	if err != nil {
		return domain.Site{}, err
	}

	site, err = u.repo.AddIfAbsent(ctx, site)
	if err != nil {
		return domain.Site{}, err
	}

	return site, nil
}
