package di

import (
	"go.uber.org/fx"
	"order/internal/domain/order"
	orderRepository "order/internal/infrastructure/repository/order"
)

var RepositoryModule = fx.Provide(
	// Order repository
	fx.Annotate(
		orderRepository.New,
		fx.As(new(order.Repository)),
	),
)
