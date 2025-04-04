package order

import (
	"context"
	"github.com/google/uuid"
)

type UseCase interface {
	Create(ctx context.Context, data CreateDto, customerToken string) (uuid.UUID, error)
	CancelByCustomer(ctx context.Context, orderID uuid.UUID, customerToken string) error
	Complete(ctx context.Context, orderID uuid.UUID, courierToken string) error
	GetByCustomer(ctx context.Context, customerToken string) ([]*OrderDto, error)
	GetCurrentByCourier(ctx context.Context, courierToken string) ([]*OrderDto, error)
}
