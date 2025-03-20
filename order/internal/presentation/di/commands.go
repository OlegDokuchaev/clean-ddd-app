package di

import (
	"context"
	"errors"
	"fmt"
	"log"
	"order/internal/presentation/commands"

	"go.uber.org/fx"
)

var CommandConsumerModule = fx.Options(
	fx.Provide(
		fx.Annotate(
			commands.NewHandler,
			fx.As(new(commands.Handler)),
		),
		fx.Annotate(
			commands.NewReader,
			fx.ParamTags(`name:"orderCommandReader"`),
			fx.As(new(commands.Reader)),
		),
		fx.Annotate(
			commands.NewWriter,
			fx.ParamTags(`name:"orderCommandResWriter"`),
			fx.As(new(commands.Writer)),
		),
		commands.NewProcessor,
	),
	fx.Invoke(setupCommandsLifecycle),
)

func setupCommandsLifecycle(lc fx.Lifecycle, processor *commands.Processor, reader commands.Reader) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			log.Println("Starting command processor and reader...")

			if err := reader.Start(ctx); err != nil {
				return err
			}
			if err := processor.Start(ctx); err != nil {
				return err
			}

			log.Println("Command components successfully started")
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Println("Stopping command components...")

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

			log.Println("All command components successfully stopped")
			return nil
		},
	})
}
