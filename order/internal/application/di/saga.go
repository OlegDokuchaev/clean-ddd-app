package di

import (
	createOrder "order/internal/application/order/saga/create_order"

	"go.uber.org/fx"
)

var SagaModule = fx.Provide(
	fx.Annotate(
		createOrder.New,
		fx.As(new(createOrder.Saga)),
	),
	fx.Annotate(
		createOrder.NewManager,
		fx.As(new(createOrder.Manager)),
	),
)
