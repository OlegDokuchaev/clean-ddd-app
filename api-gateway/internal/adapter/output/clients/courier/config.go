package courier

import (
	"fmt"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Address        string `envconfig:"COURIER_ADDRESS" required:"true"`
	TimeoutSeconds int    `envconfig:"COURIER_TIMEOUT_SECONDS" required:"true"`
}

func NewConfig() (*Config, error) {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, fmt.Errorf("failed to load client config: %w", err)
	}
	return &cfg, nil
}
