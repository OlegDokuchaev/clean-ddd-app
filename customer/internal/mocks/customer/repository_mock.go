package customer

import (
	"context"
	"github.com/google/uuid"

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

func (m *RepositoryMock) Save(ctx context.Context, customer *customerDomain.Customer) error {
	args := m.Called(ctx, customer)
	return args.Error(0)
}

func (m *RepositoryMock) GetByID(ctx context.Context, id uuid.UUID) (*customerDomain.Customer, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*customerDomain.Customer), args.Error(1)
}

func (m *RepositoryMock) GetByPhone(ctx context.Context, phone string) (*customerDomain.Customer, error) {
	args := m.Called(ctx, phone)
	return args.Get(0).(*customerDomain.Customer), args.Error(1)
}

func (m *RepositoryMock) GetByEmail(ctx context.Context, email string) (*customerDomain.Customer, error) {
	args := m.Called(ctx, email)
	return args.Get(0).(*customerDomain.Customer), args.Error(1)
}

var _ customerDomain.Repository = (*RepositoryMock)(nil)
