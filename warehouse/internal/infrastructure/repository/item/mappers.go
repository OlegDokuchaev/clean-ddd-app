package item

import (
	itemDomain "warehouse/internal/domain/item"
	"warehouse/internal/infrastructure/db/tables"
	productRepository "warehouse/internal/infrastructure/repository/product"
)

func ToDomain(model *tables.Item) *itemDomain.Item {
	return &itemDomain.Item{
		ID:      model.ID,
		Count:   model.Count,
		Version: model.Version,
		Product: productRepository.ToDomain(&model.Product),
	}
}

func ToDomains(models []*tables.Item) []*itemDomain.Item {
	domains := make([]*itemDomain.Item, 0, len(models))
	for _, model := range models {
		domains = append(domains, ToDomain(model))
	}
	return domains
}

func ToModel(domain *itemDomain.Item) *tables.Item {
	return &tables.Item{
		ID:        domain.ID,
		ProductID: domain.Product.ID,
		Count:     domain.Count,
		Version:   domain.Version,
		Product:   productRepository.ToModel(domain.Product),
	}
}
