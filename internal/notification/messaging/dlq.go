package messaging

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"
)

const (
	ErrorKindRetryExhausted = "retry_exhausted"
	ErrorKindPermanent      = "permanent"
)

type DLQMessage struct {
	Original          json.RawMessage `json:"original"`
	OriginalTopic     string          `json:"original_topic"`
	OriginalPartition int             `json:"original_partition"`
	OriginalOffset    int64           `json:"original_offset"`
	Attempts          int             `json:"attempts"`
	FirstAttemptAt    time.Time       `json:"first_attempt_at"`
	LastAttemptAt     time.Time       `json:"last_attempt_at"`
	Error             string          `json:"error"`
	ErrorKind         string          `json:"error_kind"`
}

type DLQProducer struct {
	writer *kafka.Writer
}

func NewDLQProducer(broker, topic string) *DLQProducer {
	return &DLQProducer{
		writer: &kafka.Writer{
			Addr:         kafka.TCP(broker),
			Topic:        topic,
			Balancer:     &kafka.Hash{},
			RequiredAcks: kafka.RequireAll,
			MaxAttempts:  3,
			Async:        false,
		},
	}
}

func (p *DLQProducer) Publish(ctx context.Context, event DLQMessage, key []byte) error {
	payload, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("marshalling DLQ message: %w", err)
	}

	return p.writer.WriteMessages(ctx, kafka.Message{
		Key:   key,
		Value: payload,
	})
}

func (p *DLQProducer) Close() error {
	return p.writer.Close()
}
