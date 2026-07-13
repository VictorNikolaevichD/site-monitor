package site

import (
	domain "gitlab.com/Dokuchaevvn/site-monitor/internal/domain/site"
)

type Repository interface {
	AddIfAbsent(site domain.Site) (domain.Site, error)
}

type AddCommand struct {
	URL  string
	Name string
}

type AddUseCase struct {
	repo Repository
}

func NewAddUseCase(repository Repository) *AddUseCase {
	return &AddUseCase{
		repo: repository,
	}
}

func (u *AddUseCase) Execute(command AddCommand) (domain.Site, error) {
	site, err := domain.NewSite(command.URL, command.Name)
	if err != nil {
		return domain.Site{}, err
	}

	site, err = u.repo.AddIfAbsent(site)
	if err != nil {
		return domain.Site{}, err
	}

	return site, nil
}
