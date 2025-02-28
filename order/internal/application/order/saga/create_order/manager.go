package create_order

import (
	"context"
	orderDomain "order/internal/domain/order"
)

type Manager interface {
	Create(ctx context.Context, order *orderDomain.Order)
}
