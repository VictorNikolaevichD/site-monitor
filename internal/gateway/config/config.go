package config

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

type Config struct {
	HTTPAddr        string `yaml:"http_addr"`
	GRPCAddr        string `yaml:"grpc_addr"`
	LogLevel        string `yaml:"log_level" env:"LOG_LEVEL"`
	GatewayPort     int    `yaml:"-" env:"GATEWAY_PORT"`
	MonitorGRPCAddr string `yaml:"-" env:"MONITOR_GRPC_ADDR"`
}

func Load(path string) (*Config, error) {
	_ = godotenv.Load()

	data, err := os.ReadFile(path) //nolint:gosec // G304: path is the config file
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

func (c *Config) Validate() error {
	if strings.TrimSpace(c.HTTPAddr) == "" {
		return errors.New("http_addr is required")
	}

	if strings.TrimSpace(c.GRPCAddr) == "" {
		return errors.New("grpc_addr is required")
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

func (c *Config) applyEnvOverrides() {
	if c.GatewayPort != 0 {
		c.HTTPAddr = fmt.Sprintf(":%d", c.GatewayPort)
	}
	if strings.TrimSpace(c.MonitorGRPCAddr) != "" {
		c.GRPCAddr = c.MonitorGRPCAddr
	}
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
