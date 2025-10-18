package warehouse

import (
	warehouseDto "api-gateway/internal/domain/dtos/warehouse"
	"context"
	"io"

	"github.com/google/uuid"
)

type Client interface {
	ReserveItems(ctx context.Context, items []warehouseDto.ItemInfoDto) error
	ReleaseItems(ctx context.Context, items []warehouseDto.ItemInfoDto) error
	CreateProduct(ctx context.Context, data warehouseDto.CreateProductDto) (uuid.UUID, error)
	GetAllItems(ctx context.Context, limit int, offset int) ([]*warehouseDto.ItemDto, error)
	UpdateProductImage(ctx context.Context, productID uuid.UUID, fileReader io.Reader, contentType string) error
	GetProductImage(ctx context.Context, productID uuid.UUID) (fileReader io.Reader, contentType string, err error)
}
