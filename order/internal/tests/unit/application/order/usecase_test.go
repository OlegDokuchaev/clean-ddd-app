package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"order/internal/application/order/usecase"
	orderDomain "order/internal/domain/order"
	orderMock "order/internal/mocks/order"
	createOrderMock "order/internal/mocks/order/saga/create_order"
)

type OrderUseCaseTestSuite struct {
	suite.Suite
	ctx context.Context
}

func (s *OrderUseCaseTestSuite) SetupTest() {
	s.ctx = context.Background()
}

func (s *OrderUseCaseTestSuite) TestCreate() {
	tests := []struct {
		name        string
		dto         usecase.CreateDto
		setup       func(repo *orderMock.RepositoryMock, manager *createOrderMock.ManagerMock)
		expectedErr error
	}{
		{
			name: "Success",
			dto: usecase.CreateDto{
				CustomerID: uuid.New(),
				Address:    "Test Address",
				Items:      []orderDomain.Item{},
			},
			setup: func(repo *orderMock.RepositoryMock, manager *createOrderMock.ManagerMock) {
				repo.On("Create", s.ctx, mock.Anything).Return(nil).Once()
				manager.On("Create", s.ctx, mock.Anything).Return(nil).Once()
			},
			expectedErr: nil,
		},
		{
			name: "Failure: repository error",
			dto: usecase.CreateDto{
				CustomerID: uuid.New(),
				Address:    "Test Address",
				Items:      []orderDomain.Item{},
			},
			setup: func(repo *orderMock.RepositoryMock, manager *createOrderMock.ManagerMock) {
				repo.On("Create", s.ctx, mock.Anything).Return(errors.New("repo error")).Once()
			},
			expectedErr: errors.New("repo error"),
		},
	}

	for _, tc := range tests {
		tc := tc
		s.Run(tc.name, func() {
			s.T().Parallel()
			repo := new(orderMock.RepositoryMock)
			manager := new(createOrderMock.ManagerMock)
			uc := usecase.New(repo, manager)
			tc.setup(repo, manager)

			orderID, err := uc.Create(s.ctx, tc.dto)

			if tc.expectedErr == nil {
				require.NoError(s.T(), err)
				require.NotEqual(s.T(), uuid.Nil, orderID)
			} else {
				require.Error(s.T(), err)
				require.EqualError(s.T(), err, tc.expectedErr.Error())
				require.Equal(s.T(), uuid.Nil, orderID)
			}

			repo.AssertExpectations(s.T())
			manager.AssertExpectations(s.T())
		})
	}
}

func (s *OrderUseCaseTestSuite) TestCancelByCustomer() {
	tests := []struct {
		name        string
		setup       func(repo *orderMock.RepositoryMock) *orderDomain.Order
		expectedErr error
		finalStatus orderDomain.Status
	}{
		{
			name: "Success: Order in Delivering",
			setup: func(repo *orderMock.RepositoryMock) *orderDomain.Order {
				o := orderDomain.Create(uuid.New(), "Some Address", nil)
				require.NoError(s.T(), o.NoteDelivering(uuid.New()))
				repo.On("GetByID", s.ctx, o.ID).Return(o, nil).Once()
				repo.On("Update", s.ctx, o).Return(nil).Once()
				return o
			},
			expectedErr: nil,
			finalStatus: orderDomain.CustomerCanceled,
		},
		{
			name: "Failure: repo.GetByID error",
			setup: func(repo *orderMock.RepositoryMock) *orderDomain.Order {
				o := orderDomain.Create(uuid.New(), "Some Address", nil)
				repo.On("GetByID", s.ctx, o.ID).
					Return((*orderDomain.Order)(nil), errors.New("not found")).Once()
				return o
			},
			expectedErr: errors.New("not found"),
			finalStatus: orderDomain.Created,
		},
		{
			name: "Failure: domain method error (order in Created)",
			setup: func(repo *orderMock.RepositoryMock) *orderDomain.Order {
				o := orderDomain.Create(uuid.New(), "Some Address", nil)
				repo.On("GetByID", s.ctx, o.ID).Return(o, nil).Once()
				return o
			},
			expectedErr: orderDomain.ErrUnsupportedStatusTransition,
			finalStatus: orderDomain.Created,
		},
		{
			name: "Failure: repo.Update error",
			setup: func(repo *orderMock.RepositoryMock) *orderDomain.Order {
				o := orderDomain.Create(uuid.New(), "Some Address", nil)
				require.NoError(s.T(), o.NoteDelivering(uuid.New()))
				repo.On("GetByID", s.ctx, o.ID).Return(o, nil).Once()
				repo.On("Update", s.ctx, o).Return(errors.New("update error")).Once()
				return o
			},
			expectedErr: errors.New("update error"),
			finalStatus: orderDomain.CustomerCanceled,
		},
	}
	for _, tc := range tests {
		tc := tc
		s.Run(tc.name, func() {
			s.T().Parallel()
			repo := new(orderMock.RepositoryMock)
			manager := new(createOrderMock.ManagerMock)
			uc := usecase.New(repo, manager)
			o := tc.setup(repo)

			err := uc.CancelByCustomer(s.ctx, o.ID)

			if tc.expectedErr == nil {
				require.NoError(s.T(), err)
			} else {
				require.Error(s.T(), err)
				require.EqualError(s.T(), err, tc.expectedErr.Error())
			}
			require.Equal(s.T(), tc.finalStatus, o.Status)

			repo.AssertExpectations(s.T())
			manager.AssertExpectations(s.T())
		})
	}
}

