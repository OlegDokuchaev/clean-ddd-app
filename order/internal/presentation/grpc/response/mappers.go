package response

import (
	"fmt"
	"math"
	orderDomain "order/internal/domain/order"
	orderv1 "order/internal/presentation/grpc"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func safeIntToInt32(v int) (int32, error) {
	if v > math.MaxInt32 || v < math.MinInt32 {
		return 0, fmt.Errorf("value %d out of int32 range", v)
	}
	return int32(v), nil
}

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

func ToOrderItemResponse(item orderDomain.Item) (*orderv1.OrderItem, error) {
	count32, err := safeIntToInt32(item.Count)
	if err != nil {
		return nil, err
	}

	return &orderv1.OrderItem{
		ProductId: item.ProductID.String(),
		Price:     item.Price.InexactFloat64(),
		Count:     count32,
	}, nil
}

func ToOrderItemsResponse(items []orderDomain.Item) ([]*orderv1.OrderItem, error) {
	resp := make([]*orderv1.OrderItem, 0, len(items))
	for _, it := range items {
		mapped, err := ToOrderItemResponse(it)
		if err != nil {
			return nil, err
		}
		resp = append(resp, mapped)
	}
	return resp, nil
}

func ToOrderResponse(order *orderDomain.Order) (*orderv1.Order, error) {
	items, err := ToOrderItemsResponse(order.Items)
	if err != nil {
		return nil, err
	}

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
		Items:      items,
		Delivery: &orderv1.Delivery{
			CourierId: courierID,
			Address:   order.Delivery.Address,
			Arrived:   arrived,
		},
		Created: timestamppb.New(order.Created),
	}, nil
}

func ToOrdersResponse(orders []*orderDomain.Order) ([]*orderv1.Order, error) {
	resp := make([]*orderv1.Order, 0, len(orders))
	for _, order := range orders {
		mapped, err := ToOrderResponse(order)
		if err != nil {
			return nil, err
		}
		resp = append(resp, mapped)
	}
	return resp, nil
}

func ToGetOrdersByCustomerResponse(orders []*orderDomain.Order) (*orderv1.GetOrdersByCustomerResponse, error) {
	mappedOrders, err := ToOrdersResponse(orders)
	if err != nil {
		return nil, err
	}

	return &orderv1.GetOrdersByCustomerResponse{
		Orders: mappedOrders,
	}, nil
}

func ToGetCurrentOrdersByCourierResponse(orders []*orderDomain.Order) (*orderv1.GetCurrentOrdersByCourierResponse, error) {
	mappedOrders, err := ToOrdersResponse(orders)
	if err != nil {
		return nil, err
	}

	return &orderv1.GetCurrentOrdersByCourierResponse{
		Orders: mappedOrders,
	}, nil
}
