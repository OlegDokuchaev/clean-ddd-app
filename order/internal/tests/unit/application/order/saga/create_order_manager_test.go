package saga

import (
	"context"
	"errors"
	createOrder "order/internal/application/order/saga/create_order"
	orderDomain "order/internal/domain/order"
	createOrderMock "order/internal/mocks/order/saga/create_order"
	"order/internal/tests/testutils/mothers"
	"testing"

	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"github.com/stretchr/testify/mock"
)

type CreateOrderSagaManagerTestSuite struct {
	suite.Suite
	ctx context.Context
}

func (s *CreateOrderSagaManagerTestSuite) BeforeEach(_ provider.T) {
	s.ctx = context.Background()
}

func (s *CreateOrderSagaManagerTestSuite) TestCreate(t provider.T) {
	t.Parallel()

	tests := []struct {
		name        string
		order       *orderDomain.Order
		setup       func(publisher *createOrderMock.PublisherMock)
		expectedErr error
	}{
		{
			name:  "Success",
			order: mothers.DefaultOrder(),
			setup: func(publisher *createOrderMock.PublisherMock) {
				publisher.On("PublishReserveItemsCmd", mock.Anything, mock.Anything).Return(nil).Once()
			},
			expectedErr: nil,
		},
		{
			name:  "Failure: Publisher error",
			order: mothers.DefaultOrder(),
			setup: func(publisher *createOrderMock.PublisherMock) {
				publisher.On("PublishReserveItemsCmd", mock.Anything, mock.Anything).
					Return(errors.New("publisher error")).Once()
			},
			expectedErr: nil,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t provider.T) {
			t.Parallel()

			publisher := new(createOrderMock.PublisherMock)
			manager := createOrder.NewManager(publisher)
			tc.setup(publisher)

			manager.Create(s.ctx, tc.order)

			publisher.AssertExpectations(t)
		})
	}
}

func TestCreateOrderSagaManagerTestSuite(t *testing.T) {
	suite.RunSuite(t, new(CreateOrderSagaManagerTestSuite))
}
