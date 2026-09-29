package config

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

const testDatabaseURL = "postgres://monitor:monitor@localhost:5432/site_monitor_db?sslmode=disable" //nolint:gosec // G101: test fixture, not a real secret

const validFullYAML = `
interval: 60s
http_addr: :8080
grpc_addr: :9090
log_level: info
http_timeout: 10s
database:
  max_conns: 10
  min_conns: 2
  max_conn_lifetime: 1h
  max_conn_idle_time: 30m
  connect_timeout: 5s
  query_timeout: 3s
kafka:
  broker: localhost:9092
  topic: site-check-events
`

const validMinimalYAML = `
interval: 60s
http_addr: :8080
grpc_addr: :9090
log_level: info
http_timeout: 10s
kafka:
  broker: localhost:9092
  topic: site-check-events
`

const invalidSyntaxYAML = `
interval: 60s
http_addr: [broken
`

func TestLoad_ValidFile(t *testing.T) {
	t.Setenv("DATABASE_URL", testDatabaseURL)
	t.Setenv("KAFKA_BROKER", "localhost:9092")
	t.Setenv("KAFKA_TOPIC", "site-check-events")

	cfg, err := Load(writeConfig(t, validFullYAML))
	if err != nil {
		t.Fatalf("Load() unexpected error: %v", err)
	}

	if cfg.Interval != 60*time.Second {
		t.Fatalf("Interval = %v, want 60s", cfg.Interval)
	}
	if cfg.HTTPAddr != ":8080" {
		t.Fatalf("HTTPAddr = %q, want %q", cfg.HTTPAddr, ":8080")
	}
	if cfg.LogLevel != "info" {
		t.Fatalf("LogLevel = %q, want %q", cfg.LogLevel, "info")
	}
	if cfg.HTTPTimeout != 10*time.Second {
		t.Fatalf("HTTPTimeout = %v, want 10s", cfg.HTTPTimeout)
	}
	if cfg.Database.URL != testDatabaseURL {
		t.Fatalf("Database.URL = %q, want %q", cfg.Database.URL, testDatabaseURL)
	}
	if cfg.Database.MaxConns != 10 {
		t.Fatalf("Database.MaxConns = %d, want 10", cfg.Database.MaxConns)
	}
	if cfg.Database.MinConns != 2 {
		t.Fatalf("Database.MinConns = %d, want 2", cfg.Database.MinConns)
	}
	if cfg.Database.MaxConnLifetime != time.Hour {
		t.Fatalf("Database.MaxConnLifetime = %v, want 1h", cfg.Database.MaxConnLifetime)
	}
	if cfg.Database.MaxConnIdleTime != 30*time.Minute {
		t.Fatalf("Database.MaxConnIdleTime = %v, want 30m", cfg.Database.MaxConnIdleTime)
	}
	if cfg.Database.ConnectTimeout != 5*time.Second {
		t.Fatalf("Database.ConnectTimeout = %v, want 5s", cfg.Database.ConnectTimeout)
	}
	if cfg.Database.QueryTimeout != 3*time.Second {
		t.Fatalf("Database.QueryTimeout = %v, want 3s", cfg.Database.QueryTimeout)
	}
	if cfg.Kafka.Broker != "localhost:9092" {
		t.Fatalf("Kafka.Broker = %q, want %q", cfg.Kafka.Broker, "localhost:9092")
	}
	if cfg.Kafka.Topic != "site-check-events" {
		t.Fatalf("Kafka.Topic = %q, want %q", cfg.Kafka.Topic, "site-check-events")
	}
}

