package messaging

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/segmentio/kafka-go"
)

type EventHandler func(ctx context.Context, event SiteCheckEvent)

type KafkaConsumer struct {
	reader *kafka.Reader
	logger *slog.Logger
}

func NewKafkaConsumer(broker, groupID, topic string, logger *slog.Logger) *KafkaConsumer {
	return &KafkaConsumer{
		logger: logger,
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers: []string{broker},
			GroupID: groupID,
			Topic:   topic,
		}),
	}
}

func (c *KafkaConsumer) Consume(ctx context.Context, handler EventHandler) error {
	for {
		msg, err := c.reader.FetchMessage(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
				return nil
			}
			return fmt.Errorf("error consume message: %w", err)
		}

		if err := c.handleMessage(msg, handler); err != nil {
			return fmt.Errorf("error handle message: %w", err)
		}
	}
}

func (c *KafkaConsumer) Close() error {
	return c.reader.Close()
}

func (c *KafkaConsumer) handleMessage(msg kafka.Message, handler EventHandler) error {
	processCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var event SiteCheckEvent
	err := json.Unmarshal(msg.Value, &event)
	if err != nil {
		c.logger.Error("failed to unmarshal message", "error", err)
		return c.reader.CommitMessages(processCtx, msg)
	}

	handler(processCtx, event)

	return c.reader.CommitMessages(processCtx, msg)
}
