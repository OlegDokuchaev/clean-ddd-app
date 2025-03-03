package mocks

import (
	"context"
	"github.com/stretchr/testify/mock"
	"testing"
	itemDomain "warehouse/internal/domain/item"
	outboxDomain "warehouse/internal/domain/outbox"
	productDomain "warehouse/internal/domain/product"
	"warehouse/internal/domain/uow"
	itemMock "warehouse/internal/mocks/item"
	outboxMock "warehouse/internal/mocks/outbox"
	productMock "warehouse/internal/mocks/product"
)

type UoWMock struct {
	ProductMock *productMock.RepositoryMock
	ItemMock    *itemMock.RepositoryMock
	OutboxMock  *outboxMock.RepositoryMock

	mock.Mock
}

func NewUowMock() *UoWMock {
	product := &productMock.RepositoryMock{}
	item := &itemMock.RepositoryMock{}
	outbox := &outboxMock.RepositoryMock{}
	return &UoWMock{
		ProductMock: product,
		ItemMock:    item,
		OutboxMock:  outbox,
	}
}

func (u *UoWMock) Product() productDomain.Repository {
	return u.ProductMock
}
func (u *UoWMock) Item() itemDomain.Repository {
	return u.ItemMock
}
func (u *UoWMock) Outbox() outboxDomain.Repository {
	return u.OutboxMock
}

func (u *UoWMock) Transaction(ctx context.Context, fn func(u uow.UoW) error) error {
	args := u.Called(ctx, fn)
	if len(args) == 0 {
		return fn(u)
	}
	return args.Error(0)
}

func (u *UoWMock) AssertExpectations(t *testing.T) {
	u.ProductMock.AssertExpectations(t)
	u.ItemMock.AssertExpectations(t)
	u.OutboxMock.AssertExpectations(t)
	u.Mock.AssertExpectations(t)
}

var _ uow.UoW = (*UoWMock)(nil)
