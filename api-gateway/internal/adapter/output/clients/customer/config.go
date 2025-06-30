package customer

import (
	"fmt"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Address        string `envconfig:"CUSTOMER_ADDRESS" required:"true"`
	TimeoutSeconds int    `envconfig:"CUSTOMER_TIMEOUT_SECONDS" required:"true"`
}

func NewConfig() (*Config, error) {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, fmt.Errorf("failed to load client config: %w", err)
	}
	return &cfg, nil
}
