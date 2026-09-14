package telegram

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strconv"
	"sync"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gitlab.com/Dokuchaevvn/site-monitor/internal/notification/notifier"
)

type TelegramNotifier struct {
	bot       *tgbotapi.BotAPI
	chatID    string
	logger    *slog.Logger
	mu        sync.Mutex // Нужен, чтобы в конкурентном режиме бот не уходил в 429
	notBefore time.Time  // Время, после которого можно отправлять сообщения дальше (если была 429)
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

	tn.mu.Lock()
	defer tn.mu.Unlock()

	for {
		if err := tn.waitUntil(ctx, tn.notBefore); err != nil {
			return err
		}

		_, err := tn.bot.Send(msg)
		if err == nil {
			return nil
		}

		var tgErr *tgbotapi.Error
		if !errors.As(err, &tgErr) || tgErr.Code != 429 {
			return fmt.Errorf("error sending message to telegram error: %w", err)
		}

		retryAfter := time.Duration(tgErr.ResponseParameters.RetryAfter) * time.Second
		if retryAfter <= 0 {
			retryAfter = time.Second
		}
		tn.notBefore = time.Now().Add(retryAfter)
	}
}

func (tn *TelegramNotifier) waitUntil(ctx context.Context, until time.Time) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	d := time.Until(until)
	if d <= 0 {
		return nil
	}

	t := time.NewTimer(d)
	defer t.Stop()

	select {
	case <-t.C:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
