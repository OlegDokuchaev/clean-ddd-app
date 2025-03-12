package saga

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	createOrder "order/internal/application/order/saga/create_order"
	createOrderMock "order/internal/mocks/order/saga/create_order"
	"testing"
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
		setup       func(publisher *createOrderMock.PublisherMock)
		expectedErr error
	}{
		{
			name: "Success",
			event: createOrder.ItemsReserved{
				Event: createOrder.Event{
					ID: uuid.New(),
				},
				OrderID: uuid.New(),
			},
			setup: func(publisher *createOrderMock.PublisherMock) {
				publisher.On("PublishAssignCourierCmd", s.ctx, mock.Anything).
					Return(nil).Once()
			},
			expectedErr: nil,
		},
		{
			name: "Failure: publisher error",
			event: createOrder.ItemsReserved{
				Event: createOrder.Event{
					ID: uuid.New(),
				},
				OrderID: uuid.New(),
			},
			setup: func(publisher *createOrderMock.PublisherMock) {
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
			uc := createOrder.New(publisher)
			tc.setup(publisher)

			err := uc.HandleItemsReserved(s.ctx, tc.event)

			if tc.expectedErr == nil {
				require.NoError(s.T(), err)
			} else {
				require.Error(s.T(), err)
				require.EqualError(s.T(), err, tc.expectedErr.Error())
			}

			publisher.AssertExpectations(s.T())
		})
	}
}

func (s *CreateOrderSagaTestSuite) TestHandleItemsReservationFailed() {
	tests := []struct {
		name        string
		event       createOrder.ItemsReservationFailed
		setup       func(publisher *createOrderMock.PublisherMock)
		expectedErr error
	}{
		{
			name: "Success",
			event: createOrder.ItemsReservationFailed{
				Event: createOrder.Event{
					ID: uuid.New(),
				},
				OrderID: uuid.New(),
			},
			setup: func(publisher *createOrderMock.PublisherMock) {
				publisher.On("PublishCancelOutOfStockCmd", s.ctx, mock.Anything).
					Return(nil).Once()
			},
			expectedErr: nil,
		},
		{
			name: "Failure: publisher error",
			event: createOrder.ItemsReservationFailed{
				Event: createOrder.Event{
					ID: uuid.New(),
				},
				OrderID: uuid.New(),
			},
			setup: func(publisher *createOrderMock.PublisherMock) {
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
			uc := createOrder.New(publisher)
			tc.setup(publisher)

			err := uc.HandleItemsReservationFailed(s.ctx, tc.event)

			if tc.expectedErr == nil {
				require.NoError(s.T(), err)
			} else {
				require.Error(s.T(), err)
				require.EqualError(s.T(), err, tc.expectedErr.Error())
			}

			publisher.AssertExpectations(s.T())
		})
	}
}

func (s *CreateOrderSagaTestSuite) TestHandleCourierAssignmentFailed() {
	tests := []struct {
		name        string
		event       createOrder.CourierAssignmentFailed
		setup       func(publisher *createOrderMock.PublisherMock)
		expectedErr error
	}{
		{
			name: "Success",
			event: createOrder.CourierAssignmentFailed{
				Event: createOrder.Event{
					ID: uuid.New(),
				},
				OrderID: uuid.New(),
			},
			setup: func(publisher *createOrderMock.PublisherMock) {
				publisher.On("PublishReleaseItemsCmd", s.ctx, mock.Anything).
					Return(nil).Once()
			},
			expectedErr: nil,
		},
		{
			name: "Failure: publisher error",
			event: createOrder.CourierAssignmentFailed{
				Event: createOrder.Event{
					ID: uuid.New(),
				},
				OrderID: uuid.New(),
			},
			setup: func(publisher *createOrderMock.PublisherMock) {
				publisher.On("PublishReleaseItemsCmd", s.ctx, mock.Anything).
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
			uc := createOrder.New(publisher)
			tc.setup(publisher)

			err := uc.HandleCourierAssignmentFailed(s.ctx, tc.event)

			if tc.expectedErr == nil {
				require.NoError(s.T(), err)
			} else {
				require.Error(s.T(), err)
				require.EqualError(s.T(), err, tc.expectedErr.Error())
			}

			publisher.AssertExpectations(s.T())
		})
	}
}

func (s *CreateOrderSagaTestSuite) TestHandleItemsReleased() {
	tests := []struct {
		name        string
		event       createOrder.ItemsReleased
		setup       func(publisher *createOrderMock.PublisherMock)
		expectedErr error
	}{
		{
			name: "Success",
			event: createOrder.ItemsReleased{
				Event: createOrder.Event{
					ID: uuid.New(),
				},
				OrderID: uuid.New(),
			},
			setup: func(publisher *createOrderMock.PublisherMock) {
				publisher.On("PublishCancelCourierNotFoundCmd", s.ctx, mock.Anything).
					Return(nil).Once()
			},
			expectedErr: nil,
		},
		{
			name: "Failure: publisher error",
			event: createOrder.ItemsReleased{
				Event: createOrder.Event{
					ID: uuid.New(),
				},
				OrderID: uuid.New(),
			},
			setup: func(publisher *createOrderMock.PublisherMock) {
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
			uc := createOrder.New(publisher)
			tc.setup(publisher)

			err := uc.HandleItemsReleased(s.ctx, tc.event)

			if tc.expectedErr == nil {
				require.NoError(s.T(), err)
			} else {
				require.Error(s.T(), err)
				require.EqualError(s.T(), err, tc.expectedErr.Error())
			}

			publisher.AssertExpectations(s.T())
		})
	}
}

func (s *CreateOrderSagaTestSuite) TestHandleCourierAssigned() {
	tests := []struct {
		name        string
		event       createOrder.CourierAssigned
		setup       func(publisher *createOrderMock.PublisherMock)
		expectedErr error
	}{
		{
			name: "Success",
			event: createOrder.CourierAssigned{
				Event: createOrder.Event{
					ID: uuid.New(),
				},
				OrderID:   uuid.New(),
				CourierID: uuid.New(),
			},
			setup: func(publisher *createOrderMock.PublisherMock) {
				publisher.On("PublishBeginDeliveryCmd", s.ctx, mock.Anything).
					Return(nil).Once()
			},
			expectedErr: nil,
		},
		{
			name: "Failure: publisher error",
			event: createOrder.CourierAssigned{
				Event: createOrder.Event{
					ID: uuid.New(),
				},
				OrderID:   uuid.New(),
				CourierID: uuid.New(),
			},
			setup: func(publisher *createOrderMock.PublisherMock) {
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
			uc := createOrder.New(publisher)
			tc.setup(publisher)

			err := uc.HandleCourierAssigned(s.ctx, tc.event)

			if tc.expectedErr == nil {
				require.NoError(s.T(), err)
			} else {
				require.Error(s.T(), err)
				require.EqualError(s.T(), err, tc.expectedErr.Error())
			}

			publisher.AssertExpectations(s.T())
		})
	}
}

func TestCreateOrderSagaTestSuite(t *testing.T) {
	suite.Run(t, new(CreateOrderSagaTestSuite))
}
