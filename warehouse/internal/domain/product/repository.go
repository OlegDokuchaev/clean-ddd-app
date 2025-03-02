package product

import (
	"context"
	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, product *Product) error
	GetById(ctx context.Context, productID uuid.UUID) (*Product, error)
}
