package product

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"testing"
	productApplication "warehouse/internal/application/product"
	productDomain "warehouse/internal/domain/product"
	"warehouse/internal/mocks"
)

type ProductUseCaseTestSuite struct {
	suite.Suite
	ctx context.Context
}

func (s *ProductUseCaseTestSuite) SetupTest() {
	s.ctx = context.Background()
}

func (s *ProductUseCaseTestSuite) TestCreate() {
	tests := []struct {
		name        string
		dto         productApplication.CreateDto
		setup       func(uowMock *mocks.UoWMock)
		expectedErr error
	}{
		{
			name: "Success",
			dto:  productApplication.CreateDto{Name: "ValidProduct", Price: decimal.NewFromInt(1)},
			setup: func(uow *mocks.UoWMock) {
				uow.On("Transaction", s.ctx, mock.Anything).Once()
				uow.ProductMock.On("Create", s.ctx, mock.Anything).
					Return(nil).Once()
				uow.OutboxMock.On("Create", s.ctx, mock.Anything).
					Return(nil).Once()
			},
			expectedErr: nil,
		},
		{
			name:        "Failure: Create order error (invalid price)",
			dto:         productApplication.CreateDto{Name: "ValidProduct", Price: decimal.NewFromInt(-1)},
			setup:       func(uow *mocks.UoWMock) {},
			expectedErr: productDomain.ErrInvalidProductPrice,
		},
		{
			name: "Failure: Product repository create error",
			dto:  productApplication.CreateDto{Name: "ValidProduct", Price: decimal.NewFromInt(1)},
			setup: func(uow *mocks.UoWMock) {
				uow.On("Transaction", s.ctx, mock.Anything).Once()
				uow.ProductMock.On("Create", s.ctx, mock.Anything).
					Return(errors.New("failed to build message")).Once()
			},
			expectedErr: errors.New("failed to build message"),
		},
		{
			name: "Failure: Outbox repository create error",
			dto:  productApplication.CreateDto{Name: "ValidProduct", Price: decimal.NewFromInt(1)},
			setup: func(uow *mocks.UoWMock) {
				uow.On("Transaction", s.ctx, mock.Anything).Once()
				uow.ProductMock.On("Create", s.ctx, mock.Anything).
					Return(nil).Once()
				uow.OutboxMock.On("Create", s.ctx, mock.Anything).
					Return(errors.New("failed to build message")).Once()
			},
			expectedErr: errors.New("failed to build message"),
		},
		{
			name: "Failure: UoW transaction error",
			dto:  productApplication.CreateDto{Name: "ValidProduct", Price: decimal.NewFromInt(1)},
			setup: func(uow *mocks.UoWMock) {
				uow.On("Transaction", s.ctx, mock.Anything).
					Return(errors.New("transaction error")).Once()
			},
			expectedErr: errors.New("transaction error"),
		},
	}

	for _, tc := range tests {
		tc := tc
		s.Run(tc.name, func() {
			s.T().Parallel()
			uow := mocks.NewUowMock()
			tc.setup(uow)
			usecase := productApplication.NewUseCase(uow)

			productID, err := usecase.Create(s.ctx, tc.dto)

			if tc.expectedErr == nil {
				require.NoError(s.T(), err)
				require.NotEqual(s.T(), uuid.Nil, productID)
			} else {
				require.Error(s.T(), err)
				require.EqualError(s.T(), err, tc.expectedErr.Error())
				require.Equal(s.T(), uuid.Nil, productID)
			}

			uow.AssertExpectations(s.T())
		})
	}
}

func TestProductUseCaseTestSuite(t *testing.T) {
	suite.Run(t, new(ProductUseCaseTestSuite))
}
