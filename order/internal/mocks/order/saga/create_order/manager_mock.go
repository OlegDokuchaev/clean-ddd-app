package create_order

import (
	"context"
	"github.com/stretchr/testify/mock"
	createOrder "order/internal/application/order/saga/create_order"
	orderDomain "order/internal/domain/order"
)

type ManagerMock struct {
	mock.Mock
}

func (m *ManagerMock) Create(ctx context.Context, order *orderDomain.Order) {
	m.Called(ctx, order)
}

var _ createOrder.Manager = (*ManagerMock)(nil)
