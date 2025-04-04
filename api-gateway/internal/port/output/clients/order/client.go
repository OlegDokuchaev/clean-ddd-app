package order

import (
	orderDto "api-gateway/internal/domain/dtos/order"
	"context"

	"github.com/google/uuid"
)

type Client interface {
	Create(ctx context.Context, data CreateDto) (uuid.UUID, error)
	CancelByCustomer(ctx context.Context, orderID uuid.UUID, customerID uuid.UUID) error
	Complete(ctx context.Context, orderID uuid.UUID, courierID uuid.UUID) error
	GetByCustomer(ctx context.Context, customerID uuid.UUID) ([]*orderDto.OrderDto, error)
	GetCurrentByCourier(ctx context.Context, courierID uuid.UUID) ([]*orderDto.OrderDto, error)
}
