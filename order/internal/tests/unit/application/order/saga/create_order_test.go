package saga

import (
	"context"
	"errors"
	createOrder "order/internal/application/order/saga/create_order"
	orderDomain "order/internal/domain/order"
	orderMock "order/internal/mocks/order"
	createOrderMock "order/internal/mocks/order/saga/create_order"
	"testing"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type CreateOrderSagaTestSuite struct {
	suite.Suite
	ctx context.Context
}

func (s *CreateOrderSagaTestSuite) SetupTest() {
	s.ctx = context.Background()
}

func (s *CreateOrderSagaTestSuite) TestHandleItemsReserved() {
	tests := []struct {
		name        string
		event       createOrder.ItemsReserved
		setup       func(publisher *createOrderMock.PublisherMock, repository *orderMock.RepositoryMock)
		expectedErr error
	}{
		{
			name: "Success",
			event: createOrder.ItemsReserved{
				OrderID: uuid.New(),
			},
			setup: func(publisher *createOrderMock.PublisherMock, _ *orderMock.RepositoryMock) {
				publisher.On("PublishAssignCourierCmd", s.ctx, mock.Anything).
					Return(nil).Once()
			},
			expectedErr: nil,
		},
		{
			name: "Failure: publisher error",
			event: createOrder.ItemsReserved{
				OrderID: uuid.New(),
			},
			setup: func(publisher *createOrderMock.PublisherMock, _ *orderMock.RepositoryMock) {
				publisher.On("PublishAssignCourierCmd", s.ctx, mock.Anything).
					Return(errors.New("publisher error")).Once()
			},
			expectedErr: errors.New("publisher error"),
		},
	}

	for _, tc := range tests {
		tc := tc
		s.Run(tc.name, func() {
			s.T().Parallel()
			publisher := new(createOrderMock.PublisherMock)
			repository := new(orderMock.RepositoryMock)
			uc := createOrder.New(publisher, repository)
			tc.setup(publisher, repository)

			err := uc.HandleItemsReserved(s.ctx, tc.event)

			if tc.expectedErr == nil {
				require.NoError(s.T(), err)
			} else {
				require.Error(s.T(), err)
				require.EqualError(s.T(), err, tc.expectedErr.Error())
			}

			publisher.AssertExpectations(s.T())
			repository.AssertExpectations(s.T())
		})
	}
}

func (s *CreateOrderSagaTestSuite) TestHandleItemsReservationFailed() {
	tests := []struct {
		name        string
		event       createOrder.ItemsReservationFailed
		setup       func(publisher *createOrderMock.PublisherMock, repository *orderMock.RepositoryMock)
		expectedErr error
	}{
		{
			name: "Success",
			event: createOrder.ItemsReservationFailed{
				OrderID: uuid.New(),
			},
			setup: func(publisher *createOrderMock.PublisherMock, _ *orderMock.RepositoryMock) {
				publisher.On("PublishCancelOutOfStockCmd", s.ctx, mock.Anything).
					Return(nil).Once()
			},
			expectedErr: nil,
		},
		{
			name: "Failure: publisher error",
			event: createOrder.ItemsReservationFailed{
				OrderID: uuid.New(),
			},
			setup: func(publisher *createOrderMock.PublisherMock, _ *orderMock.RepositoryMock) {
				publisher.On("PublishCancelOutOfStockCmd", s.ctx, mock.Anything).
					Return(errors.New("publisher error")).Once()
			},
			expectedErr: errors.New("publisher error"),
		},
	}

	for _, tc := range tests {
		tc := tc
		s.Run(tc.name, func() {
			s.T().Parallel()
			publisher := new(createOrderMock.PublisherMock)
			repository := new(orderMock.RepositoryMock)
			uc := createOrder.New(publisher, repository)
			tc.setup(publisher, repository)

			err := uc.HandleItemsReservationFailed(s.ctx, tc.event)

			if tc.expectedErr == nil {
				require.NoError(s.T(), err)
			} else {
				require.Error(s.T(), err)
				require.EqualError(s.T(), err, tc.expectedErr.Error())
			}

			publisher.AssertExpectations(s.T())
			repository.AssertExpectations(s.T())
		})
	}
}

func (s *CreateOrderSagaTestSuite) TestHandleCourierAssignmentFailed() {
	tests := []struct {
		name        string
		event       createOrder.CourierAssignmentFailed
		order       *orderDomain.Order
		repoErr     error
		setup       func(publisher *createOrderMock.PublisherMock, repository *orderMock.RepositoryMock, order *orderDomain.Order)
		expectedErr error
	}{
		{
			name: "Success",
			event: createOrder.CourierAssignmentFailed{
				OrderID: uuid.New(),
			},
			order: &orderDomain.Order{
				ID: uuid.New(),
				Items: []orderDomain.Item{
					{
						ProductID: uuid.New(),
						Price:     decimal.NewFromInt(100),
						Count:     2,
					},
				},
			},
			repoErr: nil,
			setup: func(publisher *createOrderMock.PublisherMock, repository *orderMock.RepositoryMock, order *orderDomain.Order) {
				publisher.On("PublishReleaseItemsCmd", s.ctx, mock.Anything).Return(nil).Once()
				repository.On("GetByID", s.ctx, mock.Anything).Return(order, nil).Once()
			},
			expectedErr: nil,
		},
		{
			name: "Failure: repository error",
			event: createOrder.CourierAssignmentFailed{
				OrderID: uuid.New(),
			},
			order:   nil,
			repoErr: errors.New("repository error"),
			setup: func(_ *createOrderMock.PublisherMock, repository *orderMock.RepositoryMock, _ *orderDomain.Order) {
				repository.On("GetByID", s.ctx, mock.Anything).
					Return((*orderDomain.Order)(nil), errors.New("repository error")).Once()
			},
			expectedErr: errors.New("repository error"),
		},
		{
			name: "Failure: publisher error",
			event: createOrder.CourierAssignmentFailed{
				OrderID: uuid.New(),
			},
			order: &orderDomain.Order{
				ID: uuid.New(),
				Items: []orderDomain.Item{
					{
						ProductID: uuid.New(),
						Price:     decimal.NewFromInt(100),
						Count:     2,
					},
				},
			},
			repoErr: nil,
			setup: func(publisher *createOrderMock.PublisherMock, repository *orderMock.RepositoryMock, order *orderDomain.Order) {
				publisher.On("PublishReleaseItemsCmd", s.ctx, mock.Anything).
					Return(errors.New("publisher error")).Once()
				repository.On("GetByID", s.ctx, mock.Anything).Return(order, nil).Once()
			},
			expectedErr: errors.New("publisher error"),
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			publisher := new(createOrderMock.PublisherMock)
			repository := new(orderMock.RepositoryMock)
			tc.setup(publisher, repository, tc.order)
			uc := createOrder.New(publisher, repository)

			err := uc.HandleCourierAssignmentFailed(s.ctx, tc.event)

			if tc.expectedErr == nil {
				require.NoError(s.T(), err)
			} else {
				require.Error(s.T(), err)
				require.EqualError(s.T(), err, tc.expectedErr.Error())
			}

			publisher.AssertExpectations(s.T())
			repository.AssertExpectations(s.T())
		})
	}
}

