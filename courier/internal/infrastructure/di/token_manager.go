package di

import (
	courierAuth "courier/internal/application/courier/auth"
	"courier/internal/infrastructure/auth"

	"go.uber.org/fx"
)

var TokenManagerModule = fx.Provide(
	// Config
	auth.NewConfig,

	// Token manager
	fx.Annotate(
		auth.NewTokenManager,
		fx.As(new(courierAuth.TokenManager)),
	),
)
