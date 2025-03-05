package courier

import (
	"context"
	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, courier *Courier) error
	GetByID(ctx context.Context, orderID uuid.UUID) (*Courier, error)
	GetAll(ctx context.Context) ([]*Courier, error)
}
