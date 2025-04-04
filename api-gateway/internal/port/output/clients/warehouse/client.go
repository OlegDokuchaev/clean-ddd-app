package warehouse

import (
	warehouseUseCase "api-gateway/internal/domain/usecases/warehouse"
	"context"

	"github.com/google/uuid"
)

type Client interface {
	ReserveItems(ctx context.Context, items []warehouseUseCase.ItemInfoDto) error
	ReleaseItems(ctx context.Context, items []warehouseUseCase.ItemInfoDto) error
	CreateProduct(ctx context.Context, data warehouseUseCase.CreateProductDto) (uuid.UUID, error)
	GetAllItems(ctx context.Context) ([]*warehouseUseCase.ItemDto, error)
}
