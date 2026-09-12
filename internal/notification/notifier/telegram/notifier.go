package telegram

import (
	"context"
	"errors"
	"log/slog"

	"gitlab.com/Dokuchaevvn/site-monitor/internal/notification/notifier"
)

type TelegramNotifier struct {
	botToken string
	chatID   string
	logger   *slog.Logger
}

func NewTelegramNotifier(botToken, chatID string, logger *slog.Logger) *TelegramNotifier {
	return &TelegramNotifier{
		botToken: botToken,
		chatID:   chatID,
		logger:   logger,
	}
}

func (tn *TelegramNotifier) Send(ctx context.Context, n notifier.Notification) error {
	// TODO: добавить реализацию отправки через бота
	return errors.New("need realisation")
}
