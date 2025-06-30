package di

import (
	courierDomain "courier/internal/domain/courier"
	courierRepository "courier/internal/infrastructure/repository/courier"

	"go.uber.org/fx"
)

var RepositoryModule = fx.Provide(
	// Courier repository
	fx.Annotate(
		courierRepository.New,
		fx.As(new(courierDomain.Repository)),
	),
)
