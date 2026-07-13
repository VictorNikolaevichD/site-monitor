package site

import (
	neturl "net/url"
	"strings"

	"github.com/google/uuid"
)

type Site struct {
	ID   uuid.UUID
	URL  string
	Name string
}

func NewSite(url string, name string) (Site, error) {
	if strings.TrimSpace(url) == "" {
		return Site{}, ErrURLRequired
	}

	u, err := neturl.Parse(url)
	if err != nil {
		return Site{}, ErrInvalidURL
	}

	if u.Scheme != "http" && u.Scheme != "https" {
		return Site{}, ErrInvalidURL
	}

	if u.Host == "" {
		return Site{}, ErrInvalidURL
	}

	if strings.TrimSpace(name) == "" {
		return Site{}, ErrNameRequired
	}

	return Site{
		ID:   uuid.New(),
		URL:  u.String(),
		Name: name,
	}, nil
}
