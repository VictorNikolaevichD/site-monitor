package dto

import (
	"time"

	"github.com/google/uuid"

	site "gitlab.com/Dokuchaevvn/site-monitor/internal/domain/site"
)

type SiteResponse struct {
	ID   uuid.UUID `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	URL  string    `json:"url" example:"https://example.com"`
	Name string    `json:"name" example:"Example"`
}

type StatusResponse struct {
	URL          string            `json:"url" example:"https://example.com"`
	Availability site.Availability `json:"availability" example:"available"`
	Code         *int              `json:"code,omitempty" example:"200"`
	CheckedAt    *time.Time        `json:"checked_at,omitempty" example:"2026-07-15T10:00:00Z"`
	Duration     *int64            `json:"duration,omitempty" example:"245"`
	Error        *string           `json:"error,omitempty" example:"context deadline exceeded"`
}

type CheckResultResponse struct {
	ID           int64             `json:"id" example:"1"`
	Availability site.Availability `json:"availability" example:"available"`
	Code         int               `json:"code" example:"200"`
	Duration     int64             `json:"duration" example:"245"`
	Error        string            `json:"error" example:""`
	CheckedAt    time.Time         `json:"checked_at" example:"2026-07-15T10:00:00Z"`
}

type PaginationMeta struct {
	Total  int `json:"total" example:"42"`
	Limit  int `json:"limit" example:"20"`
	Offset int `json:"offset" example:"0"`
}

type CheckHistoryResponse struct {
	Result     []CheckResultResponse `json:"result"`
	Pagination PaginationMeta        `json:"pagination"`
}
