package saga

import (
	"context"
	"errors"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	createOrder "order/internal/application/order/saga/create_order"
	orderDomain "order/internal/domain/order"
	createOrderMock "order/internal/mocks/order/saga/create_order"
	"order/internal/tests/testutils/mothers"
	"testing"
)

type CreateOrderSagaManagerTestSuite struct {
	suite.Suite
	ctx context.Context
}

func (s *CreateOrderSagaManagerTestSuite) SetupTest() {
	s.ctx = context.Background()
}

func (s *CreateOrderSagaManagerTestSuite) TestCreate() {
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
		s.Run(tc.name, func() {
			s.T().Parallel()
			publisher := new(createOrderMock.PublisherMock)
			manager := createOrder.NewManager(publisher)
			tc.setup(publisher)

			manager.Create(s.ctx, tc.order)

			publisher.AssertExpectations(s.T())
		})
	}
}

func TestCreateOrderSagaManagerTestSuite(t *testing.T) {
	suite.Run(t, new(CreateOrderSagaManagerTestSuite))
}
