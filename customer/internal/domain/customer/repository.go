package customer

import (
	"context"
)

type Repository interface {
	Create(ctx context.Context, customer *Customer) error
	GetByPhone(ctx context.Context, phone string) (*Customer, error)
}