func (s *CreateOrderSagaTestSuite) TestHandleItemsReleased() {
	tests := []struct {
		name        string
		event       createOrder.ItemsReleased
		setup       func(publisher *createOrderMock.PublisherMock, repository *orderMock.RepositoryMock)
		expectedErr error
	}{
		{
			name: "Success",
			event: createOrder.ItemsReleased{
				OrderID: uuid.New(),
			},
			setup: func(publisher *createOrderMock.PublisherMock, _ *orderMock.RepositoryMock) {
				publisher.On("PublishCancelCourierNotFoundCmd", s.ctx, mock.Anything).
					Return(nil).Once()
			},
			expectedErr: nil,
		},
		{
			name: "Failure: publisher error",
			event: createOrder.ItemsReleased{
				OrderID: uuid.New(),
			},
			setup: func(publisher *createOrderMock.PublisherMock, _ *orderMock.RepositoryMock) {
				publisher.On("PublishCancelCourierNotFoundCmd", s.ctx, mock.Anything).
					Return(errors.New("publisher error")).Once()
			},
			expectedErr: errors.New("publisher error"),
		},
	}

	for _, tc := range tests {
		tc := tc
		s.Run(tc.name, func() {
			s.T().Parallel()
			publisher := new(createOrderMock.PublisherMock)
			repository := new(orderMock.RepositoryMock)
			uc := createOrder.New(publisher, repository)
			tc.setup(publisher, repository)

			err := uc.HandleItemsReleased(s.ctx, tc.event)

			if tc.expectedErr == nil {
				require.NoError(s.T(), err)
			} else {
				require.Error(s.T(), err)
				require.EqualError(s.T(), err, tc.expectedErr.Error())
			}

			publisher.AssertExpectations(s.T())
			repository.AssertExpectations(s.T())
		})
	}
}

func (s *CreateOrderSagaTestSuite) TestHandleCourierAssigned() {
	tests := []struct {
		name        string
		event       createOrder.CourierAssigned
		setup       func(publisher *createOrderMock.PublisherMock, repository *orderMock.RepositoryMock)
		expectedErr error
	}{
		{
			name: "Success",
			event: createOrder.CourierAssigned{
				OrderID:   uuid.New(),
				CourierID: uuid.New(),
			},
			setup: func(publisher *createOrderMock.PublisherMock, _ *orderMock.RepositoryMock) {
				publisher.On("PublishBeginDeliveryCmd", s.ctx, mock.Anything).
					Return(nil).Once()
			},
			expectedErr: nil,
		},
		{
			name: "Failure: publisher error",
			event: createOrder.CourierAssigned{
				OrderID:   uuid.New(),
				CourierID: uuid.New(),
			},
			setup: func(publisher *createOrderMock.PublisherMock, _ *orderMock.RepositoryMock) {
				publisher.On("PublishBeginDeliveryCmd", s.ctx, mock.Anything).
					Return(errors.New("publisher error")).Once()
			},
			expectedErr: errors.New("publisher error"),
		},
	}

	for _, tc := range tests {
		tc := tc
		s.Run(tc.name, func() {
			s.T().Parallel()
			publisher := new(createOrderMock.PublisherMock)
			repository := new(orderMock.RepositoryMock)
			uc := createOrder.New(publisher, repository)
			tc.setup(publisher, repository)

			err := uc.HandleCourierAssigned(s.ctx, tc.event)

			if tc.expectedErr == nil {
				require.NoError(s.T(), err)
			} else {
				require.Error(s.T(), err)
				require.EqualError(s.T(), err, tc.expectedErr.Error())
			}

			publisher.AssertExpectations(s.T())
			repository.AssertExpectations(s.T())
		})
	}
}

func TestCreateOrderSagaTestSuite(t *testing.T) {
	suite.Run(t, new(CreateOrderSagaTestSuite))
}
