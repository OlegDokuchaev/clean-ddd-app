package di

import (
	"context"
	"warehouse/internal/infrastructure/logger"
	outboxProcessor "warehouse/internal/infrastructure/outbox"

	"go.uber.org/fx"
)

var OutboxProcessorModule = fx.Options(
	fx.Provide(
		// Outbox processor
		outboxProcessor.NewProcessor,
	),

	// Lifecycle
	fx.Invoke(setupOutboxProcessor),
)

func setupOutboxProcessor(lc fx.Lifecycle, processor *outboxProcessor.Processor, logger logger.Logger) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.Println("Starting outbox processor...")
			return processor.Start(ctx)
		},
		OnStop: func(ctx context.Context) error {
			logger.Println("Stopping outbox processor...")
			return processor.Stop()
		},
	})
}
