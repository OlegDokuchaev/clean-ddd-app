package item

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/shopspring/decimal"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	itemApplication "warehouse/internal/application/item"
	itemDomain "warehouse/internal/domain/item"
	productDomain "warehouse/internal/domain/product"
	itemRepository "warehouse/internal/infrastructure/repository/item"
	"warehouse/internal/mocks"
)

type ItemUseCaseTestSuite struct {
	suite.Suite
	ctx context.Context
}

func (s *ItemUseCaseTestSuite) SetupTest() {
	s.ctx = context.Background()
}

func (s *ItemUseCaseTestSuite) createTestItem(count int) *itemDomain.Item {
	product, _, err := productDomain.Create("test", decimal.NewFromInt(1))
	require.NoError(s.T(), err)
	item, err := itemDomain.Create(product, count)
	require.NoError(s.T(), err)
	return item
}

func (s *ItemUseCaseTestSuite) createTestItems(counts ...int) []*itemDomain.Item {
	items := make([]*itemDomain.Item, len(counts))
	for i, count := range counts {
		product, _, err := productDomain.Create(fmt.Sprintf("test-product-%d", i), decimal.NewFromInt(int64(10+i)))
		require.NoError(s.T(), err)
		item, err := itemDomain.Create(product, count)
		require.NoError(s.T(), err)
		items[i] = item
	}
	return items
}

func (s *ItemUseCaseTestSuite) TestCreate() {
	tests := []struct {
		name        string
		dto         itemApplication.CreateDto
		setup       func(uowMock *mocks.UoWMock)
		expectedErr error
	}{
		{
			name: "Success",
			dto:  itemApplication.CreateDto{ProductID: uuid.New(), Count: 10},
			setup: func(uow *mocks.UoWMock) {
				product, _, err := productDomain.Create("test", decimal.NewFromInt(1))
				require.NoError(s.T(), err)
				uow.ProductMock.On("GetByID", s.ctx, mock.Anything).Return(product, nil).Once()

				uow.ItemMock.On("Create", s.ctx, mock.Anything).
					Return(nil).Once()
			},
			expectedErr: nil,
		},
		{
			name: "Failure: Product not found",
			dto:  itemApplication.CreateDto{ProductID: uuid.New(), Count: 5},
			setup: func(uow *mocks.UoWMock) {
				uow.ProductMock.On("GetByID", s.ctx, mock.Anything).
					Return((*productDomain.Product)(nil), errors.New("product not found")).Once()
			},
			expectedErr: errors.New("product not found"),
		},
		{
			name: "Failure: Create item error (negative count)",
			dto:  itemApplication.CreateDto{ProductID: uuid.New(), Count: -1},
			setup: func(uow *mocks.UoWMock) {
				product, _, err := productDomain.Create("test", decimal.NewFromInt(1))
				require.NoError(s.T(), err)
				uow.ProductMock.On("GetByID", s.ctx, mock.Anything).Return(product, nil).Once()
			},
			expectedErr: itemDomain.ErrInvalidItemCount,
		},
		{
			name: "Failure: Item repository create error",
			dto:  itemApplication.CreateDto{ProductID: uuid.New(), Count: 10},
			setup: func(uow *mocks.UoWMock) {
				product, _, err := productDomain.Create("test", decimal.NewFromInt(1))
				require.NoError(s.T(), err)
				uow.ProductMock.On("GetByID", s.ctx, mock.Anything).Return(product, nil).Once()

				uow.ItemMock.On("Create", s.ctx, mock.Anything).
					Return(errors.New("item create error")).Once()
			},
			expectedErr: errors.New("item create error"),
		},
	}

	for _, tc := range tests {
		tc := tc
		s.Run(tc.name, func() {
			s.T().Parallel()
			uow := mocks.NewUowMock()
			tc.setup(uow)
			useCase := itemApplication.NewUseCase(uow)

			itemID, err := useCase.Create(s.ctx, tc.dto)

			if tc.expectedErr == nil {
				require.NoError(s.T(), err)
				require.NotEqual(s.T(), uuid.Nil, itemID)
			} else {
				require.Error(s.T(), err)
				require.EqualError(s.T(), err, tc.expectedErr.Error())
				require.Equal(s.T(), uuid.Nil, itemID)
			}

			uow.AssertExpectations(s.T())
		})
	}
}

