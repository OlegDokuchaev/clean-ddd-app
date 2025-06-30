package db

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	URI             string        `envconfig:"DB_URI" required:"true"`
	Database        string        `envconfig:"DB_NAME" required:"true"`
	OrderCollection string        `envconfig:"DB_ORDER_COLLECTION" required:"true"`
	ConnectTimeout  time.Duration `envconfig:"DB_CONNECT_TIMEOUT" required:"true"`
}

func NewConfig() (*Config, error) {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, fmt.Errorf("failed to load db config: %w", err)
	}
	return &cfg, nil
}
