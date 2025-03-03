package item

import (
	"context"
	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, item *Item) error
	Update(ctx context.Context, item *Item) error
	GetByID(ctx context.Context, itemID uuid.UUID) (*Item, error)
	GetAll(ctx context.Context) ([]*Item, error)
	GetAllByIDs(ctx context.Context, itemIDs ...uuid.UUID) ([]*Item, error)
}
