package site

import (
	"time"

	"github.com/google/uuid"
)

type SiteRecord struct {
	ID        uuid.UUID `db:"id"`
	URL       string    `db:"url"`
	Name      string    `db:"name"`
	IsActive  bool      `db:"is_active"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type CheckStatusRecord struct {
	ID           int64         `db:"id"`
	SiteID       uuid.UUID     `db:"site_id"`
	HTTPCode     int           `db:"http_code"`
	Duration     time.Duration `db:"duration_ns"`
	Availability bool          `db:"availability"`
	Error        string        `db:"error"`
	CheckedAt    time.Time     `db:"checked_at"`
}
