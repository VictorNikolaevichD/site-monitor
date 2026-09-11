package messaging

import (
	"log/slog"

	"github.com/segmentio/kafka-go"
)

type KafkaConsumer struct {
	reader *kafka.Reader
	logger *slog.Logger
}

func NewKafkaConsumer(broker, groupID, topic string, logger *slog.Logger) *KafkaConsumer {
	return &KafkaConsumer{
		logger: logger,
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers:     []string{broker},
			GroupID:     groupID,
			Topic: topic,
		}),
	}
}

func (c *KafkaConsumer) Close() error {
	return c.reader.Close()
}