func (s *OrderUseCaseTestSuite) TestCancelOutOfStock() {
	tests := []struct {
		name        string
		setup       func(repo *orderMock.RepositoryMock) *orderDomain.Order
		expectedErr error
		finalStatus orderDomain.Status
	}{
		{
			name: "Success: Order in Created (default)",
			setup: func(repo *orderMock.RepositoryMock) *orderDomain.Order {
				o := orderDomain.Create(uuid.New(), "Some Address", nil)
				repo.On("GetByID", s.ctx, o.ID).Return(o, nil).Once()
				repo.On("Update", s.ctx, o).Return(nil).Once()
				return o
			},
			expectedErr: nil,
			finalStatus: orderDomain.CanceledOutOfStock,
		},
		{
			name: "Failure: GetByID error",
			setup: func(repo *orderMock.RepositoryMock) *orderDomain.Order {
				o := orderDomain.Create(uuid.New(), "Some Address", nil)
				repo.On("GetByID", s.ctx, o.ID).
					Return((*orderDomain.Order)(nil), errors.New("not found")).Once()
				return o
			},
			expectedErr: errors.New("not found"),
			finalStatus: orderDomain.Created,
		},
		{
			name: "Failure: domain method error (order in Delivering)",
			setup: func(repo *orderMock.RepositoryMock) *orderDomain.Order {
				o := orderDomain.Create(uuid.New(), "Some Address", nil)
				require.NoError(s.T(), o.NoteDelivering(uuid.New()))
				repo.On("GetByID", s.ctx, o.ID).Return(o, nil).Once()
				return o
			},
			expectedErr: orderDomain.ErrUnsupportedStatusTransition,
			finalStatus: orderDomain.Delivering,
		},
		{
			name: "Failure: Update error",
			setup: func(repo *orderMock.RepositoryMock) *orderDomain.Order {
				o := orderDomain.Create(uuid.New(), "Some Address", nil)
				repo.On("GetByID", s.ctx, o.ID).Return(o, nil).Once()
				repo.On("Update", s.ctx, o).Return(errors.New("update error")).Once()
				return o
			},
			expectedErr: errors.New("update error"),
			finalStatus: orderDomain.CanceledOutOfStock,
		},
	}
	for _, tc := range tests {
		tc := tc
		s.Run(tc.name, func() {
			s.T().Parallel()
			repo := new(orderMock.RepositoryMock)
			manager := new(createOrderMock.ManagerMock)
			uc := usecase.New(repo, manager)
			o := tc.setup(repo)

			err := uc.CancelOutOfStock(s.ctx, o.ID)

			if tc.expectedErr == nil {
				require.NoError(s.T(), err)
			} else {
				require.Error(s.T(), err)
				require.EqualError(s.T(), err, tc.expectedErr.Error())
			}
			require.Equal(s.T(), tc.finalStatus, o.Status)

			repo.AssertExpectations(s.T())
			manager.AssertExpectations(s.T())
		})
	}
}

