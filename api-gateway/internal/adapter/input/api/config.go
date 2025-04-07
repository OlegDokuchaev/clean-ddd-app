package api

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Port string `envconfig:"API_PORT" required:"true"`
}

func NewConfig() (*Config, error) {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, fmt.Errorf("failed to load api config: %w", err)
	}
	return &cfg, nil
}
