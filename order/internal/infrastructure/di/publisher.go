package di

import (
	createOrder "order/internal/application/order/saga/create_order"
	createOrderPublisher "order/internal/infrastructure/publisher/saga/create_order"

	"go.uber.org/fx"
)

var PublisherModule = fx.Provide(
	// Saga publishers
	fx.Annotate(
		createOrderPublisher.NewPublisher,
		fx.ParamTags(`name:"warehouseCommandWriter"`, `name:"orderCommandWriter"`, `name:"courierCommandWriter"`),
		fx.As(new(createOrder.Publisher)),
	),
)