func (s *ItemUseCaseTestSuite) TestReserve() {
	tests := []struct {
		name        string
		setup       func(uow *mocks.UoWMock) itemApplication.ReserveDto
		expectedErr error
	}{
		{
			name: "Success",
			setup: func(uow *mocks.UoWMock) itemApplication.ReserveDto {
				items := s.createTestItems(5, 10)

				uow.ItemMock.On("GetAllByProductIDs", s.ctx, items[0].Product.ID, items[1].Product.ID).
					Return(items, nil).Once()
				uow.ItemMock.On("Update", s.ctx, mock.Anything).Return(nil).Twice()

				uow.On("Transaction", s.ctx, mock.Anything).Once()

				return itemApplication.ReserveDto{
					Items: []itemApplication.ItemDto{
						{ProductID: items[0].Product.ID, Count: 4},
						{ProductID: items[1].Product.ID, Count: 9},
					},
				}
			},
			expectedErr: nil,
		},
		{
			name: "Failure: Items not found",
			setup: func(uow *mocks.UoWMock) itemApplication.ReserveDto {
				uow.ItemMock.On("GetAllByProductIDs", s.ctx, mock.Anything, mock.Anything).
					Return([]*itemDomain.Item{}, itemRepository.ErrItemsNotFound).Once()
				uow.On("Transaction", s.ctx, mock.Anything).Once()

				return itemApplication.ReserveDto{
					Items: []itemApplication.ItemDto{
						{ProductID: uuid.New(), Count: 2},
						{ProductID: uuid.New(), Count: 3},
					},
				}
			},
			expectedErr: itemRepository.ErrItemsNotFound,
		},
		{
			name: "Failure: Reserve item error",
			setup: func(uow *mocks.UoWMock) itemApplication.ReserveDto {
				items := s.createTestItems(1, 10)

				uow.ItemMock.On("GetAllByProductIDs", s.ctx, items[0].Product.ID, items[1].Product.ID).
					Return(items, nil).Once()

				uow.On("Transaction", s.ctx, mock.Anything).Once()

				return itemApplication.ReserveDto{
					Items: []itemApplication.ItemDto{
						{ProductID: items[0].Product.ID, Count: 4},
						{ProductID: items[1].Product.ID, Count: 9},
					},
				}
			},
			expectedErr: itemDomain.ErrInvalidItemCount,
		},
		{
			name: "Failure: Update item error",
			setup: func(uow *mocks.UoWMock) itemApplication.ReserveDto {
				items := s.createTestItems(5, 10)

				uow.ItemMock.On("GetAllByProductIDs", s.ctx, items[0].Product.ID, items[1].Product.ID).
					Return(items, nil).Once()
				uow.ItemMock.On("Update", s.ctx, mock.Anything).
					Return(errors.New("update error")).Once()

				uow.On("Transaction", s.ctx, mock.Anything).Once()

				return itemApplication.ReserveDto{
					Items: []itemApplication.ItemDto{
						{ProductID: items[0].Product.ID, Count: 4},
						{ProductID: items[1].Product.ID, Count: 4},
					},
				}
			},
			expectedErr: errors.New("update error"),
		},
		{
			name: "Failure: UoW transaction error",
			setup: func(uow *mocks.UoWMock) itemApplication.ReserveDto {
				uow.On("Transaction", s.ctx, mock.Anything).Return(errors.New("transaction error")).Once()
				return itemApplication.ReserveDto{
					Items: []itemApplication.ItemDto{
						{ProductID: uuid.New(), Count: 1},
					},
				}
			},
			expectedErr: errors.New("transaction error"),
		},
	}

	for _, tc := range tests {
		tc := tc
		s.Run(tc.name, func() {
			s.T().Parallel()
			uow := mocks.NewUowMock()
			dto := tc.setup(uow)
			useCase := itemApplication.NewUseCase(uow)

			err := useCase.Reserve(s.ctx, dto)

			if tc.expectedErr == nil {
				require.NoError(s.T(), err)
			} else {
				require.Error(s.T(), err)
				require.EqualError(s.T(), err, tc.expectedErr.Error())
			}

			uow.AssertExpectations(s.T())
		})
	}
}

