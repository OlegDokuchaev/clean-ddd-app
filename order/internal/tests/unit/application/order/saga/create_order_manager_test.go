package saga

import (
	"context"
	"errors"
	createOrder "order/internal/application/order/saga/create_order"
	orderDomain "order/internal/domain/order"
	createOrderMock "order/internal/mocks/order/saga/create_order"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
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
			name: "Success",
			order: &orderDomain.Order{
				ID:         uuid.New(),
				CustomerID: uuid.New(),
				Status:     orderDomain.Created,
				Created:    time.Now(),
				Version:    uuid.New(),
				Items: []orderDomain.Item{
					{
						ProductID: uuid.New(),
						Price:     decimal.NewFromInt(100),
						Count:     2,
					},
				},
			},
			setup: func(publisher *createOrderMock.PublisherMock) {
				publisher.On("PublishReserveItemsCmd", mock.Anything, mock.Anything).Return(nil).Once()
			},
			expectedErr: nil,
		},
		{
			name: "Failure: Publisher error",
			order: &orderDomain.Order{
				ID:         uuid.New(),
				CustomerID: uuid.New(),
				Status:     orderDomain.Created,
				Created:    time.Now(),
				Version:    uuid.New(),
				Items: []orderDomain.Item{
					{
						ProductID: uuid.New(),
						Price:     decimal.NewFromInt(100),
						Count:     2,
					},
				},
			},
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
