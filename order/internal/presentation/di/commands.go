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
			processor.Start(ctx)
			log.Println("Command components successfully started")
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Println("Stopping command components...")

			// Stop the processor
			log.Println("Stopping command processor...")
			processor.Stop()

			// Stop the reader
			log.Println("Stopping command reader...")
			reader.Stop()

			log.Println("All command components successfully stopped")
			return nil
		},
	})
}
