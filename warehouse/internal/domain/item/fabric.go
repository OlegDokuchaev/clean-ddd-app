package item

import (
	"github.com/google/uuid"
	productDomain "warehouse/internal/domain/product"
)

func Create(product *productDomain.Product, count int) (*Item, error) {
	if count <= 0 {
		return nil, ErrInvalidItemCount
	}
	return &Item{
		ID:      uuid.New(),
		Count:   count,
		Product: product,
		Version: uuid.New(),
	}, nil
}
