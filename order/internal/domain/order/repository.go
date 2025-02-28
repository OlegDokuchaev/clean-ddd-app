package order

import (
	"context"
	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, order *Order) error
	Update(ctx context.Context, order *Order) error
	GetById(ctx context.Context) (*Order, error)
	GetAllByCustomer(ctx context.Context, customerID uuid.UUID) ([]*Order, error)
	GetCurrentByCourier(ctx context.Context, courierID uuid.UUID) ([]*Order, error)
}
