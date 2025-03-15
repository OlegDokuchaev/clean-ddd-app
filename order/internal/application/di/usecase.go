package di

import (
	orderUsecase "order/internal/application/order/usecase"

	"go.uber.org/fx"
)

var UseCaseModule = fx.Provide(
	fx.Annotate(
		orderUsecase.New,
		fx.As(new(orderUsecase.UseCase)),
	),
)
