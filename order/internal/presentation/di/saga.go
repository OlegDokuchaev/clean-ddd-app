package di

import (
	"context"
	"errors"
	"fmt"
	"order/internal/infrastructure/logger"
	"order/internal/presentation/saga/create_order"

	"go.uber.org/fx"
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

	Lifecycle fx.Lifecycle
	Logger    logger.Logger

	Processor       *create_order.Processor
	WarehouseReader create_order.Reader `name:"warehouseReader"`
	CourierReader   create_order.Reader `name:"courierReader"`
}) {
	in.Lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			in.Logger.Println("Starting saga readers and processor...")

			if err := in.WarehouseReader.Start(ctx); err != nil {
				return err
			}
			if err := in.CourierReader.Start(ctx); err != nil {
				return err
			}
			if err := in.Processor.Start(ctx); err != nil {
				return err
			}

			in.Logger.Println("Saga components successfully started")
			return nil
		},
		OnStop: func(ctx context.Context) error {
			in.Logger.Println("Stopping saga components...")

			var errs []error
			if err := in.Processor.Stop(); err != nil {
				errs = append(errs, fmt.Errorf("processor stop error: %w", err))
			}
			if err := in.WarehouseReader.Stop(); err != nil {
				errs = append(errs, fmt.Errorf("warehouse reader stop error: %w", err))
			}
			if err := in.CourierReader.Stop(); err != nil {
				errs = append(errs, fmt.Errorf("courier reader stop error: %w", err))
			}

			if len(errs) > 0 {
				return errors.Join(errs...)
			}

			in.Logger.Println("All saga components successfully stopped")
			return nil
		},
	})
}