func TestLoad_DefaultValues(t *testing.T) {
	t.Setenv("DATABASE_URL", testDatabaseURL)
	t.Setenv("KAFKA_BROKER", "localhost:9092")
	t.Setenv("KAFKA_TOPIC", "site-check-events")

	cfg, err := Load(writeConfig(t, validMinimalYAML))
	if err != nil {
		t.Fatalf("Load() unexpected error: %v", err)
	}

	if cfg.Interval != 60*time.Second {
		t.Fatalf("Interval = %v, want 60s", cfg.Interval)
	}
	if cfg.HTTPAddr != ":8080" {
		t.Fatalf("HTTPAddr = %q, want %q", cfg.HTTPAddr, ":8080")
	}
	if cfg.LogLevel != "info" {
		t.Fatalf("LogLevel = %q, want %q", cfg.LogLevel, "info")
	}
	if cfg.HTTPTimeout != 10*time.Second {
		t.Fatalf("HTTPTimeout = %v, want 10s", cfg.HTTPTimeout)
	}
	if cfg.Database.URL != testDatabaseURL {
		t.Fatalf("Database.URL = %q, want %q", cfg.Database.URL, testDatabaseURL)
	}
	if cfg.Database.MaxConns != 0 {
		t.Fatalf("Database.MaxConns = %d, want 0 (default)", cfg.Database.MaxConns)
	}
	if cfg.Database.MinConns != 0 {
		t.Fatalf("Database.MinConns = %d, want 0 (default)", cfg.Database.MinConns)
	}
	if cfg.Database.MaxConnLifetime != 0 {
		t.Fatalf("Database.MaxConnLifetime = %v, want 0 (default)", cfg.Database.MaxConnLifetime)
	}
	if cfg.Database.MaxConnIdleTime != 0 {
		t.Fatalf("Database.MaxConnIdleTime = %v, want 0 (default)", cfg.Database.MaxConnIdleTime)
	}
	if cfg.Database.ConnectTimeout != 0 {
		t.Fatalf("Database.ConnectTimeout = %v, want 0 (default)", cfg.Database.ConnectTimeout)
	}
	if cfg.Database.QueryTimeout != 0 {
		t.Fatalf("Database.QueryTimeout = %v, want 0 (default)", cfg.Database.QueryTimeout)
	}
	if cfg.Kafka.Broker != "localhost:9092" {
		t.Fatalf("Kafka.Broker = %q, want %q", cfg.Kafka.Broker, "localhost:9092")
	}
	if cfg.Kafka.Topic != "site-check-events" {
		t.Fatalf("Kafka.Topic = %q, want %q", cfg.Kafka.Topic, "site-check-events")
	}
}

func TestLoad_MissingFile(t *testing.T) {
	t.Setenv("DATABASE_URL", testDatabaseURL)

	_, err := Load(filepath.Join(t.TempDir(), "missing.yaml"))
	if err == nil {
		t.Fatal("Load() error = nil, want error")
	}
	if !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("Load() error = %v, want os.ErrNotExist", err)
	}
}

func TestLoad_InvalidFormat(t *testing.T) {
	t.Setenv("DATABASE_URL", testDatabaseURL)

	_, err := Load(writeConfig(t, invalidSyntaxYAML))
	if err == nil {
		t.Fatal("Load() error = nil, want error")
	}
	if !strings.Contains(err.Error(), "parse config file") {
		t.Fatalf("Load() error = %v, want parse config file", err)
	}
}

func TestLoad_EnvOverrides(t *testing.T) {
	t.Setenv("DATABASE_URL", testDatabaseURL)
	t.Setenv("CHECK_INTERVAL", "30s")
	t.Setenv("LOG_LEVEL", "debug")
	t.Setenv("HTTP_TIMEOUT", "5s")
	t.Setenv("DB_MAX_CONNS", "8")

	cfg, err := Load(writeConfig(t, validMinimalYAML))
	if err != nil {
		t.Fatalf("Load() unexpected error: %v", err)
	}

	if cfg.Interval != 30*time.Second {
		t.Fatalf("Interval = %v, want 30s", cfg.Interval)
	}
	if cfg.LogLevel != "debug" {
		t.Fatalf("LogLevel = %q, want debug", cfg.LogLevel)
	}
	if cfg.HTTPTimeout != 5*time.Second {
		t.Fatalf("HTTPTimeout = %v, want 5s", cfg.HTTPTimeout)
	}
	if cfg.Database.MaxConns != 8 {
		t.Fatalf("Database.MaxConns = %d, want 8", cfg.Database.MaxConns)
	}
}

func TestLoad_EnvPriority(t *testing.T) {
	t.Setenv("DATABASE_URL", testDatabaseURL)
	t.Setenv("CHECK_INTERVAL", "15s")
	t.Setenv("LOG_LEVEL", "error")
	t.Setenv("HTTP_TIMEOUT", "2s")
	t.Setenv("DB_MAX_CONNS", "99")
	t.Setenv("APP_PORT", "9090")
	t.Setenv("GRPC_PORT", "9091")
	t.Setenv("KAFKA_BROKER", "env-localhost:9092")
	t.Setenv("KAFKA_TOPIC", "env-site-check-events")

	cfg, err := Load(writeConfig(t, validFullYAML))
	if err != nil {
		t.Fatalf("Load() unexpected error: %v", err)
	}

	if cfg.Interval != 15*time.Second {
		t.Fatalf("Interval = %v, want 15s (env), not 60s (yaml)", cfg.Interval)
	}
	if cfg.LogLevel != "error" {
		t.Fatalf("LogLevel = %q, want %q (env), not info (yaml)", cfg.LogLevel, "error")
	}
	if cfg.HTTPTimeout != 2*time.Second {
		t.Fatalf("HTTPTimeout = %v, want 2s (env), not 10s (yaml)", cfg.HTTPTimeout)
	}
	if cfg.Database.MaxConns != 99 {
		t.Fatalf("Database.MaxConns = %d, want 99 (env), not 10 (yaml)", cfg.Database.MaxConns)
	}
	if cfg.HTTPAddr != ":9090" {
		t.Fatalf("HTTPAddr = %q, want %q (APP_PORT), not :8080 (yaml)", cfg.HTTPAddr, ":9090")
	}
	if cfg.GRPCAddr != ":9091" {
		t.Fatalf("GRPCAddr = %q, want %q (GRPC_PORT), not :9090 (yaml)", cfg.GRPCAddr, ":9091")
	}
	if cfg.Kafka.Broker != "env-localhost:9092" {
		t.Fatalf("Kafka.Broker = %q, want %q (env), not localhost:9092 (yaml)", cfg.Kafka.Broker, "env-localhost:9092")
	}
	if cfg.Kafka.Topic != "env-site-check-events" {
		t.Fatalf("Kafka.Topic = %q, want %q (env), not site-check-events (yaml)", cfg.Kafka.Topic, "env-site-check-events")
	}
}

