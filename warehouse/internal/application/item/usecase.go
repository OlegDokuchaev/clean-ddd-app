package item

import (
	"context"
	itemDomain "warehouse/internal/domain/item"

	"github.com/google/uuid"
)

type UseCase interface {
	Create(ctx context.Context, data CreateDto) (uuid.UUID, error)
	Reserve(ctx context.Context, data ReserveDto) error
	Release(ctx context.Context, data ReleaseDto) error
	GetAll(ctx context.Context) ([]*itemDomain.Item, error)
}
