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
	productModel := ToModel(product)
	res := r.db.WithContext(ctx).Create(&productModel)
	return ParseError(res.Error)
}

func (r *RepositoryImpl) GetByID(ctx context.Context, productID uuid.UUID) (*productDomain.Product, error) {
	var productModel tables.Product

	res := r.db.WithContext(ctx).
		Preload("Image").
		Where("id = ?", productID).
		First(&productModel)
	if res.Error != nil {
		return nil, ParseError(res.Error)
	}

	return ToDomain(&productModel), nil
}

func (r *RepositoryImpl) GetAll(ctx context.Context, limit, offset int) ([]*productDomain.Product, error) {
	var productModels []*tables.Product

	res := r.db.WithContext(ctx).
		Limit(limit).
		Offset(offset).
		Find(&productModels)
	if res.Error != nil {
		return nil, ParseError(res.Error)
	}

	return ToDomains(productModels), nil
}

var _ productDomain.Repository = (*RepositoryImpl)(nil)
