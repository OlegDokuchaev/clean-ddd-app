package di

import (
	"context"
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
			reader.Start(ctx)
			processor.Start(ctx)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			if err := processor.Close(); err != nil {
				return err
			}
			return nil
		},
	})
}
