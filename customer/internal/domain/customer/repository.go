package customer

import (
	"context"
	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, customer *Customer) error
	Save(ctx context.Context, customer *Customer) error
	GetByID(ctx context.Context, id uuid.UUID) (*Customer, error)
	GetByPhone(ctx context.Context, phone string) (*Customer, error)
	GetByEmail(ctx context.Context, email string) (*Customer, error)
}
