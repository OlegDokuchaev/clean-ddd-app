package create_order

import (
	"context"
	orderDomain "order/internal/domain/order"
)

type ManagerImpl struct {
	publisher Publisher
}

func NewManager(publisher Publisher) Manager {
	return &ManagerImpl{publisher: publisher}
}

func (m *ManagerImpl) Create(ctx context.Context, order *orderDomain.Order) {
	orderItems := domainItemsToOrderItems(order.Items)

	cmd := ReserveItemsCmd{
		OrderID: order.ID,
		Items:   orderItems,
	}
	_ = m.publisher.PublishReserveItemsCmd(ctx, cmd)
}

var _ Manager = (*ManagerImpl)(nil)
