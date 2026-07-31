package site

import "time"

type CheckResult struct {
	ID           int64
	Availability Availability
	Code         int
	CheckedAt    time.Time
	Duration     time.Duration
	Error        string
}
