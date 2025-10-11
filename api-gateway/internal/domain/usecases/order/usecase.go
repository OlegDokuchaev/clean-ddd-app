package order

import (
	orderDto "api-gateway/internal/domain/dtos/order"
	"context"
	"github.com/google/uuid"
)

type UseCase interface {
	Create(ctx context.Context, data orderDto.CreateDto, customerToken string) (uuid.UUID, error)
	CancelByCustomer(ctx context.Context, orderID uuid.UUID, customerToken string) error
	Complete(ctx context.Context, orderID uuid.UUID, courierToken string) error
	GetByCustomer(ctx context.Context, limit int, offset int, customerToken string) ([]*orderDto.OrderDto, error)
	GetCurrentByCourier(ctx context.Context, courierToken string) ([]*orderDto.OrderDto, error)
}
