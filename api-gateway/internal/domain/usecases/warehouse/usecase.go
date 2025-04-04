package warehouse

import (
	"context"
	"github.com/google/uuid"
)

type UseCase interface {
	ReserveItems(ctx context.Context, items []ItemInfoDto, adminToken string) error
	ReleaseItems(ctx context.Context, items []ItemInfoDto, adminToken string) error
	CreateProduct(ctx context.Context, data CreateProductDto, adminToken string) (uuid.UUID, error)
	GetAllItems(ctx context.Context) ([]*ItemDto, error)
}
