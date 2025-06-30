package di

import (
	courierApplication "courier/internal/application/courier"
	courierAuth "courier/internal/application/courier/auth"

	"go.uber.org/fx"
)

var UseCaseModule = fx.Provide(
	// Courier use case
	fx.Annotate(
		courierApplication.NewUseCase,
		fx.As(new(courierApplication.UseCase)),
	),

	// Courier auth use case
	fx.Annotate(
		courierAuth.NewUseCase,
		fx.As(new(courierAuth.UseCase)),
	),
)
