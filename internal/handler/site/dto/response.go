package dto

import (
	"time"

	"github.com/google/uuid"

	site "gitlab.com/Dokuchaevvn/site-monitor/internal/domain/site"
)

type SiteResponse struct {
	ID   uuid.UUID `json:"id"`
	URL  string    `json:"url"`
	Name string    `json:"name"`
}

type StatusResponse struct {
	URL          string            `json:"url"`
	Availability site.Availability `json:"availability"`
	Code         *int              `json:"code,omitempty"`
	CheckedAt    *time.Time        `json:"checked_at,omitempty"`
	Duration     *int64            `json:"duration"`
	Error        *string           `json:"error,omitempty"`
}
