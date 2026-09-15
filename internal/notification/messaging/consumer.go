package messaging

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/segmentio/kafka-go"

	"gitlab.com/Dokuchaevvn/site-monitor/internal/notification/config"
)

type EventHandler func(ctx context.Context, event SiteCheckEvent) error

type KafkaConsumer struct {
	reader *kafka.Reader
	dlq    *DLQProducer
	topic  string
	retry  config.Retry
	logger *slog.Logger
}

func NewKafkaConsumer(
	broker, groupID, topic, dlqTopic string,
	retry config.Retry,
	logger *slog.Logger,
) *KafkaConsumer {
	return &KafkaConsumer{
		retry:  retry,
		topic:  topic,
		logger: logger,
		dlq:    NewDLQProducer(broker, dlqTopic),
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

		maxAttempts := c.retry.MaxAttempts
		if !c.retry.Enabled {
			maxAttempts = 1
		}

		backoff := c.retry.InitialInterval
		firstAttemptAt := time.Now()

		for attempt := 1; attempt <= maxAttempts; attempt++ {
			err := c.handleMessage(msg, handler)
			if err == nil {
				break
			}

			c.logger.Error("failed to handle message",
				"attempt", attempt,
				"error", err,
				"max_attempts", maxAttempts,
			)

			permanent := errors.Is(err, ErrPermanent)
			if permanent || attempt == maxAttempts {
				kind := ErrorKindRetryExhausted
				if permanent {
					kind = ErrorKindPermanent
				}

				if err := c.publishToDLQ(ctx, msg, err, attempt, firstAttemptAt, kind); err != nil {
					return fmt.Errorf("error publish dlq: %w", err)
				}

				if err := c.reader.CommitMessages(ctx, msg); err != nil {
					return fmt.Errorf("error commit after dlq: %w", err)
				}

				break
			}

			c.logger.Info("retrying message",
				"attempt", attempt+1,
				"max_attempts", maxAttempts,
				"backoff", backoff,
			)

			if err := c.waitBackoff(ctx, backoff); err != nil {
				return nil
			}

			backoff = time.Duration(float64(backoff) * c.retry.Multiplier)
		}
	}
}

func (c *KafkaConsumer) Close() error {
	return errors.Join(c.dlq.Close(), c.reader.Close())
}

func (c *KafkaConsumer) publishToDLQ(
	ctx context.Context,
	msg kafka.Message,
	handleErr error,
	attempts int,
	firstAttemptAt time.Time,
	kind string,
) error {
	c.logger.Error("sending message to dlq",
		"error", handleErr,
		"error_kind", kind,
		"attempts", attempts,
		"offset", msg.Offset,
		"topic", c.topic,
	)

	return c.dlq.Publish(ctx, DLQMessage{
		Original:          json.RawMessage(msg.Value),
		OriginalTopic:     c.topic,
		OriginalPartition: msg.Partition,
		OriginalOffset:    msg.Offset,
		Attempts:          attempts,
		FirstAttemptAt:    firstAttemptAt,
		LastAttemptAt:     time.Now(),
		Error:             handleErr.Error(),
		ErrorKind:         kind,
	}, msg.Key)
}

func (c *KafkaConsumer) waitBackoff(ctx context.Context, d time.Duration) error {
	if err := ctx.Err(); err != nil {
		return err
	}
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

func (c *KafkaConsumer) handleMessage(msg kafka.Message, handler EventHandler) error {
	processCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var event SiteCheckEvent
	err := json.Unmarshal(msg.Value, &event)
	if err != nil {
		c.logger.Error("failed to unmarshal message", "error", err)
		return fmt.Errorf("%w: unmarshal message: %v", ErrPermanent, err)
	}

	err = handler(processCtx, event)
	if err != nil {
		c.logger.Error("failed to handle message", "error", err)
		return err
	}

	return c.reader.CommitMessages(processCtx, msg)
}
