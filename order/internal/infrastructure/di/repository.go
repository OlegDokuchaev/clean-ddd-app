package di

import (
	"order/internal/domain/order"
	orderRepository "order/internal/infrastructure/repository/order"

	"go.uber.org/fx"
)

var RepositoryModule = fx.Provide(
	// Order repository
	fx.Annotate(
		orderRepository.New,
		fx.As(new(order.Repository)),
	),
)
