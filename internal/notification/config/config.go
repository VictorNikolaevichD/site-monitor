package config

import (
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Kafka struct {
	Broker   string `env:"KAFKA_BROKER,required,notEmpty"`
	Topic    string `env:"KAFKA_TOPIC,required,notEmpty"`
	GroupID  string `env:"KAFKA_GROUP_ID,required,notEmpty"`
	DLQTopic string `env:"KAFKA_DLQ_TOPIC" envDefault:"site-check-events-dlq"`
}

type Retry struct {
	Enabled         bool          `env:"NOTIFICATION_RETRY_ENABLED" envDefault:"true"`
	MaxAttempts     int           `env:"NOTIFICATION_RETRY_MAX_ATTEMPTS" envDefault:"3"`
	InitialInterval time.Duration `env:"NOTIFICATION_RETRY_INTERVAL" envDefault:"1s"`
	Multiplier      float64       `env:"NOTIFICATION_RETRY_MULTIPLIER" envDefault:"2"`
}

type Telegram struct {
	BotToken string `env:"NOTIFICATION_TG_BOT_TOKEN"`
	ChatID   string `env:"NOTIFICATION_TG_CHAT_ID"`
}

type Config struct {
	HTTPAddr string `env:"NOTIFICATION_HTTP_ADDR" envDefault:":8081"`
	LogLevel string `env:"NOTIFICATION_LOG_LEVEL" envDefault:"info"`
	Kafka    Kafka
	Retry    Retry
	Telegram Telegram
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	var cfg Config

	if err := env.Parse(&cfg); err != nil {
		return nil, fmt.Errorf("parse env config: %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("validate config: %w", err)
	}

	return &cfg, nil
}

func (c *Config) Validate() error {
	if strings.TrimSpace(c.Telegram.BotToken) == "" {
		return errors.New("NOTIFICATION_TG_BOT_TOKEN is required")
	}

	if strings.TrimSpace(c.Telegram.ChatID) == "" {
		return errors.New("NOTIFICATION_TG_CHAT_ID is required")
	}

	if c.Retry.MaxAttempts < 1 {
		return errors.New("NOTIFICATION_RETRY_MAX_ATTEMPTS must be greater than zero")
	}

	if c.Retry.InitialInterval <= 0 {
		return errors.New("NOTIFICATION_RETRY_INTERVAL must be greater than zero")
	}

	if c.Retry.Multiplier < 1 {
		return errors.New("NOTIFICATION_RETRY_MULTIPLIER must be greater than or equal to 1")
	}

	if strings.TrimSpace(c.Kafka.DLQTopic) == "" {
		return errors.New("KAFKA_DLQ_TOPIC is required")
	}

	if _, err := parseLogLevel(c.LogLevel); err != nil {
		return err
	}

	return nil
}

func (c *Config) SlogLevel() slog.Level {
	level, err := parseLogLevel(c.LogLevel)
	if err != nil {
		return slog.LevelInfo
	}
	return level
}

func parseLogLevel(value string) (slog.Level, error) {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "debug":
		return slog.LevelDebug, nil
	case "info":
		return slog.LevelInfo, nil
	case "warn":
		return slog.LevelWarn, nil
	case "error":
		return slog.LevelError, nil
	default:
		return 0, fmt.Errorf("log_level must be one of: debug, info, warn, error")
	}
}
