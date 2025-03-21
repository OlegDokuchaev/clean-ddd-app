package di

import (
	itemDomain "warehouse/internal/domain/item"
	outboxDomain "warehouse/internal/domain/outbox"
	productDomain "warehouse/internal/domain/product"
	itemRepository "warehouse/internal/infrastructure/repository/item"
	outboxRepository "warehouse/internal/infrastructure/repository/outbox"
	productRepository "warehouse/internal/infrastructure/repository/product"

	"go.uber.org/fx"
)

var RepositoryModule = fx.Provide(
	// Item repository
	fx.Annotate(
		itemRepository.New,
		fx.As(new(itemDomain.Repository)),
	),

	// Product repository
	fx.Annotate(
		productRepository.New,
		fx.As(new(productDomain.Repository)),
	),

	// Outbox repository
	fx.Annotate(
		outboxRepository.New,
		fx.As(new(outboxDomain.Repository)),
	),
)