func (s *OrderUseCaseTestSuite) TestCancelCourierNotFound() {
	tests := []struct {
		name        string
		setup       func(repo *orderMock.RepositoryMock) *orderDomain.Order
		expectedErr error
		finalStatus orderDomain.Status
	}{
		{
			name: "Success: Order in Created",
			setup: func(repo *orderMock.RepositoryMock) *orderDomain.Order {
				o := orderDomain.Create(uuid.New(), "Some Address", nil)
				repo.On("GetByID", s.ctx, o.ID).Return(o, nil).Once()
				repo.On("Update", s.ctx, o).Return(nil).Once()
				return o
			},
			expectedErr: nil,
			finalStatus: orderDomain.CanceledCourierNotFound,
		},
		{
			name: "Failure: GetByID error",
			setup: func(repo *orderMock.RepositoryMock) *orderDomain.Order {
				o := orderDomain.Create(uuid.New(), "Some Address", nil)
				repo.On("GetByID", s.ctx, o.ID).
					Return((*orderDomain.Order)(nil), errors.New("not found")).Once()
				return o
			},
			expectedErr: errors.New("not found"),
			finalStatus: orderDomain.Created,
		},
		{
			name: "Failure: domain method error (order in Delivering)",
			setup: func(repo *orderMock.RepositoryMock) *orderDomain.Order {
				o := orderDomain.Create(uuid.New(), "Some Address", nil)
				require.NoError(s.T(), o.NoteDelivering(uuid.New()))
				repo.On("GetByID", s.ctx, o.ID).Return(o, nil).Once()
				return o
			},
			expectedErr: orderDomain.ErrUnsupportedStatusTransition,
			finalStatus: orderDomain.Delivering,
		},
		{
			name: "Failure: Update error",
			setup: func(repo *orderMock.RepositoryMock) *orderDomain.Order {
				o := orderDomain.Create(uuid.New(), "Some Address", nil)
				repo.On("GetByID", s.ctx, o.ID).Return(o, nil).Once()
				repo.On("Update", s.ctx, o).Return(errors.New("update error")).Once()
				return o
			},
			expectedErr: errors.New("update error"),
			finalStatus: orderDomain.CanceledCourierNotFound,
		},
	}
	for _, tc := range tests {
		tc := tc
		s.Run(tc.name, func() {
			s.T().Parallel()
			repo := new(orderMock.RepositoryMock)
			manager := new(createOrderMock.ManagerMock)
			uc := usecase.New(repo, manager)
			o := tc.setup(repo)

			err := uc.CancelCourierNotFound(s.ctx, o.ID)

			if tc.expectedErr == nil {
				require.NoError(s.T(), err)
			} else {
				require.Error(s.T(), err)
				require.EqualError(s.T(), err, tc.expectedErr.Error())
			}
			require.Equal(s.T(), tc.finalStatus, o.Status)

			repo.AssertExpectations(s.T())
			manager.AssertExpectations(s.T())
		})
	}
}

func (s *OrderUseCaseTestSuite) TestBeginDelivery() {
	tests := []struct {
		name        string
		setup       func(repo *orderMock.RepositoryMock) (usecase.BeginDeliveryDto, *orderDomain.Order)
		expectedErr error
		finalStatus orderDomain.Status
	}{
		{
			name: "Success: Order in Created",
			setup: func(repo *orderMock.RepositoryMock) (usecase.BeginDeliveryDto, *orderDomain.Order) {
				o := orderDomain.Create(uuid.New(), "Some Address", nil)
				dto := usecase.BeginDeliveryDto{
					OrderID:   o.ID,
					CourierID: uuid.New(),
				}
				repo.On("GetByID", s.ctx, o.ID).Return(o, nil).Once()
				repo.On("Update", s.ctx, o).Return(nil).Once()
				return dto, o
			},
			expectedErr: nil,
			finalStatus: orderDomain.Delivering,
		},
		{
			name: "Failure: domain method error (already in Delivering)",
			setup: func(repo *orderMock.RepositoryMock) (usecase.BeginDeliveryDto, *orderDomain.Order) {
				o := orderDomain.Create(uuid.New(), "Some Address", nil)
				require.NoError(s.T(), o.NoteDelivering(uuid.New()))
				dto := usecase.BeginDeliveryDto{
					OrderID:   o.ID,
					CourierID: uuid.New(),
				}
				repo.On("GetByID", s.ctx, o.ID).Return(o, nil).Once()
				return dto, o
			},
			expectedErr: orderDomain.ErrUnsupportedStatusTransition,
			finalStatus: orderDomain.Delivering,
		},
		{
			name: "Failure: repo.GetByID error",
			setup: func(repo *orderMock.RepositoryMock) (usecase.BeginDeliveryDto, *orderDomain.Order) {
				o := orderDomain.Create(uuid.New(), "Some Address", nil)
				dto := usecase.BeginDeliveryDto{
					OrderID:   o.ID,
					CourierID: uuid.New(),
				}
				repo.On("GetByID", s.ctx, o.ID).
					Return((*orderDomain.Order)(nil), errors.New("not found")).Once()
				return dto, o
			},
			expectedErr: errors.New("not found"),
			finalStatus: orderDomain.Created,
		},
		{
			name: "Failure: Update error",
			setup: func(repo *orderMock.RepositoryMock) (usecase.BeginDeliveryDto, *orderDomain.Order) {
				o := orderDomain.Create(uuid.New(), "Some Address", nil)
				dto := usecase.BeginDeliveryDto{
					OrderID:   o.ID,
					CourierID: uuid.New(),
				}
				repo.On("GetByID", s.ctx, o.ID).Return(o, nil).Once()
				repo.On("Update", s.ctx, o).Return(errors.New("update error")).Once()
				return dto, o
			},
			expectedErr: errors.New("update error"),
			finalStatus: orderDomain.Delivering,
		},
	}
	for _, tc := range tests {
		tc := tc
		s.Run(tc.name, func() {
			s.T().Parallel()
			repo := new(orderMock.RepositoryMock)
			manager := new(createOrderMock.ManagerMock)
			uc := usecase.New(repo, manager)
			dto, o := tc.setup(repo)

			err := uc.BeginDelivery(s.ctx, dto)

			if tc.expectedErr == nil {
				require.NoError(s.T(), err)
			} else {
				require.Error(s.T(), err)
				require.EqualError(s.T(), err, tc.expectedErr.Error())
			}
			require.Equal(s.T(), tc.finalStatus, o.Status)

			repo.AssertExpectations(s.T())
			manager.AssertExpectations(s.T())
		})
	}
}

