package dto

import "github.com/ViktorNikolaevichD/site-monitor/internal/notification/health"

func ToHealthResponse(report health.Report) HealthResponse {
	return HealthResponse{
		Status:      report.Status,
		Uptime:      report.Uptime.Milliseconds(),
		CurrentTime: report.CurrentTime,
		Version:     report.Version,
	}
}
