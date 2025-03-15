package response

import (
	orderDomain "order/internal/domain/order"
	orderv1 "order/internal/presentation/grpc"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func MapStatus(status orderDomain.Status) orderv1.OrderStatus {
	switch status {
	case orderDomain.Created:
		return orderv1.OrderStatus_CREATED
	case orderDomain.CanceledCourierNotFound:
		return orderv1.OrderStatus_CANCELED_COURIER_NOT_FOUND
	case orderDomain.CanceledOutOfStock:
		return orderv1.OrderStatus_CANCELED_OUT_OF_STOCK
	case orderDomain.Delivering:
		return orderv1.OrderStatus_DELIVERING
	case orderDomain.Delivered:
		return orderv1.OrderStatus_DELIVERED
	case orderDomain.CustomerCanceled:
		return orderv1.OrderStatus_CUSTOMER_CANCELED
	default:
		return orderv1.OrderStatus_CREATED
	}
}

func ToCreateOrderResponse(orderID uuid.UUID) *orderv1.CreateOrderResponse {
	return &orderv1.CreateOrderResponse{
		OrderId: orderID.String(),
	}
}

func ToEmptyResponse() *emptypb.Empty {
	return &emptypb.Empty{}
}

func ToOrderItemResponse(item orderDomain.Item) *orderv1.OrderItem {
	return &orderv1.OrderItem{
		ProductId: item.ProductID.String(),
		Price:     item.Price.InexactFloat64(),
		Count:     int32(item.Count),
	}
}

func ToOrderItemsResponse(items []orderDomain.Item) []*orderv1.OrderItem {
	resp := make([]*orderv1.OrderItem, 0, len(items))
	for _, item := range items {
		resp = append(resp, ToOrderItemResponse(item))
	}
	return resp
}

func ToOrderResponse(order *orderDomain.Order) *orderv1.Order {
	var courierID *string
	if order.Delivery.CourierID != nil {
		courierId := order.Delivery.CourierID.String()
		courierID = &courierId
	}

	var arrived *timestamppb.Timestamp
	if order.Delivery.Arrived != nil {
		arrived = timestamppb.New(*order.Delivery.Arrived)
	}

	return &orderv1.Order{
		OrderId:    order.ID.String(),
		CustomerId: order.CustomerID.String(),
		Status:     MapStatus(order.Status),
		Version:    order.Version.String(),
		Items:      ToOrderItemsResponse(order.Items),
		Delivery: &orderv1.Delivery{
			CourierId: courierID,
			Address:   order.Delivery.Address,
			Arrived:   arrived,
		},
		Created: timestamppb.New(order.Created),
	}
}

func ToOrdersResponse(orders []*orderDomain.Order) []*orderv1.Order {
	resp := make([]*orderv1.Order, 0, len(orders))
	for _, order := range orders {
		resp = append(resp, ToOrderResponse(order))
	}
	return resp
}

func ToGetOrdersByCustomerResponse(orders []*orderDomain.Order) *orderv1.GetOrdersByCustomerResponse {
	return &orderv1.GetOrdersByCustomerResponse{
		Orders: ToOrdersResponse(orders),
	}
}

func ToGetCurrentOrdersByCourierResponse(orders []*orderDomain.Order) *orderv1.GetCurrentOrdersByCourierResponse {
	return &orderv1.GetCurrentOrdersByCourierResponse{
		Orders: ToOrdersResponse(orders),
	}
}
