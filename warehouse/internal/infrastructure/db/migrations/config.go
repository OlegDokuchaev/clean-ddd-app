package migrations

import (
	"fmt"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	MigrationsPath string `envconfig:"DB_MIGRATIONS_PATH" required:"true"`
}

func NewConfig() (*Config, error) {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, fmt.Errorf("failed to load migration config: %w", err)
	}
	return &cfg, nil
}
