package create_order

import (
	"context"
	"github.com/stretchr/testify/mock"
	orderDomain "order/internal/domain/order"
)

type ManagerMock struct {
	mock.Mock
}

func (m *ManagerMock) Create(ctx context.Context, order *orderDomain.Order) {
	m.Called(ctx, order)
}
