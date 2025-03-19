package handler

import (
	"context"
	orderUsecase "order/internal/application/order/usecase"
	orderv1 "order/internal/presentation/grpc"
	"order/internal/presentation/grpc/request"
	"order/internal/presentation/grpc/response"

	"google.golang.org/protobuf/types/known/emptypb"
)

type OrderServiceHandler struct {
	orderv1.UnimplementedOrderServiceServer

	usecase orderUsecase.UseCase
}

func NewOrderServiceHandler(usecase orderUsecase.UseCase) *OrderServiceHandler {
	return &OrderServiceHandler{
		usecase: usecase,
	}
}

func (h *OrderServiceHandler) CreateOrder(
	ctx context.Context,
	req *orderv1.CreateOrderRequest,
) (*orderv1.CreateOrderResponse, error) {
	data, err := request.ToCreateDto(req)
	if err != nil {
		return nil, err
	}

	orderID, err := h.usecase.Create(ctx, data)
	if err != nil {
		return nil, response.ParseError(err)
	}

	return response.ToCreateOrderResponse(orderID), nil
}

func (h *OrderServiceHandler) CancelOrderByCustomer(ctx context.Context, req *orderv1.CancelOrderByCustomerRequest) (*emptypb.Empty, error) {
	orderID, err := request.ParseUUID(req.OrderId)
	if err != nil {
		return nil, err
	}

	if err = h.usecase.CancelByCustomer(ctx, orderID); err != nil {
		return nil, response.ParseError(err)
	}

	return response.ToEmptyResponse(), nil
}

func (h *OrderServiceHandler) CancelOrderOutOfStock(ctx context.Context, req *orderv1.CancelOrderOutOfStockRequest) (*emptypb.Empty, error) {
	orderID, err := request.ParseUUID(req.OrderId)
	if err != nil {
		return nil, err
	}

	if err = h.usecase.CancelOutOfStock(ctx, orderID); err != nil {
		return nil, response.ParseError(err)
	}

	return response.ToEmptyResponse(), nil
}

func (h *OrderServiceHandler) CancelOrderCourierNotFound(ctx context.Context, req *orderv1.CancelOrderCourierNotFoundRequest) (*emptypb.Empty, error) {
	orderID, err := request.ParseUUID(req.OrderId)
	if err != nil {
		return nil, err
	}

	if err = h.usecase.CancelCourierNotFound(ctx, orderID); err != nil {
		return nil, response.ParseError(err)
	}

	return response.ToEmptyResponse(), nil
}

func (h *OrderServiceHandler) BeginDelivery(ctx context.Context, req *orderv1.BeginDeliveryRequest) (*emptypb.Empty, error) {
	data, err := request.ToBeginDeliveryDto(req)
	if err != nil {
		return nil, err
	}

	if err = h.usecase.BeginDelivery(ctx, data); err != nil {
		return nil, response.ParseError(err)
	}

	return response.ToEmptyResponse(), nil
}

func (h *OrderServiceHandler) CompleteDelivery(ctx context.Context, req *orderv1.CompleteDeliveryRequest) (*emptypb.Empty, error) {
	orderID, err := request.ParseUUID(req.OrderId)
	if err != nil {
		return nil, err
	}

	if err = h.usecase.CompleteDelivery(ctx, orderID); err != nil {
		return nil, response.ParseError(err)
	}

	return response.ToEmptyResponse(), nil
}

func (h *OrderServiceHandler) GetOrdersByCustomer(ctx context.Context, req *orderv1.GetOrdersByCustomerRequest) (*orderv1.GetOrdersByCustomerResponse, error) {
	customerID, err := request.ParseUUID(req.CustomerId)
	if err != nil {
		return nil, err
	}

	orders, err := h.usecase.GetAllByCustomer(ctx, customerID)
	if err != nil {
		return nil, response.ParseError(err)
	}

	return response.ToGetOrdersByCustomerResponse(orders), nil
}

func (h *OrderServiceHandler) GetCurrentOrdersByCourier(ctx context.Context, req *orderv1.GetCurrentOrdersByCourierRequest) (*orderv1.GetCurrentOrdersByCourierResponse, error) {
	courierID, err := request.ParseUUID(req.CourierId)
	if err != nil {
		return nil, err
	}

	orders, err := h.usecase.GetCurrentByCourier(ctx, courierID)
	if err != nil {
		return nil, response.ParseError(err)
	}

	return response.ToGetCurrentOrdersByCourierResponse(orders), nil
}

var _ orderv1.OrderServiceServer = (*OrderServiceHandler)(nil)
