package product

import (
	"context"
	productDomain "warehouse/internal/domain/product"
	"warehouse/internal/infrastructure/db/tables"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RepositoryImpl struct {
	db *gorm.DB
}

func New(db *gorm.DB) *RepositoryImpl {
	return &RepositoryImpl{db: db}
}

func (r *RepositoryImpl) Create(ctx context.Context, product *productDomain.Product) error {
	productModel := toModel(product)
	res := r.db.WithContext(ctx).Create(&productModel)
	return ParseError(res.Error)
}

func (r *RepositoryImpl) GetByID(ctx context.Context, productID uuid.UUID) (*productDomain.Product, error) {
	var productModel tables.Product
	res := r.db.WithContext(ctx).Where("id = ?", productID).First(&productModel)
	if res.Error != nil {
		return nil, ParseError(res.Error)
	}
	return toDomain(productModel), nil
}

var _ productDomain.Repository = (*RepositoryImpl)(nil)
