package di

import (
	"context"
	"errors"
	"fmt"
	"warehouse/internal/infrastructure/logger"
	"warehouse/internal/presentation/events"

	"go.uber.org/fx"
)

var EventsModule = fx.Options(
	fx.Provide(
		// Handlers
		fx.Annotate(
			events.NewHandler,
			fx.As(new(events.Handler)),
		),

		// Readers
		fx.Annotate(
			events.NewReader,
			fx.ParamTags(`name:"productEventReader"`),
			fx.ResultTags(`name:"productReader"`),
			fx.As(new(events.Reader)),
		),

		// Processors
		fx.Annotate(
			events.NewProcessor,
			fx.ParamTags(``, `name:"productReader"`),
		),
	),

	// Lifecycle
	fx.Invoke(setupEventsLifecycle),
)

func setupEventsLifecycle(in struct {
	fx.In

	Lifecycle fx.Lifecycle
	Logger    logger.Logger

	Processor     *events.Processor
	ProductReader events.Reader `name:"productReader"`
}) {
	in.Lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			in.Logger.Println("Starting product event processor...")

			if err := in.ProductReader.Start(ctx); err != nil {
				return err
			}
			if err := in.Processor.Start(ctx); err != nil {
				return err
			}

			in.Logger.Println("Product event processor started successfully")
			return nil
		},
		OnStop: func(ctx context.Context) error {
			in.Logger.Println("Stopping product event processor...")

			var errs []error
			if err := in.Processor.Stop(); err != nil {
				errs = append(errs, fmt.Errorf("processor stop error: %w", err))
			}
			if err := in.ProductReader.Stop(); err != nil {
				errs = append(errs, fmt.Errorf("product reader stop error: %w", err))
			}

			if len(errs) > 0 {
				return errors.Join(errs...)
			}

			in.Logger.Println("Product event processor stopped successfully")
			return nil
		},
	})
}