func (s *OrderUseCaseTestSuite) TestCompleteDelivery() {
	tests := []struct {
		name        string
		setup       func(repo *orderMock.RepositoryMock) *orderDomain.Order
		expectedErr error
		finalStatus orderDomain.Status
	}{
		{
			name: "Success: Order in Delivering",
			setup: func(repo *orderMock.RepositoryMock) *orderDomain.Order {
				o := orderDomain.Create(uuid.New(), "Some Address", nil)
				require.NoError(s.T(), o.NoteDelivering(uuid.New()))
				repo.On("GetByID", s.ctx, o.ID).Return(o, nil).Once()
				repo.On("Update", s.ctx, o).Return(nil).Once()
				return o
			},
			expectedErr: nil,
			finalStatus: orderDomain.Delivered,
		},
		{
			name: "Failure: GetByID error",
			setup: func(repo *orderMock.RepositoryMock) *orderDomain.Order {
				o := orderDomain.Create(uuid.New(), "Some Address", nil)
				repo.On("GetByID", s.ctx, o.ID).
					Return((*orderDomain.Order)(nil), errors.New("not found")).Once()
				return o
			},
			expectedErr: errors.New("not found"),
			finalStatus: orderDomain.Created,
		},
		{
			name: "Failure: domain method error (order in Created)",
			setup: func(repo *orderMock.RepositoryMock) *orderDomain.Order {
				o := orderDomain.Create(uuid.New(), "Some Address", nil)
				repo.On("GetByID", s.ctx, o.ID).Return(o, nil).Once()
				return o
			},
			expectedErr: orderDomain.ErrUnsupportedStatusTransition,
			finalStatus: orderDomain.Created,
		},
		{
			name: "Failure: Update error",
			setup: func(repo *orderMock.RepositoryMock) *orderDomain.Order {
				o := orderDomain.Create(uuid.New(), "Some Address", nil)
				require.NoError(s.T(), o.NoteDelivering(uuid.New()))
				repo.On("GetByID", s.ctx, o.ID).Return(o, nil).Once()
				repo.On("Update", s.ctx, o).Return(errors.New("update error")).Once()
				return o
			},
			expectedErr: errors.New("update error"),
			finalStatus: orderDomain.Delivered,
		},
	}
	for _, tc := range tests {
		tc := tc
		s.Run(tc.name, func() {
			s.T().Parallel()
			repo := new(orderMock.RepositoryMock)
			manager := new(createOrderMock.ManagerMock)
			uc := usecase.New(repo, manager)
			o := tc.setup(repo)

			err := uc.CompleteDelivery(s.ctx, o.ID)

			if tc.expectedErr == nil {
				require.NoError(s.T(), err)
				require.NotNil(s.T(), o.Delivery.Arrived)
				require.WithinDuration(s.T(), time.Now(), *o.Delivery.Arrived, time.Second)
			} else {
				require.Error(s.T(), err)
				require.EqualError(s.T(), err, tc.expectedErr.Error())
			}
			require.Equal(s.T(), tc.finalStatus, o.Status)

			repo.AssertExpectations(s.T())
			manager.AssertExpectations(s.T())
		})
	}
}

