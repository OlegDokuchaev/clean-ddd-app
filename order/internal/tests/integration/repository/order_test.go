//go:build integration

package repository

import (
	"context"
	orderDomain "order/internal/domain/order"
	"order/internal/infrastructure/db/migrations"
	orderRepository "order/internal/infrastructure/repository/order"
	"order/internal/tests/testutils"
	"order/internal/tests/testutils/mothers"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"go.mongodb.org/mongo-driver/mongo"
)

type OrderRepositoryTestSuite struct {
	suite.Suite

	ctx context.Context

	db              *testutils.TestDB
	orderCollection *mongo.Collection
}

func (s *OrderRepositoryTestSuite) BeforeAll(t provider.T) {
	tCfg, err := testutils.NewConfig()
	t.Require().NoError(err)

	mCfg, err := migrations.NewConfig()
	t.Require().NoError(err)

	s.ctx = context.Background()

	s.db, err = testutils.NewTestDB(s.ctx, tCfg, mCfg)
	t.Require().NoError(err)

	s.orderCollection = s.db.DB.Collection(s.db.Cfg.OrderCollection)
}

func (s *OrderRepositoryTestSuite) AfterAll(t provider.T) {
	if s.db != nil {
		err := s.db.Close(s.ctx)
		t.Require().NoError(err)
	}
}

func (s *OrderRepositoryTestSuite) AfterEach(t provider.T) {
	ctx, cancel := context.WithTimeout(s.ctx, 5*time.Second)
	defer cancel()

	err := s.db.Clear(ctx)
	t.Require().NoError(err)
}

func (s *OrderRepositoryTestSuite) getRepo() orderDomain.Repository {
	return orderRepository.New(s.orderCollection)
}

func (s *OrderRepositoryTestSuite) TestCreate(t provider.T) {
	tests := []struct {
		name          string
		setup         func(repo orderDomain.Repository) *orderDomain.Order
		expectedError error
	}{
		{
			name: "Success",
			setup: func(_ orderDomain.Repository) *orderDomain.Order {
				return mothers.DefaultOrder()
			},
			expectedError: nil,
		},
		{
			name: "Failure: Order already exists",
			setup: func(repo orderDomain.Repository) *orderDomain.Order {
				order := mothers.DefaultOrder()
				err := repo.Create(s.ctx, order)
				t.Require().NoError(err)
				return order
			},
			expectedError: orderRepository.ErrOrderAlreadyExists,
		},
	}

	repo := s.getRepo()
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t provider.T) {
			order := tc.setup(repo)
			err := repo.Create(s.ctx, order)

			if tc.expectedError != nil {
				t.Require().Error(err)
				t.Require().Equal(tc.expectedError, err)
			} else {
				t.Require().NoError(err)

				createdOrder, err := repo.GetByID(s.ctx, order.ID)
				t.Require().NoError(err)
				t.Require().NotNil(createdOrder)
				t.Require().Equal(order.ID, createdOrder.ID)
			}
		})
	}
}

func (s *OrderRepositoryTestSuite) TestUpdate(t provider.T) {
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
				order := mothers.DefaultOrder()
				err := repo.Create(s.ctx, order)
				t.Require().NoError(err)
				return order
			},
			update: func(order *orderDomain.Order) {
				err := order.NoteDelivering(uuid.New())
				t.Require().NoError(err)
			},
			verify: func(updated, expected *orderDomain.Order) {
				t.Require().Equal(orderDomain.Delivering, updated.Status)
				t.Require().Equal(expected.ID, updated.ID)
			},
			expectedError: nil,
		},
		{
			name: "Success: Update delivery",
			setup: func(repo orderDomain.Repository) *orderDomain.Order {
				order := mothers.DefaultOrder()
				err := repo.Create(s.ctx, order)
				t.Require().NoError(err)
				return order
			},
			update: func(order *orderDomain.Order) {
				err := order.NoteDelivering(uuid.New())
				t.Require().NoError(err)
			},
			verify: func(updated, expected *orderDomain.Order) {
				t.Require().Equal(expected.ID, updated.ID)
				t.Require().Equal(expected.Delivery.CourierID, updated.Delivery.CourierID)
			},
			expectedError: nil,
		},
		{
			name: "Failure: Order not found",
			setup: func(repo orderDomain.Repository) *orderDomain.Order {
				return mothers.DefaultOrder()
			},
			update:        func(order *orderDomain.Order) {},
			verify:        func(updated, expected *orderDomain.Order) {},
			expectedError: orderRepository.ErrOrderNotFound,
		},
	}

	repo := s.getRepo()
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t provider.T) {
			order := tc.setup(repo)
			tc.update(order)

			err := repo.Update(s.ctx, order)

			if tc.expectedError != nil {
				t.Require().Error(err)
				t.Require().Equal(tc.expectedError, err)
			} else {
				t.Require().NoError(err)
				updatedOrder, err := repo.GetByID(s.ctx, order.ID)
				t.Require().NoError(err)
				t.Require().NotNil(updatedOrder)
				tc.verify(updatedOrder, order)
			}
		})
	}
}

