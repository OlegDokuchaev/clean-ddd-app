package product

import (
	"context"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	productDomain "warehouse/internal/domain/product"
)

type RepositoryMock struct {
	mock.Mock
}

func (r *RepositoryMock) Create(ctx context.Context, product *productDomain.Product) error {
	args := r.Called(ctx, product)
	return args.Error(0)
}

func (r *RepositoryMock) GetByID(ctx context.Context, productID uuid.UUID) (*productDomain.Product, error) {
	args := r.Called(ctx, productID)
	return args.Get(0).(*productDomain.Product), args.Error(1)
}

func (r *RepositoryMock) GetAll(ctx context.Context, limit int, offset int) ([]*productDomain.Product, error) {
	args := r.Called(ctx, limit, offset)
	return args.Get(0).([]*productDomain.Product), args.Error(1)
}

var _ productDomain.Repository = (*RepositoryMock)(nil)
