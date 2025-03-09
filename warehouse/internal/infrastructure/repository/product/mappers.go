package product

import (
	productDomain "warehouse/internal/domain/product"
	"warehouse/internal/infrastructure/db/tables"
)

func ToDomain(product tables.Product) *productDomain.Product {
	return &productDomain.Product{
		ID:      product.ID,
		Name:    product.Name,
		Price:   product.Price,
		Created: product.Created,
	}
}

func ToModel(product *productDomain.Product) tables.Product {
	return tables.Product{
		ID:      product.ID,
		Name:    product.Name,
		Price:   product.Price,
		Created: product.Created,
	}
}
