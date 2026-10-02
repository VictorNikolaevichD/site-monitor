package event

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/ViktorNikolaevichD/site-monitor/internal/notification/messaging"
	"github.com/ViktorNikolaevichD/site-monitor/internal/notification/notifier"
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

func (h *Handler) SiteEvent(ctx context.Context, event messaging.SiteCheckEvent) error {
	if event.IsAvailable {
		return nil
	}

	if err := h.notifier.Send(ctx, notifier.Notification{
		Title: "Site unavailable",
		Text:  fmt.Sprintf("URL %s with error: %s", event.URL, event.ErrorMessage),
	}); err != nil {
		h.logger.Error("Error send message", "error", err)
		if errors.Is(err, notifier.ErrPermanent) {
			return fmt.Errorf("%w: %v", messaging.ErrPermanent, err)
		}
		return err
	}

	return nil
}
