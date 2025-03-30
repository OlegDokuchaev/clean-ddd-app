package di

import (
	"context"
	"errors"
	"fmt"
	"warehouse/internal/infrastructure/logger"
	"warehouse/internal/presentation/commands"

	"go.uber.org/fx"
)

var CommandConsumerModule = fx.Options(
	fx.Provide(
		// Handlers
		fx.Annotate(
			commands.NewHandler,
			fx.As(new(commands.Handler)),
		),

		// Readers
		fx.Annotate(
			commands.NewReader,
			fx.ParamTags(`name:"warehouseCmdReader"`),
			fx.As(new(commands.Reader)),
		),

		// Writers
		fx.Annotate(
			commands.NewWriter,
			fx.ParamTags(`name:"warehouseCmdResWriter"`),
			fx.As(new(commands.Writer)),
		),

		// Processor
		commands.NewProcessor,
	),

	// Lifecycle
	fx.Invoke(setupCommandsLifecycle),
)

func setupCommandsLifecycle(lc fx.Lifecycle, processor *commands.Processor, reader commands.Reader, logger logger.Logger) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.Println("Starting command processor and reader...")

			if err := reader.Start(ctx); err != nil {
				return err
			}
			if err := processor.Start(ctx); err != nil {
				return err
			}

			logger.Println("Command components successfully started")
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Println("Shutting down command processor and reader...")

			var errs []error
			if err := processor.Stop(); err != nil {
				errs = append(errs, fmt.Errorf("processor stop error: %w", err))
			}
			if err := reader.Stop(); err != nil {
				errs = append(errs, fmt.Errorf("reader stop error: %w", err))
			}

			if len(errs) > 0 {
				return errors.Join(errs...)
			}

			logger.Println("Command components successfully stopped")
			return nil
		},
	})
}
