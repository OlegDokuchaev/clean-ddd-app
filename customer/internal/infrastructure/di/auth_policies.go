package di

import (
	customerAplication "customer/internal/application/customer"
	customerDomain "customer/internal/domain/customer"
	"customer/internal/infrastructure/auth"

	"go.uber.org/fx"
)

var AuthPoliciesModule = fx.Provide(
	NewLockoutPolicy,
	NewOtpPolicy,
)

func NewLockoutPolicy(cfg *auth.Config) customerDomain.LockoutPolicy {
	return customerDomain.LockoutPolicy{
		MaxFailed: cfg.LockoutMaxFailed,
		LockFor:   cfg.LockoutLockFor,
	}
}

func NewOtpPolicy(cfg *auth.Config) customerAplication.OtpPolicy {
	return customerAplication.OtpPolicy{
		TTL:         cfg.OtpTTL,
		MaxAttempts: cfg.OtpMaxAttempts,
	}
}
