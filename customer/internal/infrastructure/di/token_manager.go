package di

import (
	customerApplication "customer/internal/application/customer"
	"customer/internal/infrastructure/auth"

	"go.uber.org/fx"
)

var TokenManagerModule = fx.Provide(
	// Config
	auth.NewConfig,

	// Token manager
	fx.Annotate(
		auth.NewTokenManager,
		fx.As(new(customerApplication.TokenManager)),
	),
)
