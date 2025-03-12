//go:build integration

package repository

import (
	"context"
	orderDomain "order/internal/domain/order"
	"order/internal/infrastructure/db/migrations"
	orderRepository "order/internal/infrastructure/repository/order"
	"order/internal/tests/testutils"
	"testing"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type OrderRepositoryTestSuite struct {
	suite.Suite
	ctx    context.Context
	testDB *testutils.TestDB
}

func (s *OrderRepositoryTestSuite) SetupSuite() {
	config, err := migrations.NewConfig()
	require.NoError(s.T(), err)

	s.ctx = context.Background()

	s.testDB, err = testutils.NewTestDB(s.ctx, config)
	require.NoError(s.T(), err)
}

func (s *OrderRepositoryTestSuite) TearDownSuite() {
	if s.testDB != nil {
		err := s.testDB.Close(s.ctx)
		require.NoError(s.T(), err)
	}
}

func (s *OrderRepositoryTestSuite) getRepo() orderDomain.Repository {
	return orderRepository.New(s.testDB.DB)
}

func (s *OrderRepositoryTestSuite) createTestOrder() *orderDomain.Order {
	items := []orderDomain.Item{
		{
			ProductID: uuid.New(),
			Price:     decimal.NewFromInt(100),
			Count:     1,
		},
	}
	order, err := orderDomain.Create(uuid.New(), "Some Address", items)
	require.NoError(s.T(), err)
	return order
}

func (s *OrderRepositoryTestSuite) TestCreate() {
	tests := []struct {
		name          string
		setup         func(repo orderDomain.Repository) *orderDomain.Order
		expectedError error
	}{
		{
			name: "Success",
			setup: func(_ orderDomain.Repository) *orderDomain.Order {
				return s.createTestOrder()
			},
			expectedError: nil,
		},
		{
			name: "Failure: Order already exists",
			setup: func(repo orderDomain.Repository) *orderDomain.Order {
				order := s.createTestOrder()
				err := repo.Create(s.ctx, order)
				require.NoError(s.T(), err)
				return order
			},
			expectedError: orderRepository.ErrOrderAlreadyExists,
		},
	}

	repo := s.getRepo()
	for _, tc := range tests {
		s.Run(tc.name, func() {
			order := tc.setup(repo)
			err := repo.Create(s.ctx, order)

			if tc.expectedError != nil {
				require.Error(s.T(), err)
				require.Equal(s.T(), tc.expectedError, err)
			} else {
				require.NoError(s.T(), err)

				createdOrder, err := repo.GetByID(s.ctx, order.ID)
				require.NoError(s.T(), err)
				require.NotNil(s.T(), createdOrder)
				require.Equal(s.T(), order.ID, createdOrder.ID)
			}
		})
	}
}

func (s *OrderRepositoryTestSuite) TestUpdate() {
	tests := []struct {
		name          string
		setup         func(repo orderDomain.Repository) *orderDomain.Order
		update        func(order *orderDomain.Order)
		verify        func(updated, expected *orderDomain.Order)
		expectedError error
	}{
		{
			name: "Success: Update status",
			setup: func(repo orderDomain.Repository) *orderDomain.Order {
				order := s.createTestOrder()
				err := repo.Create(s.ctx, order)
				require.NoError(s.T(), err)
				return order
			},
			update: func(order *orderDomain.Order) {
				err := order.NoteDelivering(uuid.New())
				require.NoError(s.T(), err)
			},
			verify: func(updated, expected *orderDomain.Order) {
				require.Equal(s.T(), orderDomain.Delivering, updated.Status)
				require.Equal(s.T(), expected.ID, updated.ID)
			},
			expectedError: nil,
		},
		{
			name: "Success: Update delivery",
			setup: func(repo orderDomain.Repository) *orderDomain.Order {
				order := s.createTestOrder()
				err := repo.Create(s.ctx, order)
				require.NoError(s.T(), err)
				return order
			},
			update: func(order *orderDomain.Order) {
				err := order.NoteDelivering(uuid.New())
				require.NoError(s.T(), err)
			},
			verify: func(updated, expected *orderDomain.Order) {
				require.Equal(s.T(), expected.ID, updated.ID)
				require.Equal(s.T(), expected.Delivery.CourierID, updated.Delivery.CourierID)
			},
			expectedError: nil,
		},
		{
			name: "Failure: Order not found",
			setup: func(repo orderDomain.Repository) *orderDomain.Order {
				return s.createTestOrder()
			},
			update:        func(order *orderDomain.Order) {},
			verify:        func(updated, expected *orderDomain.Order) {},
			expectedError: orderRepository.ErrOrderNotFound,
		},
	}

	repo := s.getRepo()
	for _, tc := range tests {
		s.Run(tc.name, func() {
			order := tc.setup(repo)
			tc.update(order)

			err := repo.Update(s.ctx, order)

			if tc.expectedError != nil {
				require.Error(s.T(), err)
				require.Equal(s.T(), tc.expectedError, err)
			} else {
				require.NoError(s.T(), err)
				updatedOrder, err := repo.GetByID(s.ctx, order.ID)
				require.NoError(s.T(), err)
				require.NotNil(s.T(), updatedOrder)
				tc.verify(updatedOrder, order)
			}
		})
	}
}

func (s *OrderRepositoryTestSuite) TestGetByID() {
	tests := []struct {
		name          string
		setup         func(repo orderDomain.Repository) uuid.UUID
		expectedError error
	}{
		{
			name: "Success",
			setup: func(repo orderDomain.Repository) uuid.UUID {
				order := s.createTestOrder()
				err := repo.Create(s.ctx, order)
				require.NoError(s.T(), err)
				return order.ID
			},
		},
		{
			name: "Failure: Order not found",
			setup: func(_ orderDomain.Repository) uuid.UUID {
				return uuid.New()
			},
			expectedError: orderRepository.ErrOrderNotFound,
		},
	}

	repo := s.getRepo()
	for _, tc := range tests {
		s.Run(tc.name, func() {
			orderID := tc.setup(repo)

			order, err := repo.GetByID(s.ctx, orderID)

			if tc.expectedError != nil {
				require.Error(s.T(), err)
				require.Equal(s.T(), tc.expectedError, err)
			} else {
				require.NoError(s.T(), err)
				require.NotNil(s.T(), order)
				require.Equal(s.T(), orderID, order.ID)
			}
		})
	}
}

func (s *OrderRepositoryTestSuite) TestGetAllByCustomer() {
	tests := []struct {
		name          string
		customerID    uuid.UUID
		setup         func(repo orderDomain.Repository, customerID uuid.UUID) []uuid.UUID
		expectedError error
	}{
		{
			name:       "Success: One order",
			customerID: uuid.New(),
			setup: func(repo orderDomain.Repository, customerID uuid.UUID) []uuid.UUID {
				order := s.createTestOrder()
				order.CustomerID = customerID
				err := repo.Create(s.ctx, order)
				require.NoError(s.T(), err)
				return []uuid.UUID{order.ID}
			},
			expectedError: nil,
		},
		{
			name:       "Success: Multiple orders",
			customerID: uuid.New(),
			setup: func(repo orderDomain.Repository, customerID uuid.UUID) []uuid.UUID {
				var orderIDs []uuid.UUID
				for i := 0; i < 3; i++ {
					order := s.createTestOrder()
					order.CustomerID = customerID
					err := repo.Create(s.ctx, order)
					require.NoError(s.T(), err)
					orderIDs = append(orderIDs, order.ID)
				}
				return orderIDs
			},
			expectedError: nil,
		},
		{
			name:       "Success: No orders",
			customerID: uuid.New(),
			setup: func(_ orderDomain.Repository, _ uuid.UUID) []uuid.UUID {
				return []uuid.UUID{}
			},
			expectedError: nil,
		},
	}

	repo := s.getRepo()
	for _, tc := range tests {
		s.Run(tc.name, func() {
			orderIDs := tc.setup(repo, tc.customerID)

			orders, err := repo.GetAllByCustomer(s.ctx, tc.customerID)

			if tc.expectedError != nil {
				require.Error(s.T(), err)
				require.Equal(s.T(), tc.expectedError, err)
			} else {
				require.NoError(s.T(), err)
				require.Equal(s.T(), len(orderIDs), len(orders))
				for _, order := range orders {
					require.Contains(s.T(), orderIDs, order.ID)
				}
			}
		})
	}
}

func (s *OrderRepositoryTestSuite) TestGetCurrentByCourier() {
	tests := []struct {
		name          string
		courierID     uuid.UUID
		setup         func(repo orderDomain.Repository, courierID uuid.UUID) []uuid.UUID
		expectedError error
	}{
		{
			name:      "Success: One order",
			courierID: uuid.New(),
			setup: func(repo orderDomain.Repository, courierID uuid.UUID) []uuid.UUID {
				order := s.createTestOrder()
				err := order.NoteDelivering(courierID)
				require.NoError(s.T(), err)

				err = repo.Create(s.ctx, order)
				require.NoError(s.T(), err)

				return []uuid.UUID{order.ID}
			},
		},
		{
			name:      "Success: Multiple orders",
			courierID: uuid.New(),
			setup: func(repo orderDomain.Repository, courierID uuid.UUID) []uuid.UUID {
				var orderIDs []uuid.UUID
				for i := 0; i < 3; i++ {
					order := s.createTestOrder()
					err := order.NoteDelivering(courierID)
					require.NoError(s.T(), err)

					err = repo.Create(s.ctx, order)
					require.NoError(s.T(), err)

					orderIDs = append(orderIDs, order.ID)
				}
				return orderIDs
			},
		},
		{
			name:      "Success: No orders",
			courierID: uuid.New(),
			setup: func(_ orderDomain.Repository, _ uuid.UUID) []uuid.UUID {
				return []uuid.UUID{}
			},
			expectedError: nil,
		},
	}

	repo := s.getRepo()
	for _, tc := range tests {
		s.Run(tc.name, func() {
			orderIDs := tc.setup(repo, tc.courierID)

			orders, err := repo.GetCurrentByCourier(s.ctx, tc.courierID)

			if tc.expectedError != nil {
				require.Error(s.T(), err)
				require.Equal(s.T(), tc.expectedError, err)
			} else {
				require.NoError(s.T(), err)
				require.Equal(s.T(), len(orderIDs), len(orders))
				for _, order := range orders {
					require.Contains(s.T(), orderIDs, order.ID)
				}
			}
		})
	}
}

func TestOrderRepository(t *testing.T) {
	suite.Run(t, new(OrderRepositoryTestSuite))
}
