package warehouse

import (
	warehouseDto "api-gateway/internal/domain/dtos/warehouse"
	"context"
	"io"

	"github.com/google/uuid"
)

type UseCase interface {
	ReserveItems(ctx context.Context, items []warehouseDto.ItemInfoDto, adminToken string) error
	ReleaseItems(ctx context.Context, items []warehouseDto.ItemInfoDto, adminToken string) error
	CreateProduct(ctx context.Context, data warehouseDto.CreateProductDto, adminToken string) (uuid.UUID, error)
	GetAllItems(ctx context.Context) ([]*warehouseDto.ItemDto, error)
	UpdateProductImage(ctx context.Context, productID uuid.UUID, fileReader io.Reader, contentType string, adminToken string) error
}
