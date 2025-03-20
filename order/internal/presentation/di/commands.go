package di

import (
	"context"
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

			reader.Start(ctx)
			if err := processor.Start(ctx); err != nil {
				return err
			}

			log.Println("Command components successfully started")
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Println("Stopping command components...")

			err := processor.Stop()
			reader.Stop()

			if err == nil {
				log.Println("All saga components successfully stopped")
			}
			return nil
		},
	})
}
