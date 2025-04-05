package order

import (
	orderGRPC "api-gateway/gen/order/v1"
	"api-gateway/internal/adapter/output/clients/response"
	orderDto "api-gateway/internal/domain/dtos/order"
)

func toOrders(protoOrders []*orderGRPC.Order) ([]*orderDto.OrderDto, error) {
	orders := make([]*orderDto.OrderDto, 0, len(protoOrders))
	for _, protoOrder := range protoOrders {
		order, err := toOrder(protoOrder)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}
	return orders, nil
}

func toItem(protoItem *orderGRPC.OrderItem) (orderDto.ItemDto, error) {
	productId, err := response.ToUUID(protoItem.ProductId)
	if err != nil {
		return orderDto.ItemDto{}, err
	}

	return orderDto.ItemDto{
		ProductID: productId,
		Price:     response.ToDecimal(protoItem.Price),
		Count:     int(protoItem.Count),
	}, nil
}

func toItems(protoItems []*orderGRPC.OrderItem) ([]orderDto.ItemDto, error) {
	items := make([]orderDto.ItemDto, 0, len(protoItems))
	for _, protoItem := range protoItems {
		item, err := toItem(protoItem)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

func toDelivery(protoDelivery *orderGRPC.Delivery) (orderDto.DeliveryDto, error) {
	var deliveryDto orderDto.DeliveryDto

	if protoDelivery.CourierId != nil {
		courierId, err := response.ToUUID(*protoDelivery.CourierId)
		if err != nil {
			return orderDto.DeliveryDto{}, err
		}
		deliveryDto.CourierID = &courierId
	}

	if protoDelivery.Arrived != nil {
		t := protoDelivery.Arrived.AsTime()
		deliveryDto.Arrived = &t
	}

	deliveryDto.Address = protoDelivery.Address

	return deliveryDto, nil
}

func toOrder(protoOrder *orderGRPC.Order) (*orderDto.OrderDto, error) {
	orderID, err := response.ToUUID(protoOrder.OrderId)
	if err != nil {
		return nil, err
	}

	customerID, err := response.ToUUID(protoOrder.CustomerId)
	if err != nil {
		return nil, err
	}

	versionID, err := response.ToUUID(protoOrder.Version)
	if err != nil {
		return nil, err
	}

	items, err := toItems(protoOrder.Items)
	if err != nil {
		return nil, err
	}

	delivery, err := toDelivery(protoOrder.Delivery)
	if err != nil {
		return nil, err
	}

	return &orderDto.OrderDto{
		ID:         orderID,
		CustomerID: customerID,
		Status:     toOrderStatus(protoOrder.Status),
		Created:    protoOrder.Created.AsTime(),
		Version:    versionID,
		Delivery:   delivery,
		Items:      items,
	}, nil
}

func toOrderStatus(protoStatus orderGRPC.OrderStatus) orderDto.Status {
	switch protoStatus {
	case orderGRPC.OrderStatus_CREATED:
		return orderDto.Created
	case orderGRPC.OrderStatus_CANCELED_COURIER_NOT_FOUND:
		return orderDto.CanceledCourierNotFound
	case orderGRPC.OrderStatus_CANCELED_OUT_OF_STOCK:
		return orderDto.CanceledOutOfStock
	case orderGRPC.OrderStatus_DELIVERING:
		return orderDto.Delivering
	case orderGRPC.OrderStatus_DELIVERED:
		return orderDto.Delivered
	case orderGRPC.OrderStatus_CUSTOMER_CANCELED:
		return orderDto.CustomerCanceled
	default:
		return orderDto.Created
	}
}
