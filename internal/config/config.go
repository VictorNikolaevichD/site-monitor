package config

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

type Database struct {
	URL             string        `env:"DATABASE_URL,required,notEmpty"`
	MaxConns        int32         `yaml:"max_conns" env:"DB_MAX_CONNS"`
	MinConns        int32         `yaml:"min_conns" env:"DB_MIN_CONNS"`
	MaxConnLifetime time.Duration `yaml:"max_conn_lifetime" env:"DB_MAX_CONN_LIFETIME"`
	MaxConnIdleTime time.Duration `yaml:"max_conn_idle_time" env:"DB_MAX_CONN_IDLE_TIME"`
	ConnectTimeout  time.Duration `yaml:"connect_timeout" env:"DB_CONNECT_TIMEOUT"`
	QueryTimeout    time.Duration `yaml:"query_timeout" env:"DB_QUERY_TIMEOUT"`
}

type Kafka struct {
	Broker string `yaml:"broker" env:"KAFKA_BROKER"`
	Topic  string `yaml:"topic" env:"KAFKA_TOPIC"`
}

type Config struct {
	Interval    time.Duration `yaml:"interval" env:"CHECK_INTERVAL"`
	HTTPAddr    string        `yaml:"http_addr"`
	GRPCAddr    string        `yaml:"grpc_addr"`
	LogLevel    string        `yaml:"log_level" env:"LOG_LEVEL"`
	HTTPTimeout time.Duration `yaml:"http_timeout" env:"HTTP_TIMEOUT"`
	AppPort     int           `yaml:"-" env:"APP_PORT"`
	GRPCPort    int           `yaml:"-" env:"GRPC_PORT"`
	Database    Database      `yaml:"database"`
	Kafka       Kafka         `yaml:"kafka"`
}

func Load(path string) (*Config, error) {
	_ = godotenv.Load()

	data, err := os.ReadFile(path) //nolint:gosec // G304: path is the config file from the CLI
	if err != nil {
		return nil, fmt.Errorf("read config file: %w", err)
	}

	var cfg Config

	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parse config file: %w", err)
	}

	if err := env.Parse(&cfg); err != nil {
		return nil, fmt.Errorf("parse env config: %w", err)
	}

	cfg.applyEnvOverrides()

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("validate config: %w", err)
	}

	return &cfg, nil
}

func (c *Config) applyEnvOverrides() {
	if c.AppPort != 0 {
		c.HTTPAddr = fmt.Sprintf(":%d", c.AppPort)
	}

	if c.GRPCPort != 0 {
		c.GRPCAddr = fmt.Sprintf(":%d", c.GRPCPort)
	}
}

func (c *Config) Validate() error {
	if c.Interval <= 0 {
		return errors.New("interval must be greater than zero")
	}

	if strings.TrimSpace(c.HTTPAddr) == "" {
		return errors.New("http_addr is required")
	}

	if strings.TrimSpace(c.GRPCAddr) == "" {
		return errors.New("grpc_addr is required")
	}

	if c.HTTPTimeout <= 0 {
		return errors.New("http_timeout must be greater than zero")
	}

	if _, err := parseLogLevel(c.LogLevel); err != nil {
		return err
	}

	if strings.TrimSpace(c.Kafka.Broker) == "" {
		return errors.New("kafka.broker is required")
	}

	if strings.TrimSpace(c.Kafka.Topic) == "" {
		return errors.New("kafka.topic is required")
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