func (s *ItemUseCaseTestSuite) TestRelease() {
	tests := []struct {
		name        string
		setup       func(uow *mocks.UoWMock) itemApplication.ReleaseDto
		expectedErr error
	}{
		{
			name: "Success",
			setup: func(uow *mocks.UoWMock) itemApplication.ReleaseDto {
				items := s.createTestItems(5, 10)

				uow.ItemMock.On("GetAllByProductIDs", s.ctx, items[0].Product.ID, items[1].Product.ID).
					Return(items, nil).Once()
				uow.ItemMock.On("Update", s.ctx, mock.Anything).Return(nil).Twice()

				uow.On("Transaction", s.ctx, mock.Anything).Once()

				return itemApplication.ReleaseDto{
					Items: []itemApplication.ItemDto{
						{ProductID: items[0].Product.ID, Count: 4},
						{ProductID: items[1].Product.ID, Count: 9},
					},
				}
			},
			expectedErr: nil,
		},
		{
			name: "Failure: Items not found",
			setup: func(uow *mocks.UoWMock) itemApplication.ReleaseDto {
				items := s.createTestItems(5, 10)

				uow.ItemMock.On("GetAllByProductIDs", s.ctx, items[0].Product.ID, items[1].Product.ID).
					Return(items, itemRepository.ErrItemNotFound).Once()
				uow.On("Transaction", s.ctx, mock.Anything).Once()

				return itemApplication.ReleaseDto{
					Items: []itemApplication.ItemDto{
						{ProductID: items[0].Product.ID, Count: 4},
						{ProductID: items[1].Product.ID, Count: 9},
					},
				}
			},
			expectedErr: itemRepository.ErrItemNotFound,
		},
		{
			name: "Failure: Update items error",
			setup: func(uow *mocks.UoWMock) itemApplication.ReleaseDto {
				items := s.createTestItems(1, 10)

				uow.ItemMock.On("GetAllByProductIDs", s.ctx, items[0].Product.ID, items[1].Product.ID).
					Return(items, nil).Once()
				uow.ItemMock.On("Update", s.ctx, mock.Anything).Return(errors.New("update error")).Once()

				uow.On("Transaction", s.ctx, mock.Anything).Once()

				return itemApplication.ReleaseDto{
					Items: []itemApplication.ItemDto{
						{ProductID: items[0].Product.ID, Count: 4},
						{ProductID: items[1].Product.ID, Count: 9},
					},
				}
			},
			expectedErr: errors.New("update error"),
		},
		{
			name: "Failure: UoW transaction error",
			setup: func(uow *mocks.UoWMock) itemApplication.ReleaseDto {
				uow.On("Transaction", s.ctx, mock.Anything).
					Return(errors.New("transaction error")).Once()
				return itemApplication.ReleaseDto{
					Items: []itemApplication.ItemDto{
						{ProductID: uuid.New(), Count: 1},
					},
				}
			},
			expectedErr: errors.New("transaction error"),
		},
	}

	for _, tc := range tests {
		tc := tc
		s.Run(tc.name, func() {
			s.T().Parallel()
			uow := mocks.NewUowMock()
			dto := tc.setup(uow)
			useCase := itemApplication.NewUseCase(uow)

			err := useCase.Release(s.ctx, dto)

			if tc.expectedErr == nil {
				require.NoError(s.T(), err)
			} else {
				require.Error(s.T(), err)
				require.EqualError(s.T(), err, tc.expectedErr.Error())
			}

			uow.AssertExpectations(s.T())
		})
	}
}

func (s *ItemUseCaseTestSuite) TestGetAll() {
	tests := []struct {
		name        string
		setup       func(uow *mocks.UoWMock) []*itemDomain.Item
		expectedErr error
	}{
		{
			name: "Success",
			setup: func(uow *mocks.UoWMock) []*itemDomain.Item {
				var items []*itemDomain.Item
				uow.ItemMock.On("GetAll", s.ctx).Return(items, nil).Once()
				return items
			},
			expectedErr: nil,
		},
		{
			name: "Failure: Item repository get all error",
			setup: func(uow *mocks.UoWMock) []*itemDomain.Item {
				var items []*itemDomain.Item
				uow.ItemMock.On("GetAll", s.ctx).
					Return(items, errors.New("get all error")).Once()
				return items
			},
			expectedErr: errors.New("get all error"),
		},
	}

	for _, tc := range tests {
		tc := tc
		s.Run(tc.name, func() {
			s.T().Parallel()
			uow := mocks.NewUowMock()
			expectedItems := tc.setup(uow)
			useCase := itemApplication.NewUseCase(uow)

			items, err := useCase.GetAll(s.ctx)

			if tc.expectedErr == nil {
				require.NoError(s.T(), err)
			} else {
				require.Error(s.T(), err)
				require.EqualError(s.T(), err, tc.expectedErr.Error())
			}
			require.Equal(s.T(), expectedItems, items)

			uow.AssertExpectations(s.T())
		})
	}
}

func TestItemUseCaseTestSuite(t *testing.T) {
	suite.Run(t, new(ItemUseCaseTestSuite))
}
