package di

import (
	customerDomain "customer/internal/domain/customer"
	customerRepository "customer/internal/infrastructure/repository/customer"

	"go.uber.org/fx"
)

var RepositoryModule = fx.Provide(
	// Customer repository
	fx.Annotate(
		customerRepository.New,
		fx.As(new(customerDomain.Repository)),
	),
)
