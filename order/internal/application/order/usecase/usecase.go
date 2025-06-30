package usecase

import (
	"context"
	"github.com/google/uuid"
	orderDomain "order/internal/domain/order"
)

type UseCase interface {
	Create(ctx context.Context, data CreateDto) (uuid.UUID, error)
	CancelByCustomer(ctx context.Context, orderID uuid.UUID) error
	CancelOutOfStock(ctx context.Context, orderID uuid.UUID) error
	CancelCourierNotFound(ctx context.Context, orderID uuid.UUID) error
	BeginDelivery(ctx context.Context, data BeginDeliveryDto) error
	CompleteDelivery(ctx context.Context, orderID uuid.UUID) error
	GetAllByCustomer(ctx context.Context, customerID uuid.UUID) ([]*orderDomain.Order, error)
	GetCurrentByCourier(ctx context.Context, courierID uuid.UUID) ([]*orderDomain.Order, error)
}
