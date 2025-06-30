package courier

import (
	"context"
)

type Repository interface {
	Create(ctx context.Context, courier *Courier) error
	GetByPhone(ctx context.Context, phone string) (*Courier, error)
	GetAll(ctx context.Context) ([]*Courier, error)
}
