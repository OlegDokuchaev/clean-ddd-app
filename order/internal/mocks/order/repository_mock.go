package order

import (
	"context"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	orderDomain "order/internal/domain/order"
)

type RepositoryMock struct {
	mock.Mock
}

func (r *RepositoryMock) Create(ctx context.Context, order *orderDomain.Order) error {
	args := r.Called(ctx, order)
	return args.Error(0)
}

func (r *RepositoryMock) Update(ctx context.Context, order *orderDomain.Order) error {
	args := r.Called(ctx, order)
	return args.Error(0)
}

func (r *RepositoryMock) GetById(ctx context.Context, orderID uuid.UUID) (*orderDomain.Order, error) {
	args := r.Called(ctx, orderID)
	return args.Get(0).(*orderDomain.Order), args.Error(1)
}

func (r *RepositoryMock) GetAllByCustomer(ctx context.Context, customerID uuid.UUID) ([]*orderDomain.Order, error) {
	args := r.Called(ctx, customerID)
	return args.Get(0).([]*orderDomain.Order), args.Error(1)
}

func (r *RepositoryMock) GetCurrentByCourier(ctx context.Context, courierID uuid.UUID) ([]*orderDomain.Order, error) {
	args := r.Called(ctx, courierID)
	return args.Get(0).([]*orderDomain.Order), args.Error(1)
}
