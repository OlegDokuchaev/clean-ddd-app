package mothers

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	orderDomain "order/internal/domain/order"
	"order/internal/tests/testutils/builders"
)

func DefaultOrder() *orderDomain.Order {
	return builders.NewOrderBuilder().Build()
}

func OrderWithSingleItem(price int64, count int) *orderDomain.Order {
	return builders.NewOrderBuilder().
		WithItems([]orderDomain.Item{
			{
				ProductID: uuid.New(),
				Price:     decimal.NewFromInt(price),
				Count:     count,
			},
		}).
		Build()
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
