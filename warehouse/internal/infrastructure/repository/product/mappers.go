package product

import (
	productDomain "warehouse/internal/domain/product"
	"warehouse/internal/infrastructure/db/tables"

	"github.com/google/uuid"
)

func ToDomain(model *tables.Product) *productDomain.Product {
	return &productDomain.Product{
		ID:      model.ID,
		Name:    model.Name,
		Price:   model.Price,
		Created: model.Created,
		Image:   toProductImageDomain(model.Image),
	}
}

func toProductImageDomain(model tables.ProductImage) productDomain.Image {
	return productDomain.Image{
		Path: model.Path,
	}
}

func ToModel(domain *productDomain.Product) tables.Product {
	return tables.Product{
		ID:      domain.ID,
		Name:    domain.Name,
		Price:   domain.Price,
		Created: domain.Created,
		Image:   toProductImageModel(domain.ID, domain.Image),
	}
}

func toProductImageModel(productID uuid.UUID, domain productDomain.Image) tables.ProductImage {
	return tables.ProductImage{
		ID:        uuid.New(),
		ProductID: productID,
		Path:      domain.Path,
	}
}
