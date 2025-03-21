package di

import (
	itemApplication "warehouse/internal/application/item"
	productApplication "warehouse/internal/application/product"

	"go.uber.org/fx"
)

var UseCaseModule = fx.Provide(
	// Item use case
	fx.Annotate(
		itemApplication.NewUseCase,
		fx.As(new(itemApplication.UseCase)),
	),

	// Product use case
	fx.Annotate(
		productApplication.NewUseCase,
		fx.As(new(productApplication.UseCase)),
	),
)
