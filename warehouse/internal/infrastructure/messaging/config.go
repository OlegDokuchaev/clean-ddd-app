package messaging

import (
	"fmt"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Brokers  string `envconfig:"KAFKA_BROKERS" required:"true"`
	RetryMax int    `envconfig:"KAFKA_RETRY_MAX" required:"true"`
	Timeout  int    `envconfig:"KAFKA_TIMEOUT" required:"true"`
}

func NewConfig() (*Config, error) {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, fmt.Errorf("failed to load kafka config: %w", err)
	}
	return &cfg, nil
}
