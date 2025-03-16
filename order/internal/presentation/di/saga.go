package di

import (
	"context"
	"log"

	"go.uber.org/fx"

	"order/internal/infrastructure/messaging"
	"order/internal/presentation/saga/create_order"
)

var SagaModule = fx.Provide(
	messaging.NewConfig,

	fx.Annotate(
		messaging.NewWarehouseCommandResultReader,
		fx.ResultTags(`name:"warehouseKafkaReader"`),
	),

	fx.Annotate(
		messaging.NewCourierCommandResultReader,
		fx.ResultTags(`name:"courierKafkaReader"`),
	),

	fx.Annotate(
		create_order.NewReader,
		fx.ParamTags(`name:"warehouseKafkaReader"`),
		fx.ResultTags(`name:"warehouseReader"`),
	),
	fx.Annotate(
		create_order.NewReader,
		fx.ParamTags(`name:"courierKafkaReader"`),
		fx.ResultTags(`name:"courierReader"`),
	),

	create_order.NewHandler,
	fx.Annotate(
		create_order.NewProcessor,
		fx.ParamTags(``, `name:"warehouseReader"`, `name:"courierReader"`),
	),

	fx.Invoke(runProcessor),
)

func runProcessor(in struct {
	fx.In

	Lifecycle       fx.Lifecycle
	Processor       *create_order.Processor
	WarehouseReader *create_order.ReaderImpl `name:"warehouseReader"`
	CourierReader   *create_order.ReaderImpl `name:"courierReader"`
}) {
	in.Lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			log.Println("Starting saga readers and processor...")

			in.WarehouseReader.Start(ctx)
			in.CourierReader.Start(ctx)
			in.Processor.Start(ctx)

			log.Println("Saga components successfully started")
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Println("Stopping saga processor...")

			err := in.Processor.Close()
			if err != nil {
				log.Printf("Error stopping saga processor: %v", err)
				return err
			}

			log.Println("Saga processor successfully stopped")
			return nil
		},
	})
}
