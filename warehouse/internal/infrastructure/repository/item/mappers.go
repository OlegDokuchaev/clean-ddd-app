package item

import (
	itemDomain "warehouse/internal/domain/item"
	"warehouse/internal/infrastructure/db/tables"
	productRepository "warehouse/internal/infrastructure/repository/product"
)

func ToDomain(item *tables.Item) *itemDomain.Item {
	return &itemDomain.Item{
		ID:      item.ID,
		Count:   item.Count,
		Version: item.Version,
		Product: productRepository.ToDomain(item.Product),
	}
}

func ToDomains(items []*tables.Item) []*itemDomain.Item {
	domains := make([]*itemDomain.Item, 0, len(items))
	for _, item := range items {
		domains = append(domains, ToDomain(item))
	}
	return domains
}

func ToModel(item *itemDomain.Item) *tables.Item {
	return &tables.Item{
		ID:        item.ID,
		ProductID: item.Product.ID,
		Count:     item.Count,
		Version:   item.Version,
		Product:   productRepository.ToModel(item.Product),
	}
}
