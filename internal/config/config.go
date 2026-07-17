package config

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

type Site struct {
	URL  string `yaml:"url"`
	Name string `yaml:"name"`
}

type Config struct {
	Sites    []Site        `yaml:"sites"`
	Interval time.Duration `yaml:"interval"`
	HTTPAddr string        `yaml:"http_addr"`
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config file: %w", err)
	}

	var cfg Config

	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parse config file: %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("validate config: %w", err)
	}

	return &cfg, nil
}

func (c *Config) Validate() error {
	if c.Interval <= 0 {
		return errors.New("interval must be greater than zero")
	}

	if strings.TrimSpace(c.HTTPAddr) == "" {
		return errors.New("http_addr is required")
	}

	return nil
}
