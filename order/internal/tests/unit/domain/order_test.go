package domain

import (
	orderDomain "order/internal/domain/order"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type OrderDomainTestSuite struct {
	suite.Suite
}

func (s *OrderDomainTestSuite) createTestOrder() *orderDomain.Order {
	items := []orderDomain.Item{
		{
			ProductID: uuid.New(),
			Price:     decimal.NewFromInt(100),
		},
	}
	order, err := orderDomain.Create(uuid.New(), "Test Address", items)
	require.NoError(s.T(), err)
	return order
}

func (s *OrderDomainTestSuite) TestNoteCanceledByCustomer() {
	tests := []struct {
		name           string
		setup          func() *orderDomain.Order
		expectedStatus orderDomain.Status
		expectedErr    error
	}{
		{
			name: "Success: Order in Delivering",
			setup: func() *orderDomain.Order {
				ord := s.createTestOrder()
				require.NoError(s.T(), ord.NoteDelivering(uuid.New()))
				return ord
			},
			expectedStatus: orderDomain.CustomerCanceled,
			expectedErr:    nil,
		},
		{
			name: "Failure: Order remains in Created",
			setup: func() *orderDomain.Order {
				return s.createTestOrder()
			},
			expectedStatus: orderDomain.Created,
			expectedErr:    orderDomain.ErrUnsupportedStatusTransition,
		},
	}

	for _, tc := range tests {
		tc := tc
		s.Run(tc.name, func() {
			s.T().Parallel()
			order := tc.setup()

			err := order.NoteCanceledByCustomer()

			if tc.expectedErr != nil {
				require.Error(s.T(), err)
				require.ErrorIs(s.T(), err, tc.expectedErr)
			} else {
				require.NoError(s.T(), err)
			}
			require.Equal(s.T(), tc.expectedStatus, order.Status)
		})
	}
}

func (s *OrderDomainTestSuite) TestNoteCanceledOutOfStock() {
	tests := []struct {
		name           string
		setup          func() *orderDomain.Order
		expectedStatus orderDomain.Status
		expectedErr    error
	}{
		{
			name: "Success: Order in Created (default)",
			setup: func() *orderDomain.Order {
				return s.createTestOrder()
			},
			expectedStatus: orderDomain.CanceledOutOfStock,
			expectedErr:    nil,
		},
		{
			name: "Failure: Order in Delivering",
			setup: func() *orderDomain.Order {
				ord := s.createTestOrder()
				require.NoError(s.T(), ord.NoteDelivering(uuid.New()))
				return ord
			},
			expectedStatus: orderDomain.Delivering,
			expectedErr:    orderDomain.ErrUnsupportedStatusTransition,
		},
	}

	for _, tc := range tests {
		tc := tc
		s.Run(tc.name, func() {
			s.T().Parallel()
			order := tc.setup()

			err := order.NoteCanceledOutOfStock()

			if tc.expectedErr != nil {
				require.Error(s.T(), err)
				require.ErrorIs(s.T(), err, tc.expectedErr)
			} else {
				require.NoError(s.T(), err)
			}
			require.Equal(s.T(), tc.expectedStatus, order.Status)
		})
	}
}

func (s *OrderDomainTestSuite) TestNoteCanceledCourierNotFound() {
	tests := []struct {
		name           string
		setup          func() *orderDomain.Order
		expectedStatus orderDomain.Status
		expectedErr    error
	}{
		{
			name: "Success: Order in Created (default)",
			setup: func() *orderDomain.Order {
				return s.createTestOrder()
			},
			expectedStatus: orderDomain.CanceledCourierNotFound,
			expectedErr:    nil,
		},
		{
			name: "Failure: Order in Delivering",
			setup: func() *orderDomain.Order {
				ord := s.createTestOrder()
				require.NoError(s.T(), ord.NoteDelivering(uuid.New()))
				return ord
			},
			expectedStatus: orderDomain.Delivering,
			expectedErr:    orderDomain.ErrUnsupportedStatusTransition,
		},
	}

	for _, tc := range tests {
		tc := tc
		s.Run(tc.name, func() {
			s.T().Parallel()
			order := tc.setup()

			err := order.NoteCanceledCourierNotFound()

			if tc.expectedErr != nil {
				require.Error(s.T(), err)
				require.ErrorIs(s.T(), err, tc.expectedErr)
			} else {
				require.NoError(s.T(), err)
			}
			require.Equal(s.T(), tc.expectedStatus, order.Status)
		})
	}
}

func (s *OrderDomainTestSuite) TestNoteDelivering() {
	tests := []struct {
		name           string
		setup          func() *orderDomain.Order
		courierID      uuid.UUID
		expectedStatus orderDomain.Status
		expectedErr    error
	}{
		{
			name: "Success: Order in Created",
			setup: func() *orderDomain.Order {
				return s.createTestOrder()
			},
			courierID:      uuid.New(),
			expectedStatus: orderDomain.Delivering,
			expectedErr:    nil,
		},
		{
			name: "Failure: Order already in Delivering",
			setup: func() *orderDomain.Order {
				ord := s.createTestOrder()
				require.NoError(s.T(), ord.NoteDelivering(uuid.New()))
				return ord
			},
			courierID:      uuid.New(),
			expectedStatus: orderDomain.Delivering,
			expectedErr:    orderDomain.ErrUnsupportedStatusTransition,
		},
	}

	for _, tc := range tests {
		tc := tc
		s.Run(tc.name, func() {
			s.T().Parallel()
			order := tc.setup()

			err := order.NoteDelivering(tc.courierID)

			if tc.expectedErr != nil {
				require.Error(s.T(), err)
				require.ErrorIs(s.T(), err, tc.expectedErr)
			} else {
				require.NoError(s.T(), err)
				require.NotNil(s.T(), order.Delivery.CourierID)
				require.Equal(s.T(), tc.courierID, *order.Delivery.CourierID)
			}
			require.Equal(s.T(), tc.expectedStatus, order.Status)
		})
	}
}

func (s *OrderDomainTestSuite) TestNoteDelivered() {
	tests := []struct {
		name           string
		setup          func() *orderDomain.Order
		expectedStatus orderDomain.Status
		expectedErr    error
	}{
		{
			name: "Success: Order in Delivering",
			setup: func() *orderDomain.Order {
				ord := s.createTestOrder()
				require.NoError(s.T(), ord.NoteDelivering(uuid.New()))
				return ord
			},
			expectedStatus: orderDomain.Delivered,
			expectedErr:    nil,
		},
		{
			name: "Failure: Order in Created (default)",
			setup: func() *orderDomain.Order {
				return s.createTestOrder()
			},
			expectedStatus: orderDomain.Created,
			expectedErr:    orderDomain.ErrUnsupportedStatusTransition,
		},
	}

	for _, tc := range tests {
		tc := tc
		s.Run(tc.name, func() {
			s.T().Parallel()
			order := tc.setup()

			err := order.NoteDelivered()

			if tc.expectedErr != nil {
				require.Error(s.T(), err)
				require.ErrorIs(s.T(), err, tc.expectedErr)
			} else {
				require.NoError(s.T(), err)
				require.NotNil(s.T(), order.Delivery.Arrived)
				require.WithinDuration(s.T(), time.Now(), *order.Delivery.Arrived, time.Second)
			}
			require.Equal(s.T(), tc.expectedStatus, order.Status)
		})
	}
}

func TestOrderDomainTestSuite(t *testing.T) {
	suite.Run(t, new(OrderDomainTestSuite))
}
