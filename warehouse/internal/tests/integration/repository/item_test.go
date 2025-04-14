//go:build integration

package repository

import (
	"context"
	"testing"
	itemDomain "warehouse/internal/domain/item"
	productDomain "warehouse/internal/domain/product"
	"warehouse/internal/infrastructure/db/migrations"
	itemRepository "warehouse/internal/infrastructure/repository/item"
	productRepository "warehouse/internal/infrastructure/repository/product"
	"warehouse/internal/tests/testutils"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type ItemRepositoryTestSuite struct {
	suite.Suite
	ctx    context.Context
	testDB *testutils.TestDB
}

func (s *ItemRepositoryTestSuite) SetupSuite() {
	config, err := migrations.NewConfig()
	require.NoError(s.T(), err)

	s.ctx = context.Background()

	s.testDB, err = testutils.NewTestDB(s.ctx, config)
	require.NoError(s.T(), err)
}

func (s *ItemRepositoryTestSuite) TearDownSuite() {
	if s.testDB != nil {
		err := s.testDB.Close(s.ctx)
		require.NoError(s.T(), err)
	}
}

func (s *ItemRepositoryTestSuite) getItemRepo() itemDomain.Repository {
	return itemRepository.New(s.testDB.DB)
}

func (s *ItemRepositoryTestSuite) getProductRepo() productDomain.Repository {
	return productRepository.New(s.testDB.DB)
}

func (s *ItemRepositoryTestSuite) createTestProduct() *productDomain.Product {
	product, _, err := productDomain.Create("Test Product", decimal.NewFromInt(100), "test.png")
	require.NoError(s.T(), err)
	return product
}

func (s *ItemRepositoryTestSuite) createTestItem() *itemDomain.Item {
	product := s.createTestProduct()
	item, err := itemDomain.Create(product, 10)
	require.NoError(s.T(), err)
	return item
}

func (s *ItemRepositoryTestSuite) createTestItemInDb(
	productRepo productDomain.Repository,
	itemRepo itemDomain.Repository,
) *itemDomain.Item {
	item := s.createTestItem()

	err := productRepo.Create(s.ctx, item.Product)
	require.NoError(s.T(), err)
	err = itemRepo.Create(s.ctx, item)
	require.NoError(s.T(), err)

	return item
}

func (s *ItemRepositoryTestSuite) TestCreate() {
	tests := []struct {
		name          string
		setup         func(productRepo productDomain.Repository, itemRepo itemDomain.Repository) *itemDomain.Item
		expectedError error
	}{
		{
			name: "Success",
			setup: func(productRepo productDomain.Repository, itemRepo itemDomain.Repository) *itemDomain.Item {
				item := s.createTestItem()
				err := productRepo.Create(s.ctx, item.Product)
				require.NoError(s.T(), err)
				return item
			},
			expectedError: nil,
		},
		{
			name: "Failure: Product not found",
			setup: func(_ productDomain.Repository, _ itemDomain.Repository) *itemDomain.Item {
				return s.createTestItem()
			},
			expectedError: productRepository.ErrProductNotFound,
		},
		{
			name: "Failure: Product already exists",
			setup: func(productRepo productDomain.Repository, itemRepo itemDomain.Repository) *itemDomain.Item {
				return s.createTestItemInDb(productRepo, itemRepo)
			},
			expectedError: itemRepository.ErrItemAlreadyExists,
		},
	}

	productRepo := s.getProductRepo()
	itemRepo := s.getItemRepo()
	for _, tc := range tests {
		s.Run(tc.name, func() {
			item := tc.setup(productRepo, itemRepo)

			err := itemRepo.Create(s.ctx, item)

			if tc.expectedError != nil {
				require.Error(s.T(), err)
				require.Equal(s.T(), tc.expectedError, err)
			} else {
				require.NoError(s.T(), err)

				createdItem, err := itemRepo.GetByID(s.ctx, item.ID)
				require.NoError(s.T(), err)
				require.Equal(s.T(), item.ID, createdItem.ID)
			}
		})
	}
}

func (s *ItemRepositoryTestSuite) TestUpdate() {
	tests := []struct {
		name          string
		setup         func(productRepo productDomain.Repository, itemRepo itemDomain.Repository) *itemDomain.Item
		update        func(item *itemDomain.Item)
		verify        func(updated, expected *itemDomain.Item)
		expectedError error
	}{
		{
			name: "Success: Reserve item count",
			setup: func(productRepo productDomain.Repository, itemRepo itemDomain.Repository) *itemDomain.Item {
				return s.createTestItemInDb(productRepo, itemRepo)
			},
			update: func(item *itemDomain.Item) {
				err := item.Reserve(1)
				require.NoError(s.T(), err)
			},
			verify: func(updated, expected *itemDomain.Item) {
				require.Equal(s.T(), expected.ID, updated.ID)
				require.Equal(s.T(), expected.Count, updated.Count)
			},
			expectedError: nil,
		},
		{
			name: "Success: Release item count",
			setup: func(productRepo productDomain.Repository, itemRepo itemDomain.Repository) *itemDomain.Item {
				return s.createTestItemInDb(productRepo, itemRepo)
			},
			update: func(item *itemDomain.Item) {
				err := item.Release(1)
				require.NoError(s.T(), err)
			},
			verify: func(updated, expected *itemDomain.Item) {
				require.Equal(s.T(), expected.ID, updated.ID)
				require.Equal(s.T(), expected.Count, updated.Count)
			},
			expectedError: nil,
		},
		{
			name: "Failure: Order not found",
			setup: func(_ productDomain.Repository, _ itemDomain.Repository) *itemDomain.Item {
				return s.createTestItem()
			},
			update:        func(order *itemDomain.Item) {},
			verify:        func(updated, expected *itemDomain.Item) {},
			expectedError: itemRepository.ErrItemNotFound,
		},
	}

	productRepo := s.getProductRepo()
	itemRepo := s.getItemRepo()
	for _, tc := range tests {
		s.Run(tc.name, func() {
			order := tc.setup(productRepo, itemRepo)
			tc.update(order)

			err := itemRepo.Update(s.ctx, order)

			if tc.expectedError != nil {
				require.Error(s.T(), err)
				require.Equal(s.T(), tc.expectedError, err)
			} else {
				require.NoError(s.T(), err)

				updatedOrder, err := itemRepo.GetByID(s.ctx, order.ID)
				require.NoError(s.T(), err)
				require.NotNil(s.T(), updatedOrder)
				tc.verify(updatedOrder, order)
			}
		})
	}
}

func (s *ItemRepositoryTestSuite) TestGetByID() {
	tests := []struct {
		name          string
		setup         func(productRepo productDomain.Repository, itemRepo itemDomain.Repository) uuid.UUID
		expectedError error
	}{
		{
			name: "Success",
			setup: func(productRepo productDomain.Repository, itemRepo itemDomain.Repository) uuid.UUID {
				item := s.createTestItemInDb(productRepo, itemRepo)
				return item.ID
			},
			expectedError: nil,
		},
		{
			name: "Failure: Item not found",
			setup: func(_ productDomain.Repository, _ itemDomain.Repository) uuid.UUID {
				return uuid.New()
			},
			expectedError: itemRepository.ErrItemNotFound,
		},
	}

	productRepo := s.getProductRepo()
	itemRepo := s.getItemRepo()
	for _, tc := range tests {
		s.Run(tc.name, func() {
			itemID := tc.setup(productRepo, itemRepo)

			item, err := itemRepo.GetByID(s.ctx, itemID)

			if tc.expectedError != nil {
				require.Error(s.T(), err)
				require.Equal(s.T(), tc.expectedError, err)
			} else {
				require.NoError(s.T(), err)
				require.Equal(s.T(), itemID, item.ID)
			}
		})
	}
}

func (s *ItemRepositoryTestSuite) TestGetAll() {
	tests := []struct {
		name          string
		setup         func(productRepo productDomain.Repository, itemRepo itemDomain.Repository) []uuid.UUID
		expectedError error
	}{
		{
			name: "Success: One item",
			setup: func(productRepo productDomain.Repository, itemRepo itemDomain.Repository) []uuid.UUID {
				item := s.createTestItemInDb(productRepo, itemRepo)
				return []uuid.UUID{item.ID}
			},
			expectedError: nil,
		},
		{
			name: "Success: Multiple items",
			setup: func(productRepo productDomain.Repository, itemRepo itemDomain.Repository) []uuid.UUID {
				var itemIDs []uuid.UUID
				for i := 0; i < 3; i++ {
					item := s.createTestItemInDb(productRepo, itemRepo)
					itemIDs = append(itemIDs, item.ID)
				}
				return itemIDs
			},
			expectedError: nil,
		},
		{
			name: "Success: No items",
			setup: func(_ productDomain.Repository, _ itemDomain.Repository) []uuid.UUID {
				return []uuid.UUID{}
			},
			expectedError: nil,
		},
	}

	productRepo := s.getProductRepo()
	itemRepo := s.getItemRepo()
	for _, tc := range tests {
		s.Run(tc.name, func() {
			itemIDs := tc.setup(productRepo, itemRepo)

			items, err := itemRepo.GetAll(s.ctx)

			if tc.expectedError != nil {
				require.Error(s.T(), err)
				require.Equal(s.T(), tc.expectedError, err)
			} else {
				require.NoError(s.T(), err)

				returnedItemIDs := make([]uuid.UUID, 0, len(items))
				for _, item := range items {
					returnedItemIDs = append(returnedItemIDs, item.ID)
				}
				for _, itemID := range itemIDs {
					require.Contains(s.T(), returnedItemIDs, itemID)
				}
			}
		})
	}
}

func (s *ItemRepositoryTestSuite) TestGetAllByProductIDs() {
	tests := []struct {
		name          string
		setup         func(productRepo productDomain.Repository, itemRepo itemDomain.Repository) []uuid.UUID
		expectedError error
	}{
		{
			name: "Success: One product",
			setup: func(productRepo productDomain.Repository, itemRepo itemDomain.Repository) []uuid.UUID {
				item := s.createTestItemInDb(productRepo, itemRepo)
				return []uuid.UUID{item.Product.ID}
			},
			expectedError: nil,
		},
		{
			name: "Success: Multiple products",
			setup: func(productRepo productDomain.Repository, itemRepo itemDomain.Repository) []uuid.UUID {
				var productIDs []uuid.UUID
				for i := 0; i < 3; i++ {
					item := s.createTestItemInDb(productRepo, itemRepo)
					productIDs = append(productIDs, item.Product.ID)
				}
				return productIDs
			},
			expectedError: nil,
		},
		{
			name: "Success: No products",
			setup: func(_ productDomain.Repository, _ itemDomain.Repository) []uuid.UUID {
				return []uuid.UUID{}
			},
			expectedError: nil,
		},
		{
			name: "Failure: Products not found",
			setup: func(_ productDomain.Repository, _ itemDomain.Repository) []uuid.UUID {
				return []uuid.UUID{uuid.New()}
			},
			expectedError: itemRepository.ErrItemsNotFound,
		},
	}

	productRepo := s.getProductRepo()
	itemRepo := s.getItemRepo()
	for _, tc := range tests {
		s.Run(tc.name, func() {
			productIDs := tc.setup(productRepo, itemRepo)

			items, err := itemRepo.GetAllByProductIDs(s.ctx, productIDs...)

			if tc.expectedError != nil {
				require.Error(s.T(), err)
				require.Equal(s.T(), tc.expectedError, err)
			} else {
				require.NoError(s.T(), err)
				if len(productIDs) > 0 {
					require.Len(s.T(), items, len(productIDs))
					for _, item := range items {
						require.Contains(s.T(), productIDs, item.Product.ID)
					}
				} else {
					require.Empty(s.T(), items)
				}
			}
		})
	}
}

func TestItemRepository(t *testing.T) {
	suite.Run(t, new(ItemRepositoryTestSuite))
}
