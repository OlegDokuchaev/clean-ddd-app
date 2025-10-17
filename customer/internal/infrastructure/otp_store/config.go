package otp_store

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Addr         string        `envconfig:"REDIS_ADDR"`
	Username     string        `envconfig:"REDIS_USERNAME"`
	Password     string        `envconfig:"REDIS_PASSWORD"`
	DB           int           `envconfig:"REDIS_DB"`
	DialTimeout  time.Duration `envconfig:"REDIS_DIAL_TIMEOUT"`
	ReadTimeout  time.Duration `envconfig:"REDIS_READ_TIMEOUT"`
	WriteTimeout time.Duration `envconfig:"REDIS_WRITE_TIMEOUT"`
	KeyPrefix    string        `envconfig:"OTP_KEY_PREFIX"`
}

func NewConfig() (*Config, error) {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, fmt.Errorf("failed to load otp store config: %w", err)
	}
	return &cfg, nil
}