func (s *OrderRepositoryTestSuite) TestGetByID(t provider.T) {
	tests := []struct {
		name          string
		setup         func(repo orderDomain.Repository) uuid.UUID
		expectedError error
	}{
		{
			name: "Success",
			setup: func(repo orderDomain.Repository) uuid.UUID {
				order := mothers.DefaultOrder()
				err := repo.Create(s.ctx, order)
				t.Require().NoError(err)
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
		tc := tc
		t.Run(tc.name, func(t provider.T) {
			orderID := tc.setup(repo)

			order, err := repo.GetByID(s.ctx, orderID)

			if tc.expectedError != nil {
				t.Require().Error(err)
				t.Require().Equal(tc.expectedError, err)
			} else {
				t.Require().NoError(err)
				t.Require().NotNil(order)
				t.Require().Equal(orderID, order.ID)
			}
		})
	}
}

func (s *OrderRepositoryTestSuite) TestGetAllByCustomer(t provider.T) {
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
				order := mothers.DefaultOrder()
				order.CustomerID = customerID
				err := repo.Create(s.ctx, order)
				t.Require().NoError(err)
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
					order := mothers.DefaultOrder()
					order.CustomerID = customerID
					err := repo.Create(s.ctx, order)
					t.Require().NoError(err)
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
		tc := tc
		t.Run(tc.name, func(t provider.T) {
			orderIDs := tc.setup(repo, tc.customerID)

			orders, err := repo.GetAllByCustomer(s.ctx, tc.customerID)

			if tc.expectedError != nil {
				t.Require().Error(err)
				t.Require().Equal(tc.expectedError, err)
			} else {
				t.Require().NoError(err)
				t.Require().Equal(len(orderIDs), len(orders))
				for _, order := range orders {
					t.Require().Contains(orderIDs, order.ID)
				}
			}
		})
	}
}

func (s *OrderRepositoryTestSuite) TestGetCurrentByCourier(t provider.T) {
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
				order := mothers.OrderDelivering()
				order.Delivery.CourierID = &courierID

				err := repo.Create(s.ctx, order)
				t.Require().NoError(err)

				return []uuid.UUID{order.ID}
			},
		},
		{
			name:      "Success: Multiple orders",
			courierID: uuid.New(),
			setup: func(repo orderDomain.Repository, courierID uuid.UUID) []uuid.UUID {
				var orderIDs []uuid.UUID
				for i := 0; i < 3; i++ {
					order := mothers.OrderDelivering()
					order.Delivery.CourierID = &courierID

					err := repo.Create(s.ctx, order)
					t.Require().NoError(err)

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
		tc := tc
		t.Run(tc.name, func(t provider.T) {
			orderIDs := tc.setup(repo, tc.courierID)

			orders, err := repo.GetCurrentByCourier(s.ctx, tc.courierID)

			if tc.expectedError != nil {
				t.Require().Error(err)
				t.Require().Equal(tc.expectedError, err)
			} else {
				t.Require().NoError(err)
				t.Require().Equal(len(orderIDs), len(orders))
				for _, order := range orders {
					t.Require().Contains(orderIDs, order.ID)
				}
			}
		})
	}
}

func TestOrderRepository(t *testing.T) {
	suite.RunSuite(t, new(OrderRepositoryTestSuite))
}
