package di

import (
	"context"
	"go.uber.org/fx"
	"log"

	"order/internal/presentation/saga/create_order"
)

var SagaConsumerModule = fx.Options(
	fx.Provide(
		// Saga event readers
		fx.Annotate(
			create_order.NewReader,
			fx.ParamTags(`name:"warehouseCommandResultReader"`),
			fx.ResultTags(`name:"warehouseReader"`),
			fx.As(new(create_order.Reader)),
		),
		fx.Annotate(
			create_order.NewReader,
			fx.ParamTags(`name:"courierCommandResultReader"`),
			fx.ResultTags(`name:"courierReader"`),
			fx.As(new(create_order.Reader)),
		),

		// Saga handler and processor
		fx.Annotate(
			create_order.NewHandler,
			fx.As(new(create_order.Handler)),
		),
		fx.Annotate(
			create_order.NewProcessor,
			fx.ParamTags(``, `name:"warehouseReader"`, `name:"courierReader"`),
		),
	),
	fx.Invoke(runProcessor),
)

func runProcessor(in struct {
	fx.In

	Lifecycle       fx.Lifecycle
	Processor       *create_order.Processor
	WarehouseReader create_order.Reader `name:"warehouseReader"`
	CourierReader   create_order.Reader `name:"courierReader"`
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
