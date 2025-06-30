package courier

import (
	"context"
	courierDomain "courier/internal/domain/courier"
	"github.com/stretchr/testify/mock"
)

type RepositoryMock struct {
	mock.Mock
}

func (r *RepositoryMock) Create(ctx context.Context, courier *courierDomain.Courier) error {
	args := r.Called(ctx, courier)
	return args.Error(0)
}

func (r *RepositoryMock) GetByPhone(ctx context.Context, phone string) (*courierDomain.Courier, error) {
	args := r.Called(ctx, phone)
	return args.Get(0).(*courierDomain.Courier), args.Error(1)
}

func (r *RepositoryMock) GetAll(ctx context.Context) ([]*courierDomain.Courier, error) {
	args := r.Called(ctx)
	return args.Get(0).([]*courierDomain.Courier), args.Error(1)
}

var _ courierDomain.Repository = (*RepositoryMock)(nil)
