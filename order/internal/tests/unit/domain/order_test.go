package domain

import (
	orderDomain "order/internal/domain/order"
	"order/internal/tests/testutils/mothers"
	"testing"
	"time"

	"github.com/ozontech/allure-go/pkg/framework/provider"

	"github.com/google/uuid"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"github.com/shopspring/decimal"
)

type OrderDomainTestSuite struct {
	suite.Suite
}

func (s *OrderDomainTestSuite) TestCreate(t provider.T) {
	t.Parallel()

	tests := []struct {
		name        string
		CustomerID  uuid.UUID
		Address     string
		Items       []orderDomain.Item
		expectedErr error
	}{
		{
			name:       "Success",
			CustomerID: uuid.New(),
			Address:    "Test Address",
			Items: []orderDomain.Item{
				{
					ProductID: uuid.New(),
					Price:     decimal.NewFromInt(100),
					Count:     1,
				},
			},
			expectedErr: nil,
		},
		{
			name:        "Failure: Invalid items",
			CustomerID:  uuid.New(),
			Address:     "Test Address",
			Items:       []orderDomain.Item{},
			expectedErr: orderDomain.ErrInvalidItems,
		},
		{
			name:        "Failure: Invalid address",
			CustomerID:  uuid.New(),
			Address:     "",
			Items:       []orderDomain.Item{},
			expectedErr: orderDomain.ErrInvalidAddress,
		},
		{
			name:       "Failure: Invalid item price",
			CustomerID: uuid.New(),
			Address:    "Test Address",
			Items: []orderDomain.Item{
				{
					ProductID: uuid.New(),
					Price:     decimal.NewFromInt(-1),
					Count:     1,
				},
			},
			expectedErr: orderDomain.ErrInvalidItems,
		},
		{
			name:       "Failure: Invalid item count",
			CustomerID: uuid.New(),
			Address:    "Test Address",
			Items: []orderDomain.Item{
				{
					ProductID: uuid.New(),
					Price:     decimal.NewFromInt(100),
					Count:     0,
				},
			},
			expectedErr: orderDomain.ErrInvalidItems,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t provider.T) {
			t.Parallel()

			order, err := orderDomain.Create(tc.CustomerID, tc.Address, tc.Items)

			if tc.expectedErr != nil {
				t.Require().Error(err)
				t.Require().ErrorIs(err, tc.expectedErr)
			} else {
				t.Require().NoError(err)
				t.Require().NotNil(order)
			}
		})
	}
}

func (s *OrderDomainTestSuite) TestNoteCanceledByCustomer(t provider.T) {
	t.Parallel()

	tests := []struct {
		name           string
		setup          func() *orderDomain.Order
		expectedStatus orderDomain.Status
		expectedErr    error
	}{
		{
			name: "Success: Order in Delivering",
			setup: func() *orderDomain.Order {
				return mothers.OrderDelivering()
			},
			expectedStatus: orderDomain.CustomerCanceled,
			expectedErr:    nil,
		},
		{
			name: "Failure: Order remains in Created",
			setup: func() *orderDomain.Order {
				return mothers.DefaultOrder()
			},
			expectedStatus: orderDomain.Created,
			expectedErr:    orderDomain.ErrUnsupportedStatusTransition,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t provider.T) {
			t.Parallel()
			order := tc.setup()

			err := order.NoteCanceledByCustomer()

			if tc.expectedErr != nil {
				t.Require().Error(err)
				t.Require().ErrorIs(err, tc.expectedErr)
			} else {
				t.Require().NoError(err)
			}
			t.Require().Equal(tc.expectedStatus, order.Status)
		})
	}
}

func (s *OrderDomainTestSuite) TestNoteCanceledOutOfStock(t provider.T) {
	t.Parallel()

	tests := []struct {
		name           string
		setup          func() *orderDomain.Order
		expectedStatus orderDomain.Status
		expectedErr    error
	}{
		{
			name: "Success: Order in Created (default)",
			setup: func() *orderDomain.Order {
				return mothers.DefaultOrder()
			},
			expectedStatus: orderDomain.CanceledOutOfStock,
			expectedErr:    nil,
		},
		{
			name: "Failure: Order in Delivering",
			setup: func() *orderDomain.Order {
				return mothers.OrderDelivering()
			},
			expectedStatus: orderDomain.Delivering,
			expectedErr:    orderDomain.ErrUnsupportedStatusTransition,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t provider.T) {
			t.Parallel()
			order := tc.setup()

			err := order.NoteCanceledOutOfStock()

			if tc.expectedErr != nil {
				t.Require().Error(err)
				t.Require().ErrorIs(err, tc.expectedErr)
			} else {
				t.Require().NoError(err)
			}
			t.Require().Equal(tc.expectedStatus, order.Status)
		})
	}
}

