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

	LockoutMaxFailed int           `envconfig:"LOCKOUT_MAX_FAILED" required:"true"`
	LockoutLockFor   time.Duration `envconfig:"LOCKOUT_LOCK_FOR" required:"true"`

	OtpTTL         time.Duration `envconfig:"OTP_TTL" required:"true"`
	OtpMaxAttempts int           `envconfig:"OTP_MAX_ATTEMPTS" required:"true"`
}

func NewConfig() (*Config, error) {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, fmt.Errorf("failed to load auth config: %w", err)
	}
	return &cfg, nil
}
