package messaging

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/segmentio/kafka-go"
)

type KafkaProducer struct {
	writer *kafka.Writer
	logger *slog.Logger
}

func NewKafkaProducer(broker, topic string, logger *slog.Logger) *KafkaProducer {
	producer := &KafkaProducer{logger: logger}

	producer.writer = &kafka.Writer{
		Addr:         kafka.TCP(broker),
		Topic:        topic,
		Balancer:     &kafka.Hash{},
		RequiredAcks: kafka.RequireAll,
		MaxAttempts:  3,
		Async:        true,
		Completion: func(messages []kafka.Message, err error) {
			if err != nil {
				logger.Error("kafka publish failed", "error", err, "count", len(messages))
				return
			}
			logger.Debug("kafka publish succeeded", "count", len(messages))
		},
	}

	return producer
}

func (p *KafkaProducer) Publish(ctx context.Context, event SiteCheckEvent) error {
	payload, err := json.Marshal(event)
	if err != nil {
		return err
	}

	return p.writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(event.SiteID.String()),
		Value: payload,
	})
}

func (p *KafkaProducer) Close() error {
	return p.writer.Close()
}
