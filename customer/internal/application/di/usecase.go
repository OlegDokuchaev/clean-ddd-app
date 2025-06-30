package di

import (
	customerApplication "customer/internal/application/customer"

	"go.uber.org/fx"
)

var UseCaseModule = fx.Provide(
	// Customer auth use case
	fx.Annotate(
		customerApplication.NewAuthUseCase,
		fx.As(new(customerApplication.AuthUseCase)),
	),
)
