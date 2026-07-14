package site

import "time"

type Availability string

const (
	Available   Availability = "available"
	Unavailable Availability = "unavailable"
)

type CheckStatus struct {
	Availability Availability
	Code         int
	CheckedAt    time.Time
	Error        string
}
