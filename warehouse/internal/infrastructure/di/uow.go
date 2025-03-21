package di

import (
	"warehouse/internal/domain/uow"
	uowImpl "warehouse/internal/infrastructure/uow"

	"go.uber.org/fx"
)

var UowModule = fx.Provide(
	// UoW
	fx.Annotate(
		uowImpl.New,
		fx.As(new(uow.UoW)),
	),
)
