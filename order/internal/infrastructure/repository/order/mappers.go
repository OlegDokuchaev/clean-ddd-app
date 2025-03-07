package order

import (
	orderDomain "order/internal/domain/order"
	"order/internal/infrastructure/db/tables"

	"github.com/google/uuid"
)

func toOrder(orderModel tables.Order) *orderDomain.Order {
	return &orderDomain.Order{
		ID:         orderModel.ID,
		CustomerID: orderModel.CustomerID,
		Status:     orderModel.Status,
		Created:    orderModel.Created,
		Version:    orderModel.Version,
		Items:      toOrderItems(orderModel.Items),
		Delivery:   toDelivery(orderModel.Delivery),
	}
}

func toOrders(orderModels []*tables.Order) []*orderDomain.Order {
	orders := make([]*orderDomain.Order, 0, len(orderModels))
	for _, orderModel := range orderModels {
		orders = append(orders, toOrder(*orderModel))
	}
	return orders
}

func toOrderItems(items []tables.OrderItem) []orderDomain.Item {
	orderItems := make([]orderDomain.Item, 0, len(items))
	for _, item := range items {
		orderItems = append(orderItems, toOrderItem(item))
	}
	return orderItems
}

func toOrderItem(item tables.OrderItem) orderDomain.Item {
	return orderDomain.Item{
		ProductID: item.ProductID,
		Price:     item.Price,
		Count:     item.Count,
	}
}

func toDelivery(delivery tables.Delivery) orderDomain.Delivery {
	return orderDomain.Delivery{
		CourierID: delivery.CourierID,
		Address:   delivery.Address,
		Arrived:   delivery.Arrived,
	}
}

func toOrderModel(order *orderDomain.Order) tables.Order {
	return tables.Order{
		ID:         order.ID,
		CustomerID: order.CustomerID,
		Status:     order.Status,
		Created:    order.Created,
		Version:    order.Version,
		Items:      toOrderItemsModel(order.ID, order.Items),
		Delivery:   toDeliveryModel(order.ID, order.Delivery),
	}
}

func toOrderItemsModel(orderID uuid.UUID, items []orderDomain.Item) []tables.OrderItem {
	orderItems := make([]tables.OrderItem, 0, len(items))
	for _, item := range items {
		orderItems = append(orderItems, toOrderItemModel(orderID, item))
	}
	return orderItems
}

func toOrderItemModel(orderID uuid.UUID, item orderDomain.Item) tables.OrderItem {
	return tables.OrderItem{
		ID:        uuid.New(),
		OrderID:   orderID,
		ProductID: item.ProductID,
		Price:     item.Price,
		Count:     item.Count,
	}
}

func toDeliveryModel(orderID uuid.UUID, delivery orderDomain.Delivery) tables.Delivery {
	return tables.Delivery{
		ID:        uuid.New(),
		OrderID:   orderID,
		CourierID: delivery.CourierID,
		Address:   delivery.Address,
		Arrived:   delivery.Arrived,
	}
}
