package order

import (
	orderGRPC "api-gateway/gen/order/v1"
	"api-gateway/internal/adapter/output/clients/response"
	orderDto "api-gateway/internal/domain/dtos/order"
	orderClient "api-gateway/internal/port/output/clients/order"
	"context"
	"github.com/google/uuid"
)

type ClientImpl struct {
	client orderGRPC.OrderServiceClient
}

func NewClient(client orderGRPC.OrderServiceClient) orderClient.Client {
	return &ClientImpl{
		client: client,
	}
}

func (c *ClientImpl) Create(ctx context.Context, data orderClient.CreateDto) (uuid.UUID, error) {
	in := toCreateRequest(data)

	out, err := c.client.CreateOrder(ctx, in)
	if err != nil {
		return uuid.Nil, response.ParseGRPCError(err)
	}

	orderID, err := response.ToUUID(out.OrderId)
	if err != nil {
		return uuid.Nil, err
	}

	return orderID, nil
}

func (c *ClientImpl) CancelByCustomer(ctx context.Context, orderID uuid.UUID, customerID uuid.UUID) error {
	in := toCancelByCustomerRequest(orderID)

	_, err := c.client.CancelOrderByCustomer(ctx, in)
	if err != nil {
		return response.ParseGRPCError(err)
	}

	return nil
}

func (c *ClientImpl) Complete(ctx context.Context, orderID uuid.UUID, courierID uuid.UUID) error {
	in := toCompleteDeliveryRequest(orderID)

	_, err := c.client.CompleteDelivery(ctx, in)
	if err != nil {
		return response.ParseGRPCError(err)
	}

	return nil
}

func (c *ClientImpl) GetByCustomer(
	ctx context.Context,
	customerID uuid.UUID,
	limit int,
	offset int,
) ([]*orderDto.OrderDto, error) {
	in := toGetByCustomerRequest(customerID, limit, offset)

	out, err := c.client.GetOrdersByCustomer(ctx, in)
	if err != nil {
		return nil, response.ParseGRPCError(err)
	}

	orders, err := toOrders(out.Orders)
	if err != nil {
		return nil, err
	}

	return orders, nil
}

func (c *ClientImpl) GetCurrentByCourier(
	ctx context.Context,
	courierID uuid.UUID,
	limit int,
	offset int,
) ([]*orderDto.OrderDto, error) {
	in := toGetCurrentByCourierRequest(courierID, limit, offset)

	out, err := c.client.GetCurrentOrdersByCourier(ctx, in)
	if err != nil {
		return nil, response.ParseGRPCError(err)
	}

	orders, err := toOrders(out.Orders)
	if err != nil {
		return nil, err
	}

	return orders, nil
}

var _ orderClient.Client = (*ClientImpl)(nil)
