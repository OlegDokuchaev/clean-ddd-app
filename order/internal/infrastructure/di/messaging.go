package di

import (
	"order/internal/infrastructure/messaging"

	"go.uber.org/fx"
)

// MessagingModule provides configuration and clients for message exchange
var MessagingModule = fx.Options(
	fx.Provide(
		// General Kafka configuration
		messaging.NewConfig,

		// Message readers
		fx.Annotate(
			messaging.NewOrderCommandReader,
			fx.ResultTags(`name:"orderCommandReader"`),
		),
		fx.Annotate(
			messaging.NewWarehouseCommandResultReader,
			fx.ResultTags(`name:"warehouseCommandResultReader"`),
		),
		fx.Annotate(
			messaging.NewCourierCommandResultReader,
			fx.ResultTags(`name:"courierCommandResultReader"`),
		),

		// Message writers
		fx.Annotate(
			messaging.NewOrderCommandWriter,
			fx.ResultTags(`name:"orderCommandWriter"`),
		),
		fx.Annotate(
			messaging.NewWarehouseCommandWriter,
			fx.ResultTags(`name:"warehouseCommandWriter"`),
		),
		fx.Annotate(
			messaging.NewCourierCommandWriter,
			fx.ResultTags(`name:"courierCommandWriter"`),
		),
		fx.Annotate(
			messaging.NewOrderCommandResWriter,
			fx.ResultTags(`name:"orderCommandResWriter"`),
		),
	),
)
