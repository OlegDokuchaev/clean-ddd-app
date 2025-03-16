package di

import (
	"go.uber.org/fx"
	createOrder "order/internal/application/order/saga/create_order"
)

var SagaModule = fx.Provide(
	fx.Annotate(
		createOrder.New,
		fx.As(new(createOrder.Saga)),
	),
)
