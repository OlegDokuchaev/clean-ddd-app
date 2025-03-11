package outbox

import (
	"fmt"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	ProductTopic string `envconfig:"KAFKA_PRODUCT_TOPIC" required:"true"`
}

func NewConfig() (*Config, error) {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, fmt.Errorf("failed to load publisher config: %w", err)
	}
	return &cfg, nil
}
