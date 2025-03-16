package di

import (
	"context"
	"order/internal/infrastructure/messaging"
	"order/internal/presentation/commands"

	"go.uber.org/fx"
)

var CommandsModule = fx.Provide(
	messaging.NewConfig,
	messaging.NewOrderCommandWriter,
	messaging.NewOrderCommandReader,
	commands.NewHandler,
	commands.NewReader,
	commands.NewWriter,
	commands.NewProcessor,
	fx.Invoke(RunProcessor),
)

func RunProcessor(lc fx.Lifecycle, processor *commands.Processor, reader *commands.ReaderImpl) {
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
