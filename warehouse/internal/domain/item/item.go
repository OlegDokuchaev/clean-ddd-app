package item

import (
	"github.com/google/uuid"
	productDomain "warehouse/internal/domain/product"
)

type Item struct {
	ID      uuid.UUID
	Count   int
	Product *productDomain.Product
	Version uuid.UUID
}

func (i *Item) Reserve(count int) error {
	finalCount := i.Count - count
	if finalCount <= 0 {
		return ErrInvalidItemCount
	}
	i.Count = finalCount
	return nil
}

func (i *Item) Release(count int) error {
	i.Count += count
	return nil
}