func (s *OrderDomainTestSuite) TestNoteCanceledCourierNotFound(t provider.T) {
	t.Parallel()

	tests := []struct {
		name           string
		setup          func() *orderDomain.Order
		expectedStatus orderDomain.Status
		expectedErr    error
	}{
		{
			name: "Success: Order in Created (default)",
			setup: func() *orderDomain.Order {
				return mothers.DefaultOrder()
			},
			expectedStatus: orderDomain.CanceledCourierNotFound,
			expectedErr:    nil,
		},
		{
			name: "Failure: Order in Delivering",
			setup: func() *orderDomain.Order {
				return mothers.OrderDelivering()
			},
			expectedStatus: orderDomain.Delivering,
			expectedErr:    orderDomain.ErrUnsupportedStatusTransition,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t provider.T) {
			t.Parallel()
			order := tc.setup()

			err := order.NoteCanceledCourierNotFound()

			if tc.expectedErr != nil {
				t.Require().Error(err)
				t.Require().ErrorIs(err, tc.expectedErr)
			} else {
				t.Require().NoError(err)
			}
			t.Require().Equal(tc.expectedStatus, order.Status)
		})
	}
}

func (s *OrderDomainTestSuite) TestNoteDelivering(t provider.T) {
	t.Parallel()

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
				return mothers.DefaultOrder()
			},
			courierID:      uuid.New(),
			expectedStatus: orderDomain.Delivering,
			expectedErr:    nil,
		},
		{
			name: "Failure: Order already in Delivering",
			setup: func() *orderDomain.Order {
				return mothers.OrderDelivering()
			},
			courierID:      uuid.New(),
			expectedStatus: orderDomain.Delivering,
			expectedErr:    orderDomain.ErrUnsupportedStatusTransition,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t provider.T) {
			t.Parallel()
			order := tc.setup()

			err := order.NoteDelivering(tc.courierID)

			if tc.expectedErr != nil {
				t.Require().Error(err)
				t.Require().ErrorIs(err, tc.expectedErr)
			} else {
				t.Require().NoError(err)
				t.Require().NotNil(order.Delivery.CourierID)
				t.Require().Equal(tc.courierID, *order.Delivery.CourierID)
			}
			t.Require().Equal(tc.expectedStatus, order.Status)
		})
	}
}

func (s *OrderDomainTestSuite) TestNoteDelivered(t provider.T) {
	t.Parallel()

	tests := []struct {
		name           string
		setup          func() *orderDomain.Order
		expectedStatus orderDomain.Status
		expectedErr    error
	}{
		{
			name: "Success: Order in Delivering",
			setup: func() *orderDomain.Order {
				return mothers.OrderDelivering()
			},
			expectedStatus: orderDomain.Delivered,
			expectedErr:    nil,
		},
		{
			name: "Failure: Order in Created (default)",
			setup: func() *orderDomain.Order {
				return mothers.DefaultOrder()
			},
			expectedStatus: orderDomain.Created,
			expectedErr:    orderDomain.ErrUnsupportedStatusTransition,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t provider.T) {
			t.Parallel()
			order := tc.setup()

			err := order.NoteDelivered()

			if tc.expectedErr != nil {
				t.Require().Error(err)
				t.Require().ErrorIs(err, tc.expectedErr)
			} else {
				t.Require().NoError(err)
				t.Require().NotNil(order.Delivery.CourierID)
				t.Require().WithinDuration(time.Now(), *order.Delivery.Arrived, time.Second)
			}
			t.Require().Equal(tc.expectedStatus, order.Status)
		})
	}
}

func TestOrderDomainTestSuite(t *testing.T) {
	suite.RunSuite(t, new(OrderDomainTestSuite))
}
