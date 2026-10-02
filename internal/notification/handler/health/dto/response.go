package dto

import (
	"time"

	"github.com/ViktorNikolaevichD/site-monitor/internal/notification/health"
)

type HealthResponse struct {
	Status      health.Status `json:"status" example:"healthy"`
	Uptime      int64         `json:"uptime" example:"729028"`
	CurrentTime time.Time     `json:"time" example:"2026-07-15T13:03:26.9795475+03:00"`
	Version     string        `json:"version" example:"v1.0.0"`
}
