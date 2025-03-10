package product

import (
	productDomain "warehouse/internal/domain/product"
	"warehouse/internal/infrastructure/db/tables"
)

func ToDomain(model *tables.Product) *productDomain.Product {
	return &productDomain.Product{
		ID:      model.ID,
		Name:    model.Name,
		Price:   model.Price,
		Created: model.Created,
	}
}

func ToModel(domain *productDomain.Product) tables.Product {
	return tables.Product{
		ID:      domain.ID,
		Name:    domain.Name,
		Price:   domain.Price,
		Created: domain.Created,
	}
}
