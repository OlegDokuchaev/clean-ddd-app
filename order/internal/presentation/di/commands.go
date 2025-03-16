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
	commands.NewProcessor,
	fx.Invoke(RunProcessor),
)

func RunProcessor(lc fx.Lifecycle, processor *commands.Processor) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go processor.Process(ctx)
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
