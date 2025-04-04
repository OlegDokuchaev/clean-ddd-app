package warehouse

import (
	warehouseDto "api-gateway/internal/domain/dtos/warehouse"
	"context"

	"github.com/google/uuid"
)

type Client interface {
	ReserveItems(ctx context.Context, items []warehouseDto.ItemInfoDto) error
	ReleaseItems(ctx context.Context, items []warehouseDto.ItemInfoDto) error
	CreateProduct(ctx context.Context, data warehouseDto.CreateProductDto) (uuid.UUID, error)
	GetAllItems(ctx context.Context) ([]*warehouseDto.ItemDto, error)
}