func TestValidate_RequiredFields(t *testing.T) {
	testTable := []struct {
		name    string
		mutate  func(cfg *Config)
		wantErr string
	}{
		{
			name: "empty http_addr",
			mutate: func(cfg *Config) {
				cfg.HTTPAddr = ""
			},
			wantErr: "http_addr is required",
		},
		{
			name: "blank http_addr",
			mutate: func(cfg *Config) {
				cfg.HTTPAddr = "  "
			},
			wantErr: "http_addr is required",
		},
		{
			name: "empty grpc_addr",
			mutate: func(cfg *Config) {
				cfg.GRPCAddr = ""
			},
			wantErr: "grpc_addr is required",
		},
		{
			name: "blank grpc_addr",
			mutate: func(cfg *Config) {
				cfg.GRPCAddr = "  "
			},
			wantErr: "grpc_addr is required",
		},
		{
			name: "empty log_level",
			mutate: func(cfg *Config) {
				cfg.LogLevel = ""
			},
			wantErr: "log_level must be one of: debug, info, warn, error",
		},
		{
			name: "empty kafka broker",
			mutate: func(cfg *Config) {
				cfg.Kafka.Broker = ""
			},
			wantErr: "kafka.broker is required",
		},
		{
			name: "blank kafka broker",
			mutate: func(cfg *Config) {
				cfg.Kafka.Broker = "  "
			},
			wantErr: "kafka.broker is required",
		},
		{
			name: "empty kafka topic",
			mutate: func(cfg *Config) {
				cfg.Kafka.Topic = ""
			},
			wantErr: "kafka.topic is required",
		},
		{
			name: "blank kafka topic",
			mutate: func(cfg *Config) {
				cfg.Kafka.Topic = "  "
			},
			wantErr: "kafka.topic is required",
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			cfg := validConfig()
			testCase.mutate(&cfg)

			err := cfg.Validate()
			if err == nil {
				t.Fatal("Validate() error = nil, want error")
			}
			if err.Error() != testCase.wantErr {
				t.Fatalf("Validate() error = %q, want %q", err.Error(), testCase.wantErr)
			}
		})
	}
}

func TestValidate_NumericParams(t *testing.T) {
	testTable := []struct {
		name    string
		mutate  func(cfg *Config)
		wantErr string
	}{
		{
			name: "zero interval",
			mutate: func(cfg *Config) {
				cfg.Interval = 0
			},
			wantErr: "interval must be greater than zero",
		},
		{
			name: "negative interval",
			mutate: func(cfg *Config) {
				cfg.Interval = -time.Second
			},
			wantErr: "interval must be greater than zero",
		},
		{
			name: "zero http timeout",
			mutate: func(cfg *Config) {
				cfg.HTTPTimeout = 0
			},
			wantErr: "http_timeout must be greater than zero",
		},
		{
			name: "negative http timeout",
			mutate: func(cfg *Config) {
				cfg.HTTPTimeout = -time.Second
			},
			wantErr: "http_timeout must be greater than zero",
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			cfg := validConfig()
			testCase.mutate(&cfg)

			err := cfg.Validate()
			if err == nil {
				t.Fatal("Validate() error = nil, want error")
			}
			if err.Error() != testCase.wantErr {
				t.Fatalf("Validate() error = %q, want %q", err.Error(), testCase.wantErr)
			}
		})
	}
}

func writeConfig(t *testing.T, content string) string {
	t.Helper()

	path := filepath.Join(t.TempDir(), "config.yaml")
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("write config: %v", err)
	}
	return path
}

func validConfig() Config {
	return Config{
		Interval:    60 * time.Second,
		HTTPAddr:    ":8080",
		GRPCAddr:    ":9090",
		LogLevel:    "info",
		HTTPTimeout: 10 * time.Second,
		Kafka: Kafka{
			Broker: "localhost:9092",
			Topic:  "site-check-events",
		},
	}
}
