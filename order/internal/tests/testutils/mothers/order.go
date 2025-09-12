package mothers

import (
	"github.com/google/uuid"
	orderDomain "order/internal/domain/order"
	"order/internal/tests/testutils/builders"
)

func DefaultOrder() *orderDomain.Order {
	return builders.NewOrderBuilder().Build()
}

func OrderDelivering() *orderDomain.Order {
	courierID := uuid.New()
	return builders.NewOrderBuilder().
		WithStatus(orderDomain.Delivering).
		WithDelivery(orderDomain.Delivery{
			CourierID: &courierID,
			Address:   "address",
			Arrived:   nil,
		}).
		Build()
}

func ListOfOrders(n int) []*orderDomain.Order {
	orders := make([]*orderDomain.Order, 0, n)
	for i := 0; i < n; i++ {
		orders = append(orders, DefaultOrder())
	}
	return orders
}
