package domain

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	orderDomain "order/internal/domain/order"
	"testing"
	"time"
)

type OrderDomainTestSuite struct {
	suite.Suite
}

func (s *OrderDomainTestSuite) TestNoteCanceledByCustomer() {
	tests := []struct {
		name           string
		setup          func(ord *orderDomain.Order) error
		expectedStatus orderDomain.Status
		expectError    bool
	}{
		{
			name: "Success: Order in Delivering",
			setup: func(ord *orderDomain.Order) error {
				return ord.NoteDelivering(uuid.New())
			},
			expectedStatus: orderDomain.CustomerCanceled,
			expectError:    false,
		},
		{
			name: "Failure: Order remains in Created",
			setup: func(ord *orderDomain.Order) error {
				return nil
			},
			expectedStatus: orderDomain.Created,
			expectError:    true,
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			ord := orderDomain.Create(uuid.New(), "Test Address", []orderDomain.Item{})
			err := tc.setup(ord)
			require.NoError(s.T(), err)

			err = ord.NoteCanceledByCustomer()

			if tc.expectError {
				require.Error(s.T(), err)
				require.ErrorIs(s.T(), err, orderDomain.ErrUnsupportedStatusTransition)
			} else {
				require.NoError(s.T(), err)
			}
			require.Equal(s.T(), tc.expectedStatus, ord.Status)
		})
	}
}

func (s *OrderDomainTestSuite) TestNoteCanceledOutOfStock() {
	tests := []struct {
		name           string
		setup          func(ord *orderDomain.Order) error
		expectedStatus orderDomain.Status
		expectError    bool
	}{
		{
			name: "Success: Order in Created (default)",
			setup: func(ord *orderDomain.Order) error {
				return nil
			},
			expectedStatus: orderDomain.CanceledOutOfStock,
			expectError:    false,
		},
		{
			name: "Failure: Order in Delivering",
			setup: func(ord *orderDomain.Order) error {
				return ord.NoteDelivering(uuid.New())
			},
			expectedStatus: orderDomain.Delivering,
			expectError:    true,
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			ord := orderDomain.Create(uuid.New(), "Test Address", []orderDomain.Item{})
			err := tc.setup(ord)
			require.NoError(s.T(), err)

			err = ord.NoteCanceledOutOfStock()

			if tc.expectError {
				require.Error(s.T(), err)
				require.ErrorIs(s.T(), err, orderDomain.ErrUnsupportedStatusTransition)
			} else {
				require.NoError(s.T(), err)
			}
			require.Equal(s.T(), tc.expectedStatus, ord.Status)
		})
	}
}

func (s *OrderDomainTestSuite) TestNoteCanceledCourierNotFound() {
	tests := []struct {
		name           string
		setup          func(ord *orderDomain.Order) error
		expectedStatus orderDomain.Status
		expectError    bool
	}{
		{
			name: "Success: Order in Created (default)",
			setup: func(ord *orderDomain.Order) error {
				return nil
			},
			expectedStatus: orderDomain.CanceledCourierNotFound,
			expectError:    false,
		},
		{
			name: "Failure: Order in Delivering",
			setup: func(ord *orderDomain.Order) error {
				return ord.NoteDelivering(uuid.New())
			},
			expectedStatus: orderDomain.Delivering,
			expectError:    true,
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			ord := orderDomain.Create(uuid.New(), "Test Address", []orderDomain.Item{})
			err := tc.setup(ord)
			require.NoError(s.T(), err)

			err = ord.NoteCanceledCourierNotFound()

			if tc.expectError {
				require.Error(s.T(), err)
				require.ErrorIs(s.T(), err, orderDomain.ErrUnsupportedStatusTransition)
			} else {
				require.NoError(s.T(), err)
			}
			require.Equal(s.T(), tc.expectedStatus, ord.Status)
		})
	}
}

func (s *OrderDomainTestSuite) TestNoteDelivering() {
	tests := []struct {
		name           string
		setup          func(ord *orderDomain.Order) error
		courierID      uuid.UUID
		expectedStatus orderDomain.Status
		expectError    bool
	}{
		{
			name: "Success: Order in Created",
			setup: func(ord *orderDomain.Order) error {
				return nil
			},
			courierID:      uuid.New(),
			expectedStatus: orderDomain.Delivering,
			expectError:    false,
		},
		{
			name: "Failure: Order already in Delivering",
			setup: func(ord *orderDomain.Order) error {
				return ord.NoteDelivering(uuid.New())
			},
			courierID:      uuid.New(),
			expectedStatus: orderDomain.Delivering,
			expectError:    true,
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			ord := orderDomain.Create(uuid.New(), "Test Address", []orderDomain.Item{})
			err := tc.setup(ord)
			require.NoError(s.T(), err, "setup должен выполниться без ошибки")

			err = ord.NoteDelivering(tc.courierID)

			if tc.expectError {
				require.Error(s.T(), err)
				require.ErrorIs(s.T(), err, orderDomain.ErrUnsupportedStatusTransition)
			} else {
				require.NoError(s.T(), err)
			}
			require.Equal(s.T(), tc.expectedStatus, ord.Status)
			if !tc.expectError {
				require.NotNil(s.T(), ord.Delivery.CourierID)
				require.Equal(s.T(), tc.courierID, *ord.Delivery.CourierID)
			}
		})
	}
}

func (s *OrderDomainTestSuite) TestNoteDelivered() {
	tests := []struct {
		name           string
		setup          func(ord *orderDomain.Order) error
		expectedStatus orderDomain.Status
		expectError    bool
	}{
		{
			name: "Success: Order in Delivering",
			setup: func(ord *orderDomain.Order) error {
				return ord.NoteDelivering(uuid.New())
			},
			expectedStatus: orderDomain.Delivered,
			expectError:    false,
		},
		{
			name: "Failure: Order in Created (default)",
			setup: func(ord *orderDomain.Order) error {
				return nil
			},
			expectedStatus: orderDomain.Created,
			expectError:    true,
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			ord := orderDomain.Create(uuid.New(), "Test Address", []orderDomain.Item{})
			err := tc.setup(ord)
			require.NoError(s.T(), err, "setup должен выполниться без ошибки")

			err = ord.NoteDelivered()

			if tc.expectError {
				require.Error(s.T(), err)
				require.ErrorIs(s.T(), err, orderDomain.ErrUnsupportedStatusTransition)
			} else {
				require.NoError(s.T(), err)
			}
			require.Equal(s.T(), tc.expectedStatus, ord.Status)
			if !tc.expectError {
				require.NotNil(s.T(), ord.Delivery.Arrived)
				require.WithinDuration(s.T(), time.Now(), *ord.Delivery.Arrived, time.Second)
			}
		})
	}
}

func TestOrderDomainTestSuite(t *testing.T) {
	suite.Run(t, new(OrderDomainTestSuite))
}
