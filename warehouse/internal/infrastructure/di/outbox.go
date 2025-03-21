package di

import (
	"context"
	"log"
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

func setupOutboxProcessor(lc fx.Lifecycle, processor *outboxProcessor.Processor) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			log.Println("Starting outbox processor...")
			return processor.Start(ctx)
		},
		OnStop: func(ctx context.Context) error {
			log.Println("Stopping outbox processor...")
			return processor.Stop()
		},
	})
}
