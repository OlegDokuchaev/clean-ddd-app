package testutils

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type Mode string

const (
	ModeContainer Mode = "container"
	ModeReal      Mode = "real"
)

type Config struct {
	Mode Mode `envconfig:"TEST_MODE" default:"container"`

	ImapHost     string `envconfig:"IMAP_HOST"`
	ImapPort     string `envconfig:"IMAP_PORT"`
	ImapUsername string `envconfig:"IMAP_USERNAME"`
	ImapPassword string `envconfig:"IMAP_PASSWORD"`
}

func NewConfig() (*Config, error) {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, fmt.Errorf("failed to load test config: %w", err)
	}
	return &cfg, nil
}
