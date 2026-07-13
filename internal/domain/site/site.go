package site

import (
	neturl "net/url"

	"github.com/google/uuid"
	"gitlab.com/Dokuchaevvn/site-monitor/internal/errors"
)

type Site struct {
	ID   uuid.UUID
	URL  string
	Name string
}

func NewSite(url string, name string) (Site, error) {
	if url == "" {
		return Site{}, errors.ErrURLRequired
	}

	u, err := neturl.Parse(url)
	if err != nil {
		return Site{}, errors.ErrInvalidURL
	}

	if u.Scheme != "http" && u.Scheme != "https" {
		return Site{}, errors.ErrInvalidURL
	}

	if u.Host == "" {
		return Site{}, errors.ErrInvalidURL
	}

	if name == "" {
		return Site{}, errors.ErrNameRequired
	}

	return Site{
		ID:   uuid.New(),
		URL:  u.String(),
		Name: name,
	}, nil
}
