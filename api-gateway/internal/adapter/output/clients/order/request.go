package order

import (
	orderGRPC "api-gateway/gen/order/v1"
	orderDto "api-gateway/internal/domain/dtos/order"
	orderClient "api-gateway/internal/port/output/clients/order"

	"github.com/google/uuid"
)

func toOrderItem(item orderDto.ItemDto) *orderGRPC.OrderItem {
	return &orderGRPC.OrderItem{
		ProductId: item.ProductID.String(),
		Price:     item.Price.InexactFloat64(),
		Count:     int32(item.Count),
	}
}

func toOrderItems(items []orderDto.ItemDto) []*orderGRPC.OrderItem {
	orderItems := make([]*orderGRPC.OrderItem, 0, len(items))
	for _, item := range items {
		orderItems = append(orderItems, toOrderItem(item))
	}
	return orderItems
}

func toCreateRequest(data orderClient.CreateDto) *orderGRPC.CreateOrderRequest {
	return &orderGRPC.CreateOrderRequest{
		CustomerId: data.CustomerID.String(),
		Address:    data.Address,
		Items:      toOrderItems(data.Items),
	}
}

func toCancelByCustomerRequest(orderID uuid.UUID) *orderGRPC.CancelOrderByCustomerRequest {
	return &orderGRPC.CancelOrderByCustomerRequest{
		OrderId: orderID.String(),
	}
}

func toCompleteDeliveryRequest(orderID uuid.UUID) *orderGRPC.CompleteDeliveryRequest {
	return &orderGRPC.CompleteDeliveryRequest{
		OrderId: orderID.String(),
	}
}

func toGetByCustomerRequest(customerID uuid.UUID) *orderGRPC.GetOrdersByCustomerRequest {
	return &orderGRPC.GetOrdersByCustomerRequest{
		CustomerId: customerID.String(),
	}
}

func toGetCurrentByCourierRequest(courierID uuid.UUID) *orderGRPC.GetCurrentOrdersByCourierRequest {
	return &orderGRPC.GetCurrentOrdersByCourierRequest{
		CourierId: courierID.String(),
	}
}
