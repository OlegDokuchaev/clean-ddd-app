package admin

import (
	"fmt"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Token string `envconfig:"ADMIN_TOKEN" required:"true"`
}

func NewConfig() (*Config, error) {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, fmt.Errorf("failed to load admin config: %w", err)
	}
	return &cfg, nil
}
