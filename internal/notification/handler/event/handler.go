package event

import (
	"context"
	"fmt"
	"log/slog"

	"gitlab.com/Dokuchaevvn/site-monitor/internal/notification/messaging"
	"gitlab.com/Dokuchaevvn/site-monitor/internal/notification/notifier"
)

type Notifier interface {
	Send(ctx context.Context, n notifier.Notification) error
}

type Handler struct {
	notifier Notifier
	logger   *slog.Logger
}

func NewHandler(notifier Notifier, logger *slog.Logger) *Handler {
	return &Handler{
		notifier: notifier,
		logger:   logger,
	}
}

func (h *Handler) SiteEvent(ctx context.Context, event messaging.SiteCheckEvent) {
	if event.IsAvailable {
		return
	}

	// TODO: добавить ещё в текст недостоющие данные из event, которые сейчас не попадают
	if err := h.notifier.Send(ctx, notifier.Notification{
		Title: "Site unavailable",
		Text:  fmt.Sprintf("URL %s with error: %s", event.URL, event.ErrorMessage),
	}); err != nil {
		h.logger.Error("Error send message", "error", err)
	}
}
