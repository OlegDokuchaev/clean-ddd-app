package order

import (
	orderUseCase "api-gateway/internal/domain/usecases/order"
	"context"

	"github.com/google/uuid"
)

type Client interface {
	Create(ctx context.Context, data CreateDto) (uuid.UUID, error)
	CancelByCustomer(ctx context.Context, orderID uuid.UUID, customerID uuid.UUID) error
	Complete(ctx context.Context, orderID uuid.UUID, courierID uuid.UUID) error
	GetByCustomer(ctx context.Context, customerID uuid.UUID) ([]*orderUseCase.OrderDto, error)
	GetCurrentByCourier(ctx context.Context, courierID uuid.UUID) ([]*orderUseCase.OrderDto, error)
}
