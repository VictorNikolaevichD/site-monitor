package event

import (
	"context"
	"log/slog"

	"gitlab.com/Dokuchaevvn/site-monitor/internal/notification/messaging"
)

type Handler struct {
	logger *slog.Logger
}

func NewHandler(logger *slog.Logger) *Handler {
	return &Handler{
		logger: logger,
	}
}

func (h *Handler) SiteEvent(ctx context.Context, event messaging.SiteCheckEvent) {
	if !event.IsAvailable {
		h.logger.Info("site unavailable", "site_id", event.SiteID, "url", event.URL)
	}
}
