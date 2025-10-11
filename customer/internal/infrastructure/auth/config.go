package auth

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	SigningKey string        `envconfig:"JWT_SIGNING_KEY" required:"true"`
	AccessTTL  time.Duration `envconfig:"JWT_TOKEN_TTL" required:"true"`
	ResetTTL   time.Duration `envconfig:"JWT_RESET_TTL" required:"true"`
}

func NewConfig() (*Config, error) {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, fmt.Errorf("failed to load auth config: %w", err)
	}
	return &cfg, nil
}
