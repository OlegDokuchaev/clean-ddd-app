package customer

import (
	"context"

	customerDomain "customer/internal/domain/customer"

	"github.com/stretchr/testify/mock"
)

type RepositoryMock struct {
	mock.Mock
}

func (m *RepositoryMock) Create(ctx context.Context, customer *customerDomain.Customer) error {
	args := m.Called(ctx, customer)
	return args.Error(0)
}

func (m *RepositoryMock) GetByPhone(ctx context.Context, phone string) (*customerDomain.Customer, error) {
	args := m.Called(ctx, phone)
	return args.Get(0).(*customerDomain.Customer), args.Error(1)
}

var _ customerDomain.Repository = (*RepositoryMock)(nil)
