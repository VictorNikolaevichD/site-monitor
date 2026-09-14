package telegram

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strconv"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gitlab.com/Dokuchaevvn/site-monitor/internal/notification/notifier"
)

type TelegramNotifier struct {
	bot    *tgbotapi.BotAPI
	chatID string
	logger *slog.Logger
}

func NewTelegramNotifier(bot *tgbotapi.BotAPI, chatID string, logger *slog.Logger) *TelegramNotifier {
	return &TelegramNotifier{
		bot:    bot,
		chatID: chatID,
		logger: logger,
	}
}

func (tn *TelegramNotifier) Send(ctx context.Context, n notifier.Notification) error {
	chatID, err := strconv.ParseInt(tn.chatID, 10, 64)
	if err != nil {
		return fmt.Errorf("error parsing chatID error: %w", err)
	}
	msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("%s\n%s", n.Title, n.Text))

	_, err = tn.bot.Send(msg)
	if err == nil {
		return nil
	}

	var tgErr *tgbotapi.Error

	if !errors.As(err, &tgErr) || tgErr.Code != 429 {
		return fmt.Errorf("error sending message to telegram error: %w", err)
	}

	retryAfter := tgErr.ResponseParameters.RetryAfter

	select {
	case <-time.After(time.Duration(retryAfter) * time.Second):
	case <-ctx.Done():
		return ctx.Err()
	}

	_, err = tn.bot.Send(msg)
	if err != nil {
		return fmt.Errorf("error sending message to telegram after retry: %w", err)
	}

	return nil
}
