package mail_sender

import (
	"fmt"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Host     string `envconfig:"MAIL_HOST"`
	Port     int    `envconfig:"MAIL_PORT"`
	Username string `envconfig:"MAIL_USERNAME"`
	Password string `envconfig:"MAIL_PASSWORD"`
	FromName string `envconfig:"MAIL_FROM_NAME"`
	FromAddr string `envconfig:"MAIL_FROM_ADDRESS"`
}

func NewConfig() (*Config, error) {
	var cfg Config
	if err := envconfig.Process("MAIL", &cfg); err != nil {
		return nil, fmt.Errorf("failed to load mail sender config: %w", err)
	}
	return &cfg, nil
}
