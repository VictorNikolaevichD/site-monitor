package site

import "time"

type Availability string

const (
	Available   Availability = "available"
	Unavailable Availability = "unavailable"
	Unknown     Availability = "unknown"
)

type CheckStatus struct {
	Availability Availability
	Code         int
	CheckedAt    time.Time
	Duration     time.Duration
	Error        string
}
