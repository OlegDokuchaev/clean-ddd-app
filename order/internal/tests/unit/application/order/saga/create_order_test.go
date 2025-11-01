package saga

import (
	"context"
	"errors"
	createOrder "order/internal/application/order/saga/create_order"
	orderDomain "order/internal/domain/order"
	orderMock "order/internal/mocks/order"
	createOrderMock "order/internal/mocks/order/saga/create_order"
	"order/internal/tests/testutils/mothers"
	"testing"

	"github.com/ozontech/allure-go/pkg/framework/provider"

	"github.com/google/uuid"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"github.com/stretchr/testify/mock"
)

type CreateOrderSagaTestSuite struct {
	suite.Suite
	ctx context.Context
}

func (s *CreateOrderSagaTestSuite) BeforeEach(_ provider.T) {
	s.ctx = context.Background()
}

func (s *CreateOrderSagaTestSuite) TestHandleItemsReserved(t provider.T) {
	t.Parallel()

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
		t.Run(tc.name, func(t provider.T) {
			t.Parallel()

			publisher := new(createOrderMock.PublisherMock)
			repository := new(orderMock.RepositoryMock)
			uc := createOrder.New(publisher, repository)
			tc.setup(publisher, repository)

			err := uc.HandleItemsReserved(s.ctx, tc.event)

			if tc.expectedErr == nil {
				t.Require().NoError(err)
			} else {
				t.Require().Error(err)
				t.Require().EqualError(err, tc.expectedErr.Error())
			}

			publisher.AssertExpectations(t)
			repository.AssertExpectations(t)
		})
	}
}

func (s *CreateOrderSagaTestSuite) TestHandleItemsReservationFailed(t provider.T) {
	t.Parallel()

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
		t.Run(tc.name, func(t provider.T) {
			t.Parallel()

			publisher := new(createOrderMock.PublisherMock)
			repository := new(orderMock.RepositoryMock)
			uc := createOrder.New(publisher, repository)
			tc.setup(publisher, repository)

			err := uc.HandleItemsReservationFailed(s.ctx, tc.event)

			if tc.expectedErr == nil {
				t.Require().NoError(err)
			} else {
				t.Require().Error(err)
				t.Require().EqualError(err, tc.expectedErr.Error())
			}

			publisher.AssertExpectations(t)
			repository.AssertExpectations(t)
		})
	}
}

func (s *CreateOrderSagaTestSuite) TestHandleCourierAssignmentFailed(t provider.T) {
	t.Parallel()

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
			order:   mothers.DefaultOrder(),
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
			order:   mothers.DefaultOrder(),
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
		t.Run(tc.name, func(t provider.T) {
			t.Parallel()

			publisher := new(createOrderMock.PublisherMock)
			repository := new(orderMock.RepositoryMock)
			tc.setup(publisher, repository, tc.order)
			uc := createOrder.New(publisher, repository)

			err := uc.HandleCourierAssignmentFailed(s.ctx, tc.event)

			if tc.expectedErr == nil {
				t.Require().NoError(err)
			} else {
				t.Require().Error(err)
				t.Require().EqualError(err, tc.expectedErr.Error())
			}

			publisher.AssertExpectations(t)
			repository.AssertExpectations(t)
		})
	}
}

func (s *CreateOrderSagaTestSuite) TestHandleItemsReleased(t provider.T) {
	t.Parallel()

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
		t.Run(tc.name, func(t provider.T) {
			t.Parallel()

			publisher := new(createOrderMock.PublisherMock)
			repository := new(orderMock.RepositoryMock)
			uc := createOrder.New(publisher, repository)
			tc.setup(publisher, repository)

			err := uc.HandleItemsReleased(s.ctx, tc.event)

			if tc.expectedErr == nil {
				t.Require().NoError(err)
			} else {
				t.Require().Error(err)
				t.Require().EqualError(err, tc.expectedErr.Error())
			}

			publisher.AssertExpectations(t)
			repository.AssertExpectations(t)
		})
	}
}

func (s *CreateOrderSagaTestSuite) TestHandleCourierAssigned(t provider.T) {
	t.Parallel()

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
		t.Run(tc.name, func(t provider.T) {
			t.Parallel()

			publisher := new(createOrderMock.PublisherMock)
			repository := new(orderMock.RepositoryMock)
			uc := createOrder.New(publisher, repository)
			tc.setup(publisher, repository)

			err := uc.HandleCourierAssigned(s.ctx, tc.event)

			if tc.expectedErr == nil {
				t.Require().NoError(err)
			} else {
				t.Require().Error(err)
				t.Require().EqualError(err, tc.expectedErr.Error())
			}

			publisher.AssertExpectations(t)
			repository.AssertExpectations(t)
		})
	}
}

func TestCreateOrderSagaTestSuite(t *testing.T) {
	suite.RunSuite(t, new(CreateOrderSagaTestSuite))
}
