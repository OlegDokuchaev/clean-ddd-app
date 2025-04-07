package response

import (
	orderDto "api-gateway/internal/domain/dtos/order"
)

func ToOrderResponse(order *orderDto.OrderDto) OrderResponse {
	return OrderResponse{
		ID:         order.ID,
		CustomerID: order.CustomerID,
		Status:     string(order.Status),
		Created:    order.Created,
		Version:    order.Version.String(),
		Delivery:   toDeliverySchema(order.Delivery),
		Items:      toItemSchemas(order.Items),
	}
}

func ToOrdersResponse(orders []*orderDto.OrderDto) OrdersResponse {
	result := make([]OrderResponse, 0, len(orders))
	for _, order := range orders {
		result = append(result, ToOrderResponse(order))
	}
	return OrdersResponse{Orders: result}
}

func toItemSchemas(items []orderDto.ItemDto) []ItemSchema {
	result := make([]ItemSchema, 0, len(items))
	for _, item := range items {
		result = append(result, toItemSchema(item))
	}
	return result
}

func toItemSchema(item orderDto.ItemDto) ItemSchema {
	return ItemSchema{
		ProductID: item.ProductID,
		Price:     item.Price,
		Count:     item.Count,
	}
}

func toDeliverySchema(delivery orderDto.DeliveryDto) DeliverySchema {
	return DeliverySchema{
		CourierID: delivery.CourierID,
		Address:   delivery.Address,
		Arrived:   delivery.Arrived,
	}
}
