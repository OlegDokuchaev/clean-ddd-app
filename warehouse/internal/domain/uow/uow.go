package uow

import (
	"context"
	itemDomain "warehouse/internal/domain/item"
	outboxDomain "warehouse/internal/domain/outbox"
	productDomain "warehouse/internal/domain/product"
)

type UoW interface {
	Product() productDomain.Repository
	Item() itemDomain.Repository
	Outbox() outboxDomain.Repository
	Transaction(ctx context.Context, fn func(u UoW) error) error
}