func (s *OrderUseCaseTestSuite) TestGetAllByCustomer() {
	tests := []struct {
		name        string
		setup       func(repo *orderMock.RepositoryMock) (uuid.UUID, []*orderDomain.Order)
		expectedErr error
	}{
		{
			name: "Success: Get orders by customer",
			setup: func(repo *orderMock.RepositoryMock) (uuid.UUID, []*orderDomain.Order) {
				customerID := uuid.New()
				expectedOrders := []*orderDomain.Order{
					orderDomain.Create(customerID, "Address 1", []orderDomain.Item{}),
					orderDomain.Create(customerID, "Address 2", []orderDomain.Item{}),
				}
				repo.On("GetAllByCustomer", s.ctx, customerID).Return(expectedOrders, nil).Once()
				return customerID, expectedOrders
			},
			expectedErr: nil,
		},
		{
			name: "Failure: GetAllByCustomer error",
			setup: func(repo *orderMock.RepositoryMock) (uuid.UUID, []*orderDomain.Order) {
				customerID := uuid.New()
				var expectedOrders []*orderDomain.Order
				repo.On("GetAllByCustomer", s.ctx, customerID).
					Return(expectedOrders, errors.New("get error")).Once()
				return customerID, expectedOrders
			},
			expectedErr: errors.New("get error"),
		},
	}
	for _, tc := range tests {
		tc := tc
		s.Run(tc.name, func() {
			s.T().Parallel()
			repo := new(orderMock.RepositoryMock)
			manager := new(createOrderMock.ManagerMock)
			uc := usecase.New(repo, manager)
			customerID, expectedOrders := tc.setup(repo)

			orders, err := uc.GetAllByCustomer(s.ctx, customerID)

			if tc.expectedErr == nil {
				require.NoError(s.T(), err)
				require.Equal(s.T(), expectedOrders, orders)
			} else {
				require.Error(s.T(), err)
				require.EqualError(s.T(), err, tc.expectedErr.Error())
			}

			repo.AssertExpectations(s.T())
			manager.AssertExpectations(s.T())
		})
	}
}

func (s *OrderUseCaseTestSuite) TestGetCurrentByCourier() {
	tests := []struct {
		name        string
		setup       func(repo *orderMock.RepositoryMock) (uuid.UUID, []*orderDomain.Order)
		expectedErr error
	}{
		{
			name: "Success: Get current orders by courier",
			setup: func(repo *orderMock.RepositoryMock) (uuid.UUID, []*orderDomain.Order) {
				courierID := uuid.New()
				expectedOrders := []*orderDomain.Order{
					orderDomain.Create(uuid.New(), "Address 1", []orderDomain.Item{}),
					orderDomain.Create(uuid.New(), "Address 2", []orderDomain.Item{}),
				}
				repo.On("GetCurrentByCourier", s.ctx, courierID).Return(expectedOrders, nil).Once()
				return courierID, expectedOrders
			},
			expectedErr: nil,
		},
		{
			name: "Failure: GetCurrentByCourier error",
			setup: func(repo *orderMock.RepositoryMock) (uuid.UUID, []*orderDomain.Order) {
				courierID := uuid.New()
				var expectedOrders []*orderDomain.Order
				repo.On("GetCurrentByCourier", s.ctx, courierID).
					Return(expectedOrders, errors.New("get error")).Once()
				return courierID, expectedOrders
			},
			expectedErr: errors.New("get error"),
		},
	}
	for _, tc := range tests {
		tc := tc
		s.Run(tc.name, func() {
			s.T().Parallel()
			repo := new(orderMock.RepositoryMock)
			manager := new(createOrderMock.ManagerMock)
			uc := usecase.New(repo, manager)
			courierID, expectedOrders := tc.setup(repo)

			orders, err := uc.GetCurrentByCourier(s.ctx, courierID)

			if tc.expectedErr == nil {
				require.NoError(s.T(), err)
				require.Equal(s.T(), expectedOrders, orders)
			} else {
				require.Error(s.T(), err)
				require.EqualError(s.T(), err, tc.expectedErr.Error())
			}

			repo.AssertExpectations(s.T())
			manager.AssertExpectations(s.T())
		})
	}
}

func TestOrderUseCaseTestSuite(t *testing.T) {
	suite.Run(t, new(OrderUseCaseTestSuite))
}
