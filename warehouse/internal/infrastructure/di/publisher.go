package di

import (
	"warehouse/internal/domain/outbox"
	outboxPublisher "warehouse/internal/infrastructure/publisher/outbox"

	"go.uber.org/fx"
)

var PublisherModule = fx.Provide(
	// Outbox publisher
	fx.Annotate(
		outboxPublisher.NewPublisher,
		fx.ParamTags(`name:"productEventWriter"`),
		fx.As(new(outbox.Publisher)),
	),
)
