package di

import (
	"go.uber.org/fx"
	"order/internal/domain/order"
	"order/internal/infrastructure/db"
	orderRepository "order/internal/infrastructure/repository/order"
)

var RepositoryModule = fx.Provide(
	db.NewConfig,
	db.NewDB,
	fx.Annotate(
		orderRepository.New,
		fx.As(new(order.Repository)),
	),
)
