package health

import (
	"context"
	"time"
)

type Status string

const (
	StatusHealthy   Status = "healthy"
	StatusUnhealthy Status = "unhealthy"
)

type dependency interface {
	Name() string
	Check(ctx context.Context) error
}

type Report struct {
	Status      Status
	Uptime      time.Duration
	CurrentTime time.Time
	Version     string
}

type HealthChecker struct {
	version      string
	startedAt    time.Time
	dependencies []dependency
}

func NewHealthChecker(version string, startedAt time.Time, dependencies ...dependency) *HealthChecker {
	return &HealthChecker{
		version:      version,
		startedAt:    startedAt,
		dependencies: dependencies,
	}
}

func (h *HealthChecker) Check(ctx context.Context) Report {
	report := Report{
		Uptime:      time.Since(h.startedAt),
		CurrentTime: time.Now(),
		Version:     h.version,
	}
	for _, d := range h.dependencies {
		err := d.Check(ctx)
		if err != nil {
			report.Status = StatusUnhealthy
			return report
		}
	}

	report.Status = StatusHealthy
	return report
}
