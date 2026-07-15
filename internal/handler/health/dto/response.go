package dto

import (
	"time"

	"gitlab.com/Dokuchaevvn/site-monitor/internal/health"
)

type HealthResponse struct {
	Status      health.Status `json:"status"`
	Uptime      int64         `json:"uptime"`
	CurrentTime time.Time     `json:"time"`
	Version     string        `json:"version"`
}
