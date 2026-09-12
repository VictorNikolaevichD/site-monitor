package notifier

import (
	"context"
	"errors"
	"log/slog"
)

type Notifier interface {
	Send(ctx context.Context, n Notification) error
}

type MultiNotifier struct {
	notifiers []Notifier
	logger    *slog.Logger
}

func NewMultiNotifier(logger *slog.Logger, notifiers ...Notifier) *MultiNotifier {
	return &MultiNotifier{
		notifiers: notifiers,
		logger:    logger,
	}
}

func (m *MultiNotifier) Send(ctx context.Context, n Notification) error {
	var errs []error

	for _, ntf := range m.notifiers {
		if err := ntf.Send(ctx, n); err != nil {
			m.logger.Error("failed to send notification", "error", err)
			errs = append(errs, err)
		}
	}

	return errors.Join(errs...)
}
