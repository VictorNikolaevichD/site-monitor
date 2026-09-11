package messaging

import (
	"time"

	"github.com/google/uuid"
)

type SiteCheckEvent struct {
	SiteID       uuid.UUID `json:"site_id"`
	URL          string    `json:"url"`
	StatusCode   int       `json:"status_code"`
	IsAvailable  bool      `json:"is_available"`
	ResponseTime int64     `json:"response_time"`
	CheckedAt    time.Time `json:"checked_at"`
	ErrorMessage string    `json:"error_message,omitempty"`
}
