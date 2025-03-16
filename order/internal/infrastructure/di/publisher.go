package di

import (
	"order/internal/infrastructure/messaging"
	createOrderPublisher "order/internal/infrastructure/publisher/saga/create_order"

	"go.uber.org/fx"
)

var PublisherModule = fx.Provide(
	messaging.NewConfig,

	fx.Annotate(
		messaging.NewOrderCommandWriter,
		fx.ResultTags(`name:"orderWriter"`),
	),
	fx.Annotate(
		messaging.NewWarehouseCommandWriter,
		fx.ResultTags(`name:"warehouseWriter"`),
	),
	fx.Annotate(
		messaging.NewCourierCommandWriter,
		fx.ResultTags(`name:"courierWriter"`),
	),

	fx.Annotate(
		createOrderPublisher.NewPublisher,
		fx.ParamTags(`name:"warehouseWriter"`, `name:"orderWriter"`, `name:"courierWriter"`),
	),
)
